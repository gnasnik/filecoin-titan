package scheduler

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rsa"
	"database/sql"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"time"

	"github.com/Filecoin-Titan/titan/node/modules/dtypes"
	"github.com/Filecoin-Titan/titan/node/scheduler/nat"
	"github.com/Filecoin-Titan/titan/node/scheduler/user"
	"github.com/Filecoin-Titan/titan/node/scheduler/validation"
	"github.com/Filecoin-Titan/titan/node/scheduler/workload"

	"go.uber.org/fx"

	"github.com/Filecoin-Titan/titan/node/config"
	"github.com/Filecoin-Titan/titan/node/scheduler/assets"

	"github.com/Filecoin-Titan/titan/api"
	"github.com/Filecoin-Titan/titan/api/types"
	"github.com/Filecoin-Titan/titan/node/common"
	"github.com/Filecoin-Titan/titan/node/handler"
	"github.com/Filecoin-Titan/titan/node/scheduler/node"
	"github.com/filecoin-project/go-jsonrpc/auth"
	logging "github.com/ipfs/go-log/v2"

	titanrsa "github.com/Filecoin-Titan/titan/node/rsa"
	"github.com/Filecoin-Titan/titan/node/scheduler/sync"
	"golang.org/x/xerrors"
)

var log = logging.Logger("scheduler")

// Scheduler represents a scheduler node in a distributed system.
type Scheduler struct {
	fx.In

	*common.CommonAPI
	*EdgeUpdateManager
	dtypes.ServerID

	NodeManager            *node.Manager
	ValidationMgr          *validation.Manager
	AssetManager           *assets.Manager
	NatManager             *nat.Manager
	DataSync               *sync.DataSync
	SchedulerCfg           *config.SchedulerCfg
	SetSchedulerConfigFunc dtypes.SetSchedulerConfigFunc
	GetSchedulerConfigFunc dtypes.GetSchedulerConfigFunc
	WorkloadManager        *workload.Manager

	PrivateKey *rsa.PrivateKey
}

var _ api.Scheduler = &Scheduler{}

type jwtPayload struct {
	Allow  []auth.Permission
	NodeID string
}

// nodeConnect processes a node connect request with the given options and node type.
func (s *Scheduler) nodeConnect(ctx context.Context, opts *types.ConnectOptions, nodeType types.NodeType) error {
	remoteAddr := handler.GetRemoteAddr(ctx)
	nodeID := handler.GetNodeID(ctx)

	alreadyConnect := true

	cNode := s.NodeManager.GetNode(nodeID)
	if cNode == nil {
		if err := s.NodeManager.NodeExists(nodeID, nodeType); err != nil {
			return xerrors.Errorf("node: %s, type: %d, error: %w", nodeID, nodeType, err)
		}
		cNode = node.New()
		alreadyConnect = false
	}
	cNode.SetToken(opts.Token)

	log.Infof("node connected %s, address:%s", nodeID, remoteAddr)

	err := cNode.ConnectRPC(remoteAddr, nodeType)
	if err != nil {
		return xerrors.Errorf("nodeConnect ConnectRPC err:%s", err.Error())
	}

	if !alreadyConnect {
		// init node info
		nodeInfo, err := cNode.API.GetNodeInfo(context.Background())
		if err != nil {
			log.Errorf("nodeConnect NodeInfo err:%s", err.Error())
			return err
		}

		if nodeID != nodeInfo.NodeID {
			return xerrors.Errorf("nodeID mismatch %s, %s", nodeID, nodeInfo.NodeID)
		}

		nodeInfo.Type = nodeType
		nodeInfo.SchedulerID = s.ServerID

		pStr, err := s.NodeManager.LoadNodePublicKey(nodeID)
		if err != nil && err != sql.ErrNoRows {
			return xerrors.Errorf("load node port %s err : %s", nodeID, err.Error())
		}

		publicKey, err := titanrsa.Pem2PublicKey([]byte(pStr))
		if err != nil {
			return xerrors.Errorf("load node port %s err : %s", nodeID, err.Error())
		}

		oldInfo, err := s.NodeManager.LoadNodeInfo(nodeID)
		if err != nil && err != sql.ErrNoRows {
			return xerrors.Errorf("load node online duration %s err : %s", nodeID, err.Error())
		}

		if oldInfo != nil {
			// init node info
			nodeInfo.OnlineDuration = oldInfo.OnlineDuration
			nodeInfo.UploadTraffic = oldInfo.UploadTraffic
			nodeInfo.DownloadTraffic = oldInfo.DownloadTraffic
			nodeInfo.PortMapping = oldInfo.PortMapping
		}
		nodeInfo.ExternalIP, _, err = net.SplitHostPort(remoteAddr)
		if err != nil {
			return xerrors.Errorf("SplitHostPort err:%s", err.Error())
		}

		cNode.SetPublicKey(publicKey)
		cNode.SetTCPPort(opts.TcpServerPort)
		cNode.SetRemoteAddr(remoteAddr)

		cNode.NodeInfo = &nodeInfo

		err = s.NodeManager.NodeOnline(cNode)
		if err != nil {
			log.Errorf("nodeConnect err:%s,nodeID:%s", err.Error(), nodeID)
			return err
		}
	}

	if nodeType == types.NodeEdge {
		go s.NatManager.DetermineEdgeNATType(context.Background(), nodeID)
	}

	s.DataSync.AddNodeToList(nodeID)

	return nil
}

