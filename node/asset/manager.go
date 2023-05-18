package asset

import (
	"bytes"
	"context"
	"crypto"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"syscall"
	"time"

	"github.com/Filecoin-Titan/titan/api"
	"github.com/Filecoin-Titan/titan/api/types"
	"github.com/Filecoin-Titan/titan/node/asset/fetcher"
	titanindex "github.com/Filecoin-Titan/titan/node/asset/index"
	"github.com/Filecoin-Titan/titan/node/asset/storage"
	titanrsa "github.com/Filecoin-Titan/titan/node/rsa"
	validate "github.com/Filecoin-Titan/titan/node/validation"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore"
	format "github.com/ipfs/go-ipld-format"
	legacy "github.com/ipfs/go-ipld-legacy"
	"github.com/ipfs/go-libipfs/blocks"
	carv2 "github.com/ipld/go-car/v2"
	"github.com/ipld/go-car/v2/blockstore"
	"github.com/multiformats/go-multihash"
	"golang.org/x/xerrors"
)

const maxSizeOfCache = 128

// assetWaiter is used by Manager to store waiting assets for pulling
type assetWaiter struct {
	Root   cid.Cid
	Dss    []*types.CandidateDownloadInfo
	puller *assetPuller
}

// Manager is the struct that manages asset pulling and store
type Manager struct {
	// root cid of asset
	waitList     []*assetWaiter
	waitListLock *sync.Mutex
	pullCh       chan bool
	bFetcher     fetcher.BlockFetcher
	lru          *lruCache
	storage.Storage
	api.Scheduler
	pullParallel int
	pullTimeout  int
	pullRetry    int
}

// ManagerOptions is the struct that contains options for Manager
type ManagerOptions struct {
	Storage      storage.Storage
	BFetcher     fetcher.BlockFetcher
	SchedulerAPI api.Scheduler
	PullParallel int
	PullTimeout  int
	PullRetry    int
}

// NewManager creates a new instance of Manager
func NewManager(opts *ManagerOptions) (*Manager, error) {
	lru, err := newLRUCache(opts.Storage, maxSizeOfCache)
	if err != nil {
		return nil, err
	}

	m := &Manager{
		waitList:     make([]*assetWaiter, 0),
		waitListLock: &sync.Mutex{},
		pullCh:       make(chan bool),
		Storage:      opts.Storage,
		bFetcher:     opts.BFetcher,
		lru:          lru,
		Scheduler:    opts.SchedulerAPI,
		pullParallel: opts.PullParallel,
		pullTimeout:  opts.PullTimeout,
		pullRetry:    opts.PullRetry,
	}

	m.restoreWaitListFromStore()

	go m.start()

	return m, nil
}

// startTick is a helper function that is used to check waitList and save puller
func (m *Manager) startTick() {
	for {
		time.Sleep(10 * time.Second)

		if len(m.waitList) > 0 {
			puller := m.puller()
			if puller != nil {
				if err := m.savePuller(puller); err != nil {
					log.Error("save puller error:%s", err.Error())
				}

				log.Debugf("total block %d, done block %d, total size %d, done size %d",
					len(puller.blocksPulledSuccessList)+len(puller.blocksWaitList),
					len(puller.blocksPulledSuccessList),
					puller.totalSize,
					puller.doneSize)
			}
		}

	}
}

// triggerPuller is a helper function that is used to trigger asset downloads
func (m *Manager) triggerPuller() {
	select {
	case m.pullCh <- true:
	default:
	}
}

// start is a helper function that starts the Manager and begins downloading assets
func (m *Manager) start() {
	if m.bFetcher == nil {
		log.Panic("m.bFetcher == nil")
	}

	go m.startTick()

	// delay 15 second to pull asset if exist waitList
	time.AfterFunc(15*time.Second, m.triggerPuller)

	for {
		<-m.pullCh
		m.pullAssets()
	}
}

// pullAssets pulls all assets that are waiting to be pulled
func (m *Manager) pullAssets() {
	for len(m.waitList) > 0 {
		m.doPullAsset()
	}
}

// doPullAsset pulls a single asset from the waitList
func (m *Manager) doPullAsset() {
	cw := m.headFromWaitList()
	if cw == nil {
		return
	}
	defer m.removeAssetFromWaitList(cw.Root)

	opts := &pullerOptions{cw.Root, cw.Dss, m.Storage, m.bFetcher, m.pullParallel, m.pullTimeout, m.pullRetry}
	assetPuller, err := m.restoreAssetPullerOrNew(opts)
	if err != nil {
		log.Errorf("restore asset puller error:%s", err)
		return
	}

	cw.puller = assetPuller
	err = assetPuller.pullAsset()
	if err != nil {
		log.Errorf("pull asset error: %s", err)
	}

	m.onPullAssetFinish(assetPuller)
}

