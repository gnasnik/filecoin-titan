package recordfile

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"github.com/Filecoin-Titan/titan/lib/lotuscli"
	"github.com/filecoin-project/go-address"
	commp "github.com/filecoin-project/go-fil-commp-hashhash"
	"github.com/filecoin-project/go-fil-markets/storagemarket"
	"github.com/filecoin-project/go-state-types/abi"
	lotusapi "github.com/filecoin-project/lotus/api"
	"github.com/google/uuid"
	"github.com/ipfs/go-cid"
	cbg "github.com/whyrusleeping/cbor-gen"
	"golang.org/x/crypto/sha3"
	"io"
	"os"
	"path"
	"time"

	"github.com/Filecoin-Titan/titan/api/types"
	"github.com/Filecoin-Titan/titan/lib/carutil"
	"github.com/Filecoin-Titan/titan/node/modules/dtypes"
	"github.com/Filecoin-Titan/titan/node/scheduler/db"
	"github.com/Filecoin-Titan/titan/node/scheduler/leadership"
	commcid "github.com/filecoin-project/go-fil-commcid"
	lotustypes "github.com/filecoin-project/lotus/chain/types"
	logging "github.com/ipfs/go-log/v2"
	"golang.org/x/xerrors"
)

var log = logging.Logger("recordfile")

const (
	// Process 1000 pieces of validation result data at a time
	vResultLimit = 1000

	oneDay = 24 * time.Hour

	BufSize = (4 << 20) / 128 * 127

	BlockGasLimit = int64(10_000_000_000)
)

// Manager is the node manager responsible for managing the online nodes
type Manager struct {
	*db.SQLDB
	dtypes.ServerID // scheduler server id
	leadershipMgr   *leadership.Manager
	lotusCli        *lotuscli.Client
}

// NewManager creates a new instance of the node manager
func NewManager(sdb *db.SQLDB, serverID dtypes.ServerID, lmgr *leadership.Manager, cli *lotuscli.Client) *Manager {
	mgr := &Manager{
		SQLDB:         sdb,
		leadershipMgr: lmgr,
		ServerID:      serverID,
		lotusCli:      cli,
	}

	go mgr.startTimer()

	return mgr
}

func (m *Manager) startTimer() {
	now := time.Now()

	nextTime := time.Date(now.Year(), now.Month(), now.Day(), 2, 0, 0, 0, now.Location())
	if now.After(nextTime) {
		nextTime = nextTime.Add(oneDay)
	}

	duration := nextTime.Sub(now)

	timer := time.NewTimer(duration)
	defer timer.Stop()

	for {
		<-timer.C

		log.Debugln("start timer...")
		m.handleValidationResultSaveToFiles()
		// Packaging

		timer.Reset(oneDay)
	}
}

func (m *Manager) handleValidationResultSaveToFiles() {
	if !m.leadershipMgr.RequestAndBecomeMaster() {
		return
	}

	defer func() {
		log.Infoln("handleValidationResultSaveToFiles end")
	}()

	log.Infoln("handleValidationResultSaveToFiles start")

	// do handle validation result
	for {
		ids, ds, err := m.loadResults()
		if err != nil {
			log.Errorf("loadResults err:%s", err.Error())
			return
		}

		if len(ids) == 0 {
			return
		}

		m.saveValidationResultToFiles(ds)

		err = m.UpdateFileSavedStatus(ids)
		if err != nil {
			log.Errorf("UpdateFileSavedStatus err:%s", err.Error())
			return
		}
	}
}

func (m *Manager) loadResults() ([]int, map[string]string, error) {
	rows, err := m.LoadUnSavedValidationResults(vResultLimit)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	ids := make([]int, 0)
	ds := make(map[string]string)

	for rows.Next() {
		vInfo := &types.ValidationResultInfo{}
		err = rows.StructScan(vInfo)
		if err != nil {
			log.Errorf("loadResults StructScan err: %s", err.Error())
			continue
		}

		ids = append(ids, vInfo.ID)

		str := fmt.Sprintf("RoundID:%s, NodeID:%s, ValidatorID:%s, Profit:%.2f, ValidationCID:%s,EndTime:%s \n",
			vInfo.RoundID, vInfo.NodeID, vInfo.ValidatorID, vInfo.Profit, vInfo.Cid, vInfo.EndTime.String())

		filename := vInfo.EndTime.Format("20060102")
		data := ds[filename]
		ds[filename] = fmt.Sprintf("%s%s", data, str)
	}

	return ids, ds, nil
}

func (m *Manager) saveValidationResultToFiles(ds map[string]string) {
	saveFile := func(data, filename string) error {
		writer, err := NewWriter(DirectoryNameValidation, filename)
		if err != nil {
			return xerrors.Errorf("NewWriter err:%s", err.Error())
		}
		defer writer.Close()

		err = writer.WriteData(data)
		if err != nil {
			return xerrors.Errorf("WriteData err:%s", err.Error())
		}

		return nil
	}

	for filename, data := range ds {
		err := saveFile(filename, data)
		if err != nil {
			log.Errorf("filename:%s , saveFile err:%s", filename, err.Error())
		}
	}
}

