package assets

import (
	"context"
	"crypto"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/Filecoin-Titan/titan/api/types"
	"github.com/filecoin-project/go-statemachine"
	"github.com/ipfs/go-datastore"

	"github.com/Filecoin-Titan/titan/node/modules/dtypes"

	"github.com/Filecoin-Titan/titan/node/cidutil"
	titanrsa "github.com/Filecoin-Titan/titan/node/rsa"
	"github.com/Filecoin-Titan/titan/node/scheduler/db"
	"github.com/Filecoin-Titan/titan/node/scheduler/node"
	logging "github.com/ipfs/go-log/v2"
	"golang.org/x/xerrors"
)

var log = logging.Logger("asset")

const (
	// The number of pull replica in the first stage
	seedReplicaCount = 1
	// Interval to get asset pull progress from node (Unit:Second)
	pullProgressInterval = 20 * time.Second
	// Maximum number of concurrent asset pulls
	assetPullTaskLimit = 10
	// Maximum number of replicas per asset
	assetEdgeReplicasLimit = 100
	// The number of retries to select the pull asset node
	selectNodeRetryLimit = 3
	// If the node disk size is greater than this value, pulling will not continue
	maxNodeDiskUsage = 95.0
	// Number of asset buckets in assets view
	numAssetBuckets = 128
	// When the node is offline for more than this value, the scheduler will assign other nodes to pull the assets to increase the reliability of the assets
	maxNodeOfflineTime = 24 * time.Hour
	// Total bandwidth limit provided by the asset (The larger the bandwidth provided, the more backups are required)
	assetBandwidthLimit = 10000 // unit:MiB/s
	// If the node does not reply more than once, the asset pull timeout is determined.
	assetTimeoutLimit = 3
)

// Manager manages asset replicas
type Manager struct {
	nodeMgr            *node.Manager // node manager
	stateMachineWait   sync.WaitGroup
	assetStateMachines *statemachine.StateGroup
	lock               sync.Mutex
	pullingAssets      map[string]int                // Assignments where assets are being pulled
	config             dtypes.GetSchedulerConfigFunc // scheduler config
	*db.SQLDB
}

// NewManager returns a new AssetManager instance
func NewManager(nodeManager *node.Manager, ds datastore.Batching, configFunc dtypes.GetSchedulerConfigFunc, sdb *db.SQLDB) *Manager {
	m := &Manager{
		nodeMgr:       nodeManager,
		pullingAssets: make(map[string]int),
		config:        configFunc,
		SQLDB:         sdb,
	}

	// state machine initialization
	m.stateMachineWait.Add(1)
	m.assetStateMachines = statemachine.New(ds, m, AssetPullingInfo{})

	return m
}

// Start initializes and starts the asset state machine and associated tickers
func (m *Manager) Start(ctx context.Context) {
	if err := m.initStateMachines(ctx); err != nil {
		log.Errorf("restartStateMachines err: %s", err.Error())
	}
	go m.startCheckAssetsTimer(ctx)
	go m.startCheckPullProgressesTimer(ctx)
}

// Terminate stops the asset state machine
func (m *Manager) Terminate(ctx context.Context) error {
	return m.assetStateMachines.Stop(ctx)
}

// startCheckAssetsTimer Periodically Check for expired assets, check for missing replicas of assets
func (m *Manager) startCheckAssetsTimer(ctx context.Context) {
	now := time.Now()

	nextTime := time.Date(now.Year(), now.Month(), now.Day(), 2, 0, 0, 0, now.Location())
	if now.After(nextTime) {
		nextTime = nextTime.Add(24 * time.Hour)
	}

	duration := nextTime.Sub(now)

	timer := time.NewTimer(duration)
	defer timer.Stop()

	for {
		<-timer.C

		log.Debugln("start assets check ")

		m.processExpiredAssets()
		m.processMissingAssetReplicas()

		timer.Reset(24 * time.Hour)
	}
}

// startCheckPullProgressesTimer Periodically gets asset pull progress
func (m *Manager) startCheckPullProgressesTimer(ctx context.Context) {
	ticker := time.NewTicker(pullProgressInterval)
	defer ticker.Stop()

	for {
		<-ticker.C
		m.retrieveNodePullProgresses()
	}
}