// headFromWaitList returns the first assetWaiter in waitList, which is the oldest asset waiting to be downloaded

func (m *Manager) headFromWaitList() *assetWaiter {
	m.waitListLock.Lock()
	defer m.waitListLock.Unlock()

	if len(m.waitList) == 0 {
		return nil
	}
	return m.waitList[0]
}

// removeAssetFromWaitList removes an assetWaiter from waitList by the root CID
// and returns the removed assetWaiter if it exists, otherwise returns nil
func (m *Manager) removeAssetFromWaitList(root cid.Cid) *assetWaiter {
	m.waitListLock.Lock()
	defer m.waitListLock.Unlock()

	if len(m.waitList) == 0 {
		return nil
	}

	for i, cw := range m.waitList {
		if cw.Root.Hash().String() == root.Hash().String() {
			m.waitList = append(m.waitList[:i], m.waitList[i+1:]...)

			if err := m.saveWaitList(); err != nil {
				log.Errorf("save wait list error: %s", err.Error())
			}
			return cw
		}
	}

	return nil
}

// addToWaitList adds an assetWaiter to waitList if the asset with the root CID is not already waiting to be downloaded
func (m *Manager) addToWaitList(root cid.Cid, dss []*types.CandidateDownloadInfo) {
	m.waitListLock.Lock()
	defer m.waitListLock.Unlock()

	for _, waiter := range m.waitList {
		if waiter.Root.Hash().String() == root.Hash().String() {
			return
		}
	}

	cw := &assetWaiter{Root: root, Dss: dss}
	m.waitList = append(m.waitList, cw)

	if err := m.saveWaitList(); err != nil {
		log.Errorf("save wait list error: %s", err.Error())
	}

	m.triggerPuller()
}

// savePuller saves the assetPuller to the data store if it exists and has not completed downloading
func (m *Manager) savePuller(puller *assetPuller) error {
	if puller == nil || puller.isPulledComplete() {
		return nil
	}

	buf, err := puller.encode()
	if err != nil {
		return err
	}
	return m.StorePuller(puller.root, buf)
}

// onPullAssetFinish is called when an assetPuller finishes downloading an asset
func (m *Manager) onPullAssetFinish(puller *assetPuller) {
	log.Debugf("onPullAssetFinish, asset %s", puller.root.String())

	if puller.isPulledComplete() {
		if err := m.DeletePuller(puller.root); err != nil && !os.IsNotExist(err) {
			log.Errorf("remove asset puller error:%s", err.Error())
		}

		blockCountOfAsset := uint32(len(puller.blocksPulledSuccessList))
		if err := m.SetBlockCount(context.Background(), puller.root, blockCountOfAsset); err != nil {
			log.Errorf("set block count error:%s", err.Error())
		}

		if err := m.StoreAsset(context.Background(), puller.root); err != nil {
			log.Errorf("store asset error: %s", err.Error())
		}

		if err := m.AddAssetToView(context.Background(), puller.root); err != nil {
			log.Errorf("add asset to view error: %s", err.Error())
		}

	} else {
		if err := m.savePuller(puller); err != nil {
			log.Errorf("save puller error:%s", err.Error())
		}
	}

	if err := m.submitPullerWorkloadReport(puller); err != nil {
		log.Errorf("submitPullerWorkloadReport error %s", err.Error())
	}
}

// saveWaitList encodes the waitList and stores it in the datastore.
func (m *Manager) saveWaitList() error {
	data, err := encode(&m.waitList)
	if err != nil {
		return err
	}

	return m.StoreWaitList(data)
}

// restoreWaitListFromStore retrieves the waitList from the datastore and decodes it.
func (m *Manager) restoreWaitListFromStore() {
	data, err := m.GetWaitList()
	if err != nil {
		if err != datastore.ErrNotFound {
			log.Errorf("restoreWaitListFromStore error:%s", err)
		}
		return
	}

	if len(data) == 0 {
		return
	}

	err = decode(data, &m.waitList)
	if err != nil {
		log.Errorf("restoreWaitListFromStore error:%s", err)
		return
	}

	log.Debugf("restoreWaitListFromStore:%#v", m.waitList)
}

// waitListLen returns the number of items in the waitList.
func (m *Manager) waitListLen() int {
	return len(m.waitList)
}