func generateCAR(ctx context.Context, inputFile string, outDir string) (types.GeneratedCarInfo, error) {
	stat, err := os.Stat(inputFile)
	if err != nil {
		return types.GeneratedCarInfo{}, err
	}

	outFilename := uuid.New().String() + ".car"
	outPath := path.Join(outDir, outFilename)
	carF, err := os.Create(outPath)
	if err != nil {
		return types.GeneratedCarInfo{}, err
	}
	cp := new(commp.Calc)
	writer := bufio.NewWriterSize(io.MultiWriter(carF, cp), BufSize)
	input := []carutil.Finfo{
		{
			Path:  inputFile,
			Size:  stat.Size(),
			Start: 0,
			End:   stat.Size(),
		},
	}

	_, rootCid, _, err := carutil.GenerateCar(ctx, input, outDir, "", writer)
	if err != nil {
		return types.GeneratedCarInfo{}, err
	}
	err = writer.Flush()
	if err != nil {
		return types.GeneratedCarInfo{}, err
	}
	err = carF.Close()
	if err != nil {
		return types.GeneratedCarInfo{}, err
	}
	rawCommP, pieceSize, err := cp.Digest()
	if err != nil {
		return types.GeneratedCarInfo{}, err
	}

	commCid, err := commcid.DataCommitmentV1ToCID(rawCommP)
	if err != nil {
		return types.GeneratedCarInfo{}, err
	}

	carFilePath := path.Join(outDir, commCid.String()+".car")
	err = os.Rename(outPath, carFilePath)
	if err != nil {
		return types.GeneratedCarInfo{}, err
	}

	out := types.GeneratedCarInfo{
		DataCid:   rootCid,
		PieceCid:  commCid.String(),
		PieceSize: pieceSize,
		Path:      carFilePath,
	}

	return out, nil
}

// TODO: invoke smart contract instead of rpc call
func createDealProposal(client *lotuscli.Client, act string, from string, carInfo types.GeneratedCarInfo) (*cid.Cid, error) {
	maddr, err := address.NewFromString(act)
	if err != nil {
		return nil, fmt.Errorf("parsing address %s: %w", act, err)
	}

	_, err = client.ImportFile(lotusapi.FileRef{
		Path:  carInfo.Path,
		IsCAR: true,
	})

	if err != nil {
		return nil, xerrors.Errorf("importing file: %w", err)
	}

	dataCid, err := cid.Decode(carInfo.DataCid)
	if err != nil {
		return nil, xerrors.Errorf("getting data cid from bytes: %w", err)
	}

	pieceCid, err := cid.Decode(carInfo.PieceCid)
	if err != nil {
		return nil, xerrors.Errorf("getting piece cid from bytes: %w", err)
	}

	paddedSize := abi.PaddedPieceSize(carInfo.PieceSize)

	ref := &storagemarket.DataRef{
		TransferType: storagemarket.TTGraphsync,
		Root:         dataCid,
		PieceCid:     &pieceCid,
		PieceSize:    paddedSize.Unpadded(),
	}

	addr, err := address.NewFromString(from)
	if err != nil {
		return nil, xerrors.Errorf("failed to parse 'from' address: %w", err)
	}

	sdParams := &lotusapi.StartDealParams{
		Data:              ref,
		Wallet:            addr,
		Miner:             maddr,
		EpochPrice:        lotustypes.NewInt(1000000),
		MinBlocksDuration: 518400, // 6 months
		FastRetrieval:     true,
		VerifiedDeal:      false,
	}

	return client.StartDealProposal(sdParams)
}

func InvokeContractByFuncName(ctx context.Context, client *lotuscli.Client, fromEthAddr string, idEthAddr string, funcSignature string, inputData []byte) ([]byte, error) {
	offset := 32
	length := len(inputData)

	startOffset := buildInputFromUint64(uint64(offset))
	lengthData := buildInputFromUint64(uint64(length))
	entryPoint := CalcFuncSignature(funcSignature)
	payload := inputDataFromArray(inputData)

	calldata := append(entryPoint, startOffset...)
	calldata = append(calldata, lengthData...)
	calldata = append(calldata, payload...)

	var buffer bytes.Buffer
	err := cbg.WriteByteArray(&buffer, calldata)
	if err != nil {
		return nil, err
	}
	calldata = buffer.Bytes()

	from, err := address.NewFromString(fromEthAddr)
	if err != nil {
		return nil, xerrors.Errorf("from address is not a filecoin or eth address", err)
	}

	to, err := address.NewFromString(idEthAddr)
	if err != nil {
		return nil, xerrors.Errorf("to address is not a filecoin or eth address", err)
	}

	res, err := client.MpoolPushMessage(ctx, from, to, calldata)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// CalcFuncSignature function signatures are the first 4 bytes of the hash of the function name and types
func CalcFuncSignature(funcName string) []byte {
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write([]byte(funcName))
	hash := hasher.Sum(nil)
	return hash[:4]
}

// convert a simple byte array into input data which is a left padded 64 byte array
func inputDataFromArray(input []byte) []byte {
	inputData := make([]byte, 64)
	copy(inputData[:len(input)], input[:])
	return inputData
}

func buildInputFromUint64(number uint64) []byte {
	// Convert the number to a binary uint64 array
	binaryNumber := make([]byte, 32)
	binary.BigEndian.PutUint64(binaryNumber[24:], number)
	return binaryNumber
}