func (m *Manager) setAssetTimeout(hash string) error {
	// update replicas status
	err := m.UpdateReplicasStatusToFailed(hash)
	if err != nil {
		return xerrors.Errorf("setAssetTimeout %s UpdateReplicasStatusToFailed err:%s", hash, err.Error())
	}

	err = m.assetStateMachines.Send(AssetHash(hash), PullFailed{error: xerrors.New("node pull asset response timeout")})
	if err != nil {
		return xerrors.Errorf("setAssetTimeout %s send time out err:%s", hash, err.Error())
	}

	return nil
}

func (m *Manager) retrieveNodePullProgresses() {
	nodePulls := make(map[string][]string)

	// Process pulling assets
	for hash := range m.pullingAssets {
		if m.pullingAssets[hash] >= assetTimeoutLimit {
			m.setAssetTimeout(hash)
			continue
		}
		m.pullingAssets[hash]++

		cid, err := cidutil.HashToCID(hash)
		if err != nil {
			log.Errorf("retrieveNodePullProgresses %s HashString2CIDString err:%s", hash, err.Error())
			continue
		}

		nodes, err := m.LoadUnfinishedPullAssetNodes(hash)
		if err != nil {
			log.Errorf("retrieveNodePullProgresses %s LoadUnfinishedPullAssetNodes err:%s", hash, err.Error())
			continue
		}

		for _, nodeID := range nodes {
			list := nodePulls[nodeID]
			nodePulls[nodeID] = append(list, cid)
		}

	}

	getCP := func(nodeID string, cids []string) {
		// request node
		result, err := m.requestNodePullProgresses(nodeID, cids)
		if err != nil {
			log.Errorf("retrieveNodePullProgresses %s requestNodePullProgresses err:%s", nodeID, err.Error())
			return
		}

		// update asset info
		m.updateAssetPullResults(nodeID, result)
	}

	for nodeID, cids := range nodePulls {
		go getCP(nodeID, cids)
	}
}

func (m *Manager) requestNodePullProgresses(nodeID string, cids []string) (result *types.PullResult, err error) {
	node := m.nodeMgr.GetNode(nodeID)
	if node == nil {
		err = xerrors.Errorf("node %s not found", nodeID)
		return
	}

	result, err = node.GetAssetProgresses(context.Background(), cids)
	return
}

// CreateAssetPullTask creates a new asset pull task
func (m *Manager) CreateAssetPullTask(info *types.PullAssetReq, userID string) error {
	// Waiting for state machine initialization
	m.stateMachineWait.Wait()

	if len(m.pullingAssets) >= assetPullTaskLimit {
		return xerrors.Errorf("The asset in the pulling exceeds the limit %d, please wait", assetPullTaskLimit)
	}

	if info.Replicas > assetEdgeReplicasLimit {
		return xerrors.Errorf("The number of replicas %d exceeds the limit %d", info.Replicas, assetEdgeReplicasLimit)
	}

	if info.Bandwidth > assetBandwidthLimit {
		return xerrors.Errorf("The number of bandwidthDown %d exceeds the limit %d", info.Bandwidth, assetBandwidthLimit)
	}

	log.Infof("asset event: %s, add asset replica: %d,expiration: %s", info.CID, info.Replicas, info.Expiration.String())

	assetRecord, err := m.LoadAssetRecord(info.Hash)
	if err != nil && err != sql.ErrNoRows {
		return xerrors.Errorf("LoadAssetRecord err:%s", err.Error())
	}

	if assetRecord == nil {
		record := &types.AssetRecord{
			Hash:                  info.Hash,
			CID:                   info.CID,
			ServerID:              m.nodeMgr.ServerID,
			NeedEdgeReplica:       info.Replicas,
			NeedCandidateReplicas: int64(m.GetCandidateReplicaCount()),
			Expiration:            info.Expiration,
			NeedBandwidth:         info.Bandwidth,
			State:                 SeedSelect.String(),
		}

		event := &types.AssetEventInfo{Hash: info.Hash, Event: types.AssetEventAdd, Requester: userID}

		err = m.SaveAssetRecord(record, event)
		if err != nil {
			return xerrors.Errorf("SaveRecordOfAsset err:%s", err.Error())
		}

		// create asset task
		return m.assetStateMachines.Send(AssetHash(info.Hash), PullAssetRestart{})
	}

	if exist, _ := m.assetStateMachines.Has(AssetHash(assetRecord.Hash)); !exist {
		return xerrors.Errorf("No operation rights, the asset belongs to another scheduler %s", assetRecord.Hash)
	}

	// Check if the asset is in servicing state
	if assetRecord.State != Servicing.String() && assetRecord.State != Remove.String() {
		return xerrors.Errorf("asset state is %s", assetRecord.State)
	}

	if assetRecord.State == Remove.String() {
		assetRecord.NeedEdgeReplica = 0
		assetRecord.NeedBandwidth = 0
	}

	if info.Replicas <= assetRecord.NeedEdgeReplica && info.Bandwidth <= assetRecord.NeedBandwidth {
		return xerrors.New("No increase in the number of replicas or bandwidth")
	}

	assetRecord.NeedEdgeReplica = info.Replicas
	assetRecord.Expiration = info.Expiration
	assetRecord.NeedBandwidth = info.Bandwidth

	return m.replenishAssetReplicas(assetRecord, 0, userID, "", SeedSelect)
}