// Puller returns the asset puller associated with the first waiting item in the waitList.
func (m *Manager) puller() *assetPuller {
	for _, cw := range m.waitList {
		if cw.puller != nil {
			return cw.puller
		}
	}
	return nil
}

// DeleteAsset removes an asset from the datastore and the waitList.
func (m *Manager) DeleteAsset(root cid.Cid) error {
	// remove lru puller
	m.lru.remove(root)

	if _, err := m.deleteAssetFromWaitList(root); err != nil {
		return err
	}

	if ok, err := m.PullerExists(root); err == nil && ok {
		m.DeletePuller(root)
	}

	if err := m.Storage.DeleteAsset(root); err != nil {
		if e, ok := err.(*os.PathError); !ok {
			return err
		} else if e.Err != syscall.ENOENT {
			return err
		}

	}

	return m.RemoveAssetFromView(context.Background(), root)
}

// restoreAssetPullerOrNew retrieves the asset puller associated with the given root CID, or creates a new one.
func (m *Manager) restoreAssetPullerOrNew(opts *pullerOptions) (*assetPuller, error) {
	data, err := m.GetPuller(opts.root)
	if err != nil && !os.IsNotExist(err) {
		log.Errorf("get asset puller error %s", err.Error())
		return nil, err
	}

	cc := newAssetPuller(opts)
	if len(data) > 0 {
		err = cc.decode(data)
		if err != nil {
			return nil, err
		}

		// cover new download sources
		if len(opts.dss) > 0 {
			cc.downloadSources = opts.dss
		}
	}
	return cc, nil
}

// deleteAssetFromWaitList removes an asset from the waitList.
// return true if exist in waitList
func (m *Manager) deleteAssetFromWaitList(root cid.Cid) (bool, error) {
	if c := m.removeAssetFromWaitList(root); c != nil {
		if c.puller != nil {
			err := c.puller.cancelPulling()
			if err != nil {
				return false, err
			}
		}

		return true, nil
	}
	return false, nil
}

// cachedStatus returns the asset status of a given root CID
func (m *Manager) assetStatus(root cid.Cid) (types.ReplicaStatus, error) {
	if ok, err := m.AssetExists(root); err == nil && ok {
		return types.ReplicaStatusSucceeded, nil
	}

	for _, aw := range m.waitList {
		if aw.Root.Hash().String() == root.Hash().String() {
			if aw.puller != nil {
				return types.ReplicaStatusPulling, nil
			}
			return types.ReplicaStatusWaiting, nil
		}
	}

	return types.ReplicaStatusFailed, nil
}

// progressForAssetPulledFailed returns the progress of a failed asset pull operation
func (m *Manager) progressForAssetPulledFailed(root cid.Cid) (*types.AssetPullProgress, error) {
	progress := &types.AssetPullProgress{
		CID:    root.String(),
		Status: types.ReplicaStatusFailed,
	}

	data, err := m.GetPuller(root)
	if os.IsNotExist(err) {
		return progress, nil
	}

	if err != nil {
		return nil, err
	}

	cc := &assetPuller{}
	err = cc.decode(data)
	if err != nil {
		return nil, err
	}

	progress.BlocksCount = len(cc.blocksPulledSuccessList) + len(cc.blocksWaitList)
	progress.DoneBlocksCount = len(cc.blocksPulledSuccessList)
	progress.Size = int64(cc.totalSize)
	progress.DoneSize = int64(cc.doneSize)

	return progress, nil
}

// GetBlock returns the block with the given CID from the LRU cache
func (m *Manager) GetBlock(ctx context.Context, root, block cid.Cid) (blocks.Block, error) {
	return m.lru.getBlock(ctx, root, block)
}

// HasBlock checks if a block with the given CID exists in the LRU cache
func (m *Manager) HasBlock(ctx context.Context, root, block cid.Cid) (bool, error) {
	return m.lru.hasBlock(ctx, root, block)
}

// GetBlocksOfAsset returns a random selection of blocks for the given root CID
// return map, key is random number, value is cid string
func (m *Manager) GetBlocksOfAsset(root cid.Cid, randomSeed int64, randomCount int) ([]string, error) {
	r := rand.New(rand.NewSource(randomSeed))

	idx, err := m.lru.assetIndex(root)
	if err != nil {
		return nil, err
	}

	multiIndex, ok := idx.(*titanindex.MultiIndexSorted)
	if !ok {
		return nil, xerrors.Errorf("idx is not MultiIndexSorted")
	}

	sizeOfBuckets := multiIndex.BucketCount()
	rets := make([]string, 0, randomCount)

	for i := 0; i < randomCount; i++ {
		index := r.Intn(int(sizeOfBuckets))
		records, err := multiIndex.GetBucketRecords(uint32(index))
		if err != nil {
			return nil, err
		}

		if len(records) == 0 {
			return nil, xerrors.Errorf("record is empty")
		}

		index = r.Intn(len(records))
		record := records[index]
		rets = append(rets, record.Cid.String())
	}

	return rets, nil
}