// NodeValidationResult processes the validation result for a node
func (s *Scheduler) NodeValidationResult(ctx context.Context, r io.Reader, sign string) error {
	validator := handler.GetNodeID(ctx)
	node := s.NodeManager.GetNode(validator)
	if node == nil {
		return fmt.Errorf("node %s not online", validator)
	}

	signBuf, err := hex.DecodeString(sign)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	rsa := titanrsa.New(crypto.SHA256, crypto.SHA256.New())
	err = rsa.VerifySign(node.PublicKey(), signBuf, data)
	if err != nil {
		return err
	}

	result := &api.ValidationResult{}
	buffer := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buffer)
	err = dec.Decode(result)
	if err != nil {
		return err
	}

	result.Validator = validator
	s.ValidationMgr.HandleResult(result)

	return nil
}

// TriggerElection triggers a single election for validators.
func (s *Scheduler) TriggerElection(ctx context.Context) error {
	s.ValidationMgr.StartElection()
	return nil
}

// GetValidationResults retrieves a list of validation results.
func (s *Scheduler) GetValidationResults(ctx context.Context, startTime, endTime time.Time, pageNumber, pageSize int) (*types.ListValidationResultRsp, error) {
	svm, err := s.NodeManager.LoadValidationResultInfos(startTime, endTime, pageNumber, pageSize)
	if err != nil {
		return nil, err
	}

	return svm, nil
}

// GetSchedulerPublicKey get server publicKey
func (s *Scheduler) GetSchedulerPublicKey(ctx context.Context) (string, error) {
	if s.PrivateKey == nil {
		return "", fmt.Errorf("scheduler private key not exist")
	}

	publicKey := s.PrivateKey.PublicKey
	pem := titanrsa.PublicKey2Pem(&publicKey)
	return string(pem), nil
}

// GetValidationInfo  get information related to validation and election
func (s *Scheduler) GetValidationInfo(ctx context.Context) (*types.ValidationInfo, error) {
	eTime := s.ValidationMgr.GetNextElectionTime()

	return &types.ValidationInfo{
		NextElectionTime: eTime,
	}, nil
}