// replenishAssetReplicas updates the existing asset replicas if needed
func (m *Manager) replenishAssetReplicas(assetRecord *types.AssetRecord, replenishReplicas int64, userID, details string, state AssetState) error {
	log.Warnf("replenishAssetReplicas SaveRecordOfAsset : %d", replenishReplicas)

	record := &types.AssetRecord{
		Hash:                  assetRecord.Hash,
		CID:                   assetRecord.CID,
		ServerID:              m.nodeMgr.ServerID,
		NeedEdgeReplica:       assetRecord.NeedEdgeReplica,
		NeedCandidateReplicas: int64(m.GetCandidateReplicaCount()),
		Expiration:            assetRecord.Expiration,
		ReplenishReplicas:     replenishReplicas,
		NeedBandwidth:         assetRecord.NeedBandwidth,
		State:                 state.String(),
	}

	event := &types.AssetEventInfo{Hash: assetRecord.Hash, Event: types.AssetEventAdd, Requester: userID, Details: details}

	err := m.SaveAssetRecord(record, event)
	if err != nil {
		return xerrors.Errorf("SaveRecordOfAsset err:%s", err.Error())
	}

	return m.assetStateMachines.Send(AssetHash(assetRecord.Hash), PullAssetRestart{})
}

// RestartPullAssets restarts asset pulls
func (m *Manager) RestartPullAssets(hashes []types.AssetHash) error {
	if len(m.pullingAssets)+len(hashes) >= assetPullTaskLimit {
		return xerrors.Errorf("The asset in the pulling exceeds the limit %d, please wait", assetPullTaskLimit)
	}

	for _, hash := range hashes {
		if exist, _ := m.assetStateMachines.Has(AssetHash(hash)); !exist {
			continue
		}

		err := m.assetStateMachines.Send(AssetHash(hash), PullAssetRestart{})
		if err != nil {
			log.Errorf("RestartPullAssets send err:%s", err.Error())
		}
	}

	return nil
}

// RemoveReplica remove a replica for node
func (m *Manager) RemoveReplica(cid, hash, nodeID, userID string) error {
	err := m.DeleteAssetReplica(hash, nodeID, &types.AssetEventInfo{Hash: hash, Event: types.AssetEventRemove, Requester: userID, Details: nodeID})
	if err != nil {
		return xerrors.Errorf("RemoveReplica %s DeleteAssetReplica err: %s", hash, err.Error())
	}

	// asset view
	err = m.removeAssetFromView(nodeID, cid)
	if err != nil {
		return xerrors.Errorf("RemoveReplica %s removeAssetFromView err: %s", hash, err.Error())
	}

	go m.requestAssetDelete(nodeID, cid)

	return nil
}

// RemoveAsset removes an asset
func (m *Manager) RemoveAsset(cid, hash, userID string) error {
	if exist, _ := m.assetStateMachines.Has(AssetHash(hash)); !exist {
		return xerrors.Errorf("not found asset %s", hash)
	}

	return m.assetStateMachines.Send(AssetHash(hash), AssetForceState{State: Remove, Requester: userID})
}