// AddLostAsset adds a lost asset to the Manager's waitList if it is not already present in the storage
func (m *Manager) AddLostAsset(root cid.Cid) error {
	if has, err := m.AssetExists(root); err != nil {
		return err
	} else if has {
		return nil
	}

	switch types.RunningNodeType {
	case types.NodeCandidate:
		m.addToWaitList(root, nil)
	case types.NodeEdge:
		downloadInfos, err := m.GetCandidateDownloadInfos(context.Background(), root.String())
		if err != nil {
			return xerrors.Errorf("get candidate download infos: %w", err.Error())
		}
		m.addToWaitList(root, downloadInfos)
	default:
		return fmt.Errorf("not support node type:%s", types.RunningNodeType)
	}

	return nil
}

// GetAssetsOfBucket retrieves the list of assets in a given bucket ID from the storage
func (m *Manager) GetAssetsOfBucket(ctx context.Context, bucketID uint32, isRemote bool) ([]cid.Cid, error) {
	if !isRemote {
		return m.Storage.GetAssetsInBucket(ctx, bucketID)
	}

	multiHashes, err := m.GetAssetListForBucket(ctx, bucketID)
	if err != nil {
		return nil, err
	}

	cidList := make([]cid.Cid, 0, len(multiHashes))
	for _, multiHash := range multiHashes {
		hash, err := multihash.FromHexString(multiHash)
		if err != nil {
			log.Errorf("new multi hash from string error %s", err.Error())
			continue
		}
		cidList = append(cidList, cid.NewCidV0(hash))
	}
	return cidList, nil
}

// GetAssetForValidation returns a new instance of asset based on a given random seed
func (m *Manager) GetAssetForValidation(ctx context.Context, randomSeed int64) (validate.Asset, error) {
	return NewRandomCheck(randomSeed, m.Storage, m.lru)
}

func (m *Manager) ScanBlocks(ctx context.Context, root cid.Cid) error {
	reader, err := m.GetAsset(root)
	if err != nil {
		return err
	}

	f, ok := reader.(*os.File)
	if !ok {
		return xerrors.Errorf("can not convert asset %s reader to file", root.String())
	}

	bs, err := blockstore.NewReadOnly(f, nil, carv2.ZeroLengthSectionAsEOF(true))
	if err != nil {
		return err
	}

	block, err := bs.Get(ctx, root)
	if err != nil {
		return xerrors.Errorf("get block %s error %w", root.String(), err)
	}

	node, err := legacy.DecodeNode(context.Background(), block)
	if err != nil {
		log.Errorf("decode block error:%s", err.Error())
		return err
	}

	if len(node.Links()) > 0 {
		return m.getNodes(ctx, bs, node.Links())
	}

	return nil

}

func (m *Manager) getNodes(ctx context.Context, bs *blockstore.ReadOnly, links []*format.Link) error {
	for _, link := range links {
		block, err := bs.Get(context.Background(), link.Cid)
		if err != nil {
			return xerrors.Errorf("get block %s error %w", link.Cid.String(), err)
		}

		node, err := legacy.DecodeNode(context.Background(), block)
		if err != nil {
			log.Errorf("decode block error:%s", err.Error())
			return err
		}

		if len(node.Links()) > 0 {
			err = m.getNodes(ctx, bs, node.Links())
			if err != nil {
				return err
			}
		}

	}

	return nil
}

func (m *Manager) submitPullerWorkloadReport(puller *assetPuller) error {
	if types.RunningNodeType != types.NodeEdge {
		return nil
	}

	buf, err := puller.encodeWorkloadReports()
	if err != nil {
		return err
	}

	// TODO: update and get scheduler publicKey from same place
	pem, err := m.GetSchedulerPublicKey(context.Background())
	if err != nil {
		return err
	}

	publicKey, err := titanrsa.Pem2PublicKey([]byte(pem))
	if err != nil {
		return err
	}

	titanRsa := titanrsa.New(crypto.SHA256, crypto.SHA256.New())
	cipherText, err := titanRsa.Encrypt(buf, publicKey)
	if err != nil {
		return err
	}

	return m.SubmitUserWorkloadReport(context.Background(), bytes.NewBuffer(cipherText))
}