// SubmitUserWorkloadReport submits report of workload for User Asset Download
func (s *Scheduler) SubmitUserWorkloadReport(ctx context.Context, r io.Reader) error {
	nodeID := handler.GetNodeID(ctx)
	node := s.NodeManager.GetNode(nodeID)

	log.Warnf("SubmitUserWorkloadReport node:%s", nodeID)

	cipherText, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	titanRsa := titanrsa.New(crypto.SHA256, crypto.SHA256.New())
	data, err := titanRsa.Decrypt(cipherText, s.PrivateKey)
	if err != nil {
		return xerrors.Errorf("decrypt error: %w", err)
	}

	return s.WorkloadManager.HandleUserWorkload(data, node)
}

// SubmitNodeWorkloadReport submits report of workload for node Asset Download
func (s *Scheduler) SubmitNodeWorkloadReport(ctx context.Context, r io.Reader) error {
	nodeID := handler.GetNodeID(ctx)
	node := s.NodeManager.GetNode(nodeID)
	if node == nil {
		return xerrors.Errorf("node %s not exists", nodeID)
	}

	log.Warnf("SubmitNodeWorkloadReport node:%s", nodeID)

	report := &types.NodeWorkloadReport{}
	dec := gob.NewDecoder(r)
	err := dec.Decode(report)
	if err != nil {
		return xerrors.Errorf("decode data to NodeWorkloadReport error: %w", err)
	}

	titanRsa := titanrsa.New(crypto.SHA256, crypto.SHA256.New())
	if err = titanRsa.VerifySign(node.PublicKey(), report.Sign, report.CipherText); err != nil {
		return xerrors.Errorf("verify sign error: %w", err)
	}

	data, err := titanRsa.Decrypt(report.CipherText, s.PrivateKey)
	if err != nil {
		return xerrors.Errorf("decrypt error: %w", err)
	}

	return s.WorkloadManager.HandleNodeWorkload(data, node)
}

// GetWorkloadRecords retrieves a list of workload results.
func (s *Scheduler) GetWorkloadRecords(ctx context.Context, startTime, endTime time.Time, limit, offset int) (*types.ListWorkloadRecordRsp, error) {
	return s.NodeManager.LoadWorkloadRecords(startTime, endTime, limit, offset)
}

// GetWorkloadRecord retrieves workload result.
func (s *Scheduler) GetWorkloadRecord(ctx context.Context, tokenID string) (*types.WorkloadRecord, error) {
	return s.NodeManager.LoadWorkloadRecord(tokenID)
}

// AllocateStorage allocates storage space.
func (s *Scheduler) AllocateStorage(ctx context.Context, userID string) error {
	u := &user.User{ID: userID}
	return u.AllocateStorage(ctx, s.SchedulerCfg.UserFreeStorageSize)
}

// CreateAPIKey creates a key for the client API.
func (s *Scheduler) CreateAPIKey(ctx context.Context, userID, keyName string) (string, error) {
	u := &user.User{ID: userID}
	return u.CreateAPIKey(ctx, keyName)
}

// CreateAsset creates an asset with car CID, car name, and car size.
func (s *Scheduler) CreateAsset(ctx context.Context, req *types.CreateAssetReq) (*types.CreateAssetRsp, error) {
	u := &user.User{ID: req.UserID}
	return u.CreateAsset(ctx, req.AssetID, req.AssetName, req.AssetSize)
}

// ListAssets lists the assets of the user.
func (s *Scheduler) ListAssets(ctx context.Context, userID string) ([]*types.AssetProperty, error) {
	u := &user.User{ID: userID}
	return u.ListAssets(ctx)
}

// DeleteAssets deletes the assets of the user.
func (s *Scheduler) DeleteAssets(ctx context.Context, userID string, assetCIDs []string) error {
	u := &user.User{ID: userID}
	return u.DeleteAssets(ctx, assetCIDs)
}

// ShareAssets shares the assets of the user.
func (s *Scheduler) ShareAssets(ctx context.Context, userID string, assetCIDs []string) ([]string, error) {
	u := &user.User{ID: userID}
	return u.ShareAssets(ctx, assetCIDs)
}