// updateAssetPullResults updates asset pull results
func (m *Manager) updateAssetPullResults(nodeID string, result *types.PullResult) {
	isCandidate := false
	pullingCount := 0

	nodeInfo := m.nodeMgr.GetNode(nodeID)
	if nodeInfo != nil {
		isCandidate = nodeInfo.Type == types.NodeCandidate
		// update node info
		nodeInfo.DiskUsage = result.DiskUsage
		defer nodeInfo.SetCurPullingCount(pullingCount)
	}

	for _, progress := range result.Progresses {
		log.Debugf("updateAssetPullResults node_id: %s, status: %d, size: %d/%d, cid: %s ", nodeID, progress.Status, progress.DoneSize, progress.Size, progress.CID)

		hash, err := cidutil.CIDToHash(progress.CID)
		if err != nil {
			log.Errorf("%s cid to hash err:%s", progress.CID, err.Error())
			continue
		}

		exist, _ := m.assetStateMachines.Has(AssetHash(hash))
		if !exist {
			continue
		}

		m.resetAssetNoResponseCount(hash)

		if progress.Status == types.ReplicaStatusWaiting {
			pullingCount++
			continue
		}

		// save replica info to db
		cInfo := &types.ReplicaInfo{
			Status:   progress.Status,
			DoneSize: progress.DoneSize,
			Hash:     hash,
			NodeID:   nodeID,
		}

		err = m.UpdateUnfinishedReplica(cInfo)
		if err != nil {
			log.Errorf("updateAssetPullResults %s UpdateUnfinishedReplica err:%s", nodeID, err.Error())
			continue
		}

		if progress.Status == types.ReplicaStatusPulling {
			pullingCount++

			err = m.assetStateMachines.Send(AssetHash(hash), InfoUpdate{
				Blocks: int64(progress.BlocksCount),
				Size:   progress.Size,
			})
			if err != nil {
				log.Errorf("updateAssetPullResults %s statemachine send err:%s", nodeID, err.Error())
			}

			continue
		}

		// asset view
		err = m.addAssetToView(nodeID, progress.CID)
		if err != nil {
			log.Errorf("updateAssetPullResults %s addAssetToView err:%s", nodeID, err.Error())
			continue
		}

		err = m.assetStateMachines.Send(AssetHash(hash), PulledResult{
			ResultInfo: &NodePulledResult{
				NodeID:      nodeID,
				Status:      int64(progress.Status),
				BlocksCount: int64(progress.BlocksCount),
				Size:        progress.Size,
				IsCandidate: isCandidate,
			},
			Msg: "result",
		})
		if err != nil {
			log.Errorf("updateAssetPullResults %s statemachine send err:%s", nodeID, err.Error())
			continue
		}
	}
}

// Reset the number of no response asset tasks
func (m *Manager) resetAssetNoResponseCount(hash string) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.pullingAssets[hash] = 0
}

func (m *Manager) stopAssetCount(key string) {
	m.lock.Lock()
	defer m.lock.Unlock()

	delete(m.pullingAssets, key)
}

// UpdateAssetExpiration updates the asset expiration for a given CID
func (m *Manager) UpdateAssetExpiration(cid string, t time.Time) error {
	hash, err := cidutil.CIDToHash(cid)
	if err != nil {
		return err
	}

	log.Infof("asset event %s, reset asset expiration:%s", cid, t.String())

	err = m.UpdateAssetRecordExpiry(hash, t, m.nodeMgr.ServerID)
	if err != nil {
		return err
	}

	return nil
}

// processMissingAssetReplicas checks for missing replicas of assets and adds missing replicas
func (m *Manager) processMissingAssetReplicas() {
	aRows, err := m.LoadAllAssetRecords(m.nodeMgr.ServerID)
	if err != nil {
		log.Errorf("LoadAllAssetRecords err:%s", err.Error())
		return
	}
	defer aRows.Close()

	// loading asset records
	for aRows.Next() {
		cInfo := &types.AssetRecord{}
		err = aRows.StructScan(cInfo)
		if err != nil {
			log.Errorf("asset StructScan err: %s", err.Error())
			continue
		}

		if cInfo.State != Servicing.String() && cInfo.State != EdgesFailed.String() {
			continue
		}

		effectiveEdges, details, err := m.checkAssetReliability(cInfo.Hash)
		if err != nil {
			log.Errorf("checkAssetReliability err: %s", err.Error())
			continue
		}

		missingEdges := cInfo.NeedEdgeReplica - int64(effectiveEdges)

		if missingEdges <= 0 && cInfo.State == Servicing.String() {
			// Asset are healthy and do not need to be replenish replicas
			continue
		}

		// do replenish replicas
		err = m.replenishAssetReplicas(cInfo, missingEdges, string(m.nodeMgr.ServerID), details, CandidatesSelect)
		if err != nil {
			log.Errorf("replenishAssetReplicas err: %s", err.Error())
			continue
		}
	}
}

// Check the reliability of assets
func (m *Manager) checkAssetReliability(hash string) (effectiveEdges int, details string, outErr error) {
	// loading asset replicas
	rRows, outErr := m.LoadReplicasByHash(hash, []types.ReplicaStatus{types.ReplicaStatusSucceeded})
	if outErr != nil {
		log.Errorf("checkAssetReliability LoadReplicasByHash err:%s", outErr.Error())
		return
	}
	defer rRows.Close()

	details = "offline node:"

	for rRows.Next() {
		rInfo := &types.ReplicaInfo{}
		err := rRows.StructScan(rInfo)
		if err != nil {
			log.Errorf("checkAssetReliability StructScan err: %s", err.Error())
			continue
		}

		// Are the nodes unreliable
		nodeID := rInfo.NodeID
		lastSeen, err := m.LoadNodeLastSeenTime(nodeID)
		if err != nil {
			log.Errorf("checkAssetReliability LoadLastSeenOfNode err: %s", err.Error())
			continue
		}

		if rInfo.IsCandidate {
			continue
		}

		if lastSeen.Add(maxNodeOfflineTime).After(time.Now()) {
			effectiveEdges++
		} else {
			details = fmt.Sprintf("%s%s,", details, nodeID)
		}
	}

	return
}

// processExpiredAssets checks for expired assets and removes them
func (m *Manager) processExpiredAssets() {
	records, err := m.LoadExpiredAssetRecords(m.nodeMgr.ServerID)
	if err != nil {
		log.Errorf("LoadExpiredAssetRecords err:%s", err.Error())
		return
	}

	for _, record := range records {
		// do remove
		err = m.RemoveAsset(record.CID, record.Hash, string(m.nodeMgr.ServerID))
		log.Infof("the asset cid(%s) has expired, being removed, err: %v", record.CID, err)
	}
}

// requestAssetDelete notifies a node to delete an asset by its CID
func (m *Manager) requestAssetDelete(nodeID, cid string) error {
	node := m.nodeMgr.GetNode(nodeID)
	if node != nil {
		return node.DeleteAsset(context.Background(), cid)
	}

	return xerrors.Errorf("node %s not found", nodeID)
}

// GetCandidateReplicaCount get the candidate replica count from the configuration
func (m *Manager) GetCandidateReplicaCount() int {
	cfg, err := m.config()
	if err != nil {
		log.Errorf("get schedulerConfig err:%s", err.Error())
		return seedReplicaCount
	}

	return seedReplicaCount + cfg.CandidateReplicas
}

// GetAssetRecordInfo get the asset record info for cid
func (m *Manager) GetAssetRecordInfo(cid string) (*types.AssetRecord, error) {
	hash, err := cidutil.CIDToHash(cid)
	if err != nil {
		return nil, err
	}

	dInfo, err := m.LoadAssetRecord(hash)
	if err != nil {
		return nil, err
	}

	dInfo.ReplicaInfos, err = m.LoadAssetReplicas(hash)
	if err != nil {
		log.Errorf("GetAssetRecordInfo hash:%s, LoadAssetReplicas err:%s", hash, err.Error())
	}

	return dInfo, err
}

// saveReplicaInformation stores replica information for nodes
func (m *Manager) saveReplicaInformation(nodes map[string]*node.Node, hash string, isCandidate bool) error {
	// save replica info
	replicaInfos := make([]*types.ReplicaInfo, 0)

	for _, node := range nodes {
		replicaInfos = append(replicaInfos, &types.ReplicaInfo{
			NodeID:      node.NodeID,
			Status:      types.ReplicaStatusWaiting,
			Hash:        hash,
			IsCandidate: isCandidate,
		})
	}

	return m.BatchInitReplicas(replicaInfos)
}

// getDownloadSources gets download sources for a given CID
func (m *Manager) getDownloadSources(cid string, nodes []string) []*types.CandidateDownloadInfo {
	titanRsa := titanrsa.New(crypto.SHA256, crypto.SHA256.New())
	workloadRecords := make([]*types.WorkloadRecord, 0)
	sources := make([]*types.CandidateDownloadInfo, 0)
	for _, nodeID := range nodes {
		cNode := m.nodeMgr.GetCandidateNode(nodeID)
		if cNode == nil {
			continue
		}

		token, tkPayload, err := cNode.Token(cid, titanRsa, m.nodeMgr.PrivateKey)
		if err != nil {
			continue
		}
		workloadRecord := &types.WorkloadRecord{TokenPayload: *tkPayload, Status: types.WorkloadStatusCreate}
		workloadRecords = append(workloadRecords, workloadRecord)

		source := &types.CandidateDownloadInfo{
			NodeID:  nodeID,
			Address: cNode.DownloadAddr(),
			Tk:      token,
		}

		sources = append(sources, source)
	}

	if len(workloadRecords) > 0 {
		if err := m.SaveWorkloadRecord(workloadRecords); err != nil {
			log.Errorf("SaveWorkloadRecord error: %s", err.Error())
		}
	}
	return sources
}

// chooseCandidateNodesForAssetReplica selects candidate nodes to pull asset replicas
func (m *Manager) chooseCandidateNodesForAssetReplica(count int, filterNodes []string) (map[string]*node.Node, string) {
	str := fmt.Sprintf("need node:%d , filter node:%d , cur node:%d , randNum : ", count, len(filterNodes), m.nodeMgr.Candidates)

	selectMap := make(map[string]*node.Node)
	if count <= 0 {
		return selectMap, str
	}

	if len(filterNodes) >= m.nodeMgr.Candidates {
		return selectMap, str
	}

	filterMap := make(map[string]struct{})
	for _, nodeID := range filterNodes {
		filterMap[nodeID] = struct{}{}
	}

	num := count * selectNodeRetryLimit

	for i := 0; i < num; i++ {
		node, rNum := m.nodeMgr.GetRandomCandidate()
		str = fmt.Sprintf("%s%d,", str, rNum)

		if node == nil {
			continue
		}
		nodeID := node.NodeID

		if _, exist := filterMap[nodeID]; exist {
			continue
		}

		if node.DiskUsage > maxNodeDiskUsage {
			continue
		}

		selectMap[nodeID] = node
		if len(selectMap) >= count {
			break
		}
	}

	return selectMap, str
}

// chooseEdgeNodesForAssetReplica selects edge nodes to pull asset replicas
func (m *Manager) chooseEdgeNodesForAssetReplica(count int, bandwidthDown int64, filterNodes []string) (map[string]*node.Node, string) {
	str := fmt.Sprintf("need node:%d , filter node:%d , cur node:%d , randNum : ", count, len(filterNodes), m.nodeMgr.Edges)
	selectMap := make(map[string]*node.Node)
	if count <= 0 {
		count = 1
	}

	if len(filterNodes) >= m.nodeMgr.Edges {
		return selectMap, str
	}

	filterMap := make(map[string]struct{})
	for _, nodeID := range filterNodes {
		filterMap[nodeID] = struct{}{}
	}

	for i := 0; i < count*selectNodeRetryLimit; i++ {
		node, rNum := m.nodeMgr.GetRandomEdge()
		str = fmt.Sprintf("%s%d,", str, rNum)
		if node == nil {
			continue
		}
		nodeID := node.NodeID

		if _, exist := filterMap[nodeID]; exist {
			continue
		}

		if node.DiskUsage > maxNodeDiskUsage {
			continue
		}

		if _, exist := selectMap[nodeID]; exist {
			continue
		}

		bandwidthDown -= int64(node.BandwidthDown)
		selectMap[nodeID] = node
		if len(selectMap) >= count && bandwidthDown <= 0 {
			break
		}
	}

	return selectMap, str
}
