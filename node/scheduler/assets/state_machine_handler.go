package assets

import (
	"math"
	"time"

	"github.com/Filecoin-Titan/titan/api/types"
	"github.com/alecthomas/units"
	"github.com/filecoin-project/go-statemachine"
	"golang.org/x/xerrors"
)

var (
	// MinRetryTime defines the minimum time duration between retries
	MinRetryTime = 1 * time.Minute

	// MaxRetryCount defines the maximum number of retries allowed
	MaxRetryCount = 3
)

// failedCoolDown is called when a retry needs to be attempted and waits for the specified time duration
func failedCoolDown(ctx statemachine.Context, info AssetPullingInfo) error {
	retryStart := time.Now().Add(MinRetryTime)
	if time.Now().Before(retryStart) {
		log.Debugf("%s(%s), waiting %s before retrying", info.State, info.Hash, time.Until(retryStart))
		select {
		case <-time.After(time.Until(retryStart)):
		case <-ctx.Context().Done():
			return ctx.Context().Err()
		}
	}

	return nil
}

// handleSeedSelect handles the selection of seed nodes for asset pull
func (m *Manager) handleSeedSelect(ctx statemachine.Context, info AssetPullingInfo) error {
	log.Debugf("handle select seed: %s", info.CID)

	if len(info.CandidateReplicaSucceeds) >= seedReplicaCount {
		// The number of candidate node replicas has reached the requirement
		return ctx.Send(SkipStep{})
	}

	// find nodes
	nodes, str := m.chooseCandidateNodes(seedReplicaCount, info.CandidateReplicaSucceeds)
	if len(nodes) < 1 {
		return ctx.Send(SelectFailed{error: xerrors.Errorf("node not found; %s", str)})
	}

	// save to db
	err := m.saveReplicaInformation(nodes, info.Hash.String(), true)
	if err != nil {
		return ctx.Send(SelectFailed{error: err})
	}

	m.resetAssetNoResponseCount(info.Hash.String())

	// send a cache request to the node
	go func() {
		for _, node := range nodes {
			err := node.PullAsset(ctx.Context(), info.CID, nil)
			if err != nil {
				log.Errorf("%s pull asset err:%s", node.NodeID, err.Error())
				continue
			}

			node.IncrCurPullingCount(1)
		}
	}()

	return ctx.Send(PullRequestSent{CandidateWaitings: int64(len(nodes))})
}

// handleSeedPulling handles the asset pulling process of seed nodes
func (m *Manager) handleSeedPulling(ctx statemachine.Context, info AssetPullingInfo) error {
	log.Debugf("handle seed pulling, %s", info.CID)

	if len(info.CandidateReplicaSucceeds) >= seedReplicaCount {
		return ctx.Send(PullSucceed{})
	}

	if info.CandidateWaitings == 0 {
		return ctx.Send(PullFailed{error: xerrors.New("node pull failed")})
	}

	return nil
}

// handleCandidatesSelect handles the selection of candidate nodes for asset pull
func (m *Manager) handleCandidatesSelect(ctx statemachine.Context, info AssetPullingInfo) error {
	log.Debugf("handle candidates select, %s", info.CID)

	sources := m.getDownloadSources(info.CID, info.CandidateReplicaSucceeds)
	if len(sources) < 1 {
		return ctx.Send(SelectFailed{error: xerrors.New("source node not found")})
	}

	needCount := info.CandidateReplicas - int64(len(info.CandidateReplicaSucceeds))
	if needCount < 1 {
		// The number of candidate node replicas has reached the requirement
		return ctx.Send(SkipStep{})
	}

	// find nodes
	nodes, str := m.chooseCandidateNodes(int(needCount), info.CandidateReplicaSucceeds)
	if len(nodes) < 1 {
		return ctx.Send(SelectFailed{error: xerrors.Errorf("node not found; %s", str)})
	}

	// save to db
	err := m.saveReplicaInformation(nodes, info.Hash.String(), true)
	if err != nil {
		return ctx.Send(SelectFailed{error: err})
	}

	m.resetAssetNoResponseCount(info.Hash.String())

	// send a pull request to the node
	go func() {
		for _, node := range nodes {
			err := node.PullAsset(ctx.Context(), info.CID, sources)
			if err != nil {
				log.Errorf("%s pull asset err:%s", node.NodeID, err.Error())
				continue
			}

			node.IncrCurPullingCount(1)
		}
	}()

	return ctx.Send(PullRequestSent{CandidateWaitings: int64(len(nodes))})
}

// handleCandidatesPulling handles the asset pulling process of candidate nodes
func (m *Manager) handleCandidatesPulling(ctx statemachine.Context, info AssetPullingInfo) error {
	log.Debugf("handle candidates pulling, %s", info.CID)

	if int64(len(info.CandidateReplicaSucceeds)) >= info.CandidateReplicas {
		return ctx.Send(PullSucceed{})
	}

	if info.CandidateWaitings == 0 {
		return ctx.Send(PullFailed{error: xerrors.New("node pull failed")})
	}

	return nil
}

func (m *Manager) getCurBandwidthUp(nodes []string) int64 {
	bandwidthUp := int64(0)
	for _, nodeID := range nodes {
		node := m.nodeMgr.GetEdgeNode(nodeID)
		if node != nil {
			bandwidthUp += int64(node.BandwidthUp)
		}
	}

	// B to MiB
	v := int64(math.Ceil(float64(bandwidthUp) / float64(units.MiB)))
	return v
}

// handleEdgesSelect handles the selection of edge nodes for asset pull
func (m *Manager) handleEdgesSelect(ctx statemachine.Context, info AssetPullingInfo) error {
	log.Debugf("handle edges select , %s", info.CID)

	needCount := info.EdgeReplicas - int64(len(info.EdgeReplicaSucceeds))

	if info.ReplenishReplicas > 0 {
		// Check to node offline while replenishing the temporary replica
		needCount = info.ReplenishReplicas
	}

	needBandwidth := info.Bandwidth - m.getCurBandwidthUp(info.EdgeReplicaSucceeds)

	if needCount < 1 && needBandwidth <= 0 {
		// The number of edge node replications and the total downlink bandwidth are met
		return ctx.Send(SkipStep{})
	}

	sources := m.getDownloadSources(info.CID, info.CandidateReplicaSucceeds)
	if len(sources) < 1 {
		return ctx.Send(SelectFailed{error: xerrors.New("source node not found")})
	}

	// find nodes
	nodes, str := m.chooseEdgeNodes(int(needCount), needBandwidth, info.EdgeReplicaSucceeds)
	if len(nodes) < 1 {
		return ctx.Send(SelectFailed{error: xerrors.Errorf("node not found; %s", str)})
	}

	// save to db
	err := m.saveReplicaInformation(nodes, info.Hash.String(), false)
	if err != nil {
		return ctx.Send(SelectFailed{error: err})
	}

	m.resetAssetNoResponseCount(info.Hash.String())

	// send a pull request to the node
	go func() {
		for _, node := range nodes {
			err := node.PullAsset(ctx.Context(), info.CID, sources)
			if err != nil {
				log.Errorf("%s pull asset err:%s", node.NodeID, err.Error())
				continue
			}

			node.IncrCurPullingCount(1)
		}
	}()

	return ctx.Send(PullRequestSent{EdgeWaitings: int64(len(nodes))})
}

// handleEdgesPulling handles the asset pulling process of edge nodes
func (m *Manager) handleEdgesPulling(ctx statemachine.Context, info AssetPullingInfo) error {
	log.Debugf("handle edges pulling, %s", info.CID)

	needBandwidth := info.Bandwidth - m.getCurBandwidthUp(info.EdgeReplicaSucceeds)

	if int64(len(info.EdgeReplicaSucceeds)) >= info.EdgeReplicas && needBandwidth <= 0 {
		return ctx.Send(PullSucceed{})
	}

	if info.EdgeWaitings == 0 {
		return ctx.Send(PullFailed{error: xerrors.New("node pull failed")})
	}

	return nil
}

// handleServicing asset pull completed and in service status
func (m *Manager) handleServicing(ctx statemachine.Context, info AssetPullingInfo) error {
	log.Infof("handle servicing: %s", info.CID)
	m.stopAssetCount(info.Hash.String())

	// remove fail replicas
	return m.DeleteUnfinishedReplicas(info.Hash.String())
}

// handlePullsFailed handles the failed state of asset pulling and retries if necessary
func (m *Manager) handlePullsFailed(ctx statemachine.Context, info AssetPullingInfo) error {
	m.stopAssetCount(info.Hash.String())

	if info.RetryCount >= int64(MaxRetryCount) {
		log.Infof("handle pulls failed: %s, retry count: %d", info.CID, info.RetryCount)
		return nil
	}

	log.Debugf("handle pulls failed: %s, retries: %d", info.CID, info.RetryCount)

	if err := failedCoolDown(ctx, info); err != nil {
		return err
	}

	return ctx.Send(AssetRePull{})
}

func (m *Manager) handleRemove(ctx statemachine.Context, info AssetPullingInfo) error {
	log.Infof("handle remove: %s", info.CID)
	m.stopAssetCount(info.Hash.String())

	hash := info.Hash.String()
	cid := info.CID

	cInfos, err := m.LoadAssetReplicas(hash)
	if err != nil {
		return xerrors.Errorf("RemoveAsset %s LoadAssetReplicas err:%s", info.CID, err.Error())
	}

	err = m.DeleteAssetRecord(hash, m.nodeMgr.ServerID, &types.AssetEventInfo{Hash: hash, Event: types.AssetEventRemove, Requester: info.Requester}, Remove.String())
	if err != nil {
		return xerrors.Errorf("RemoveAsset %s DeleteAssetRecord err: %s", hash, err.Error())
	}

	for _, cInfo := range cInfos {
		// asset view
		err = m.removeAssetFromView(cInfo.NodeID, cid)
		if err != nil {
			log.Errorf("RemoveAsset %s removeAssetFromView err:%s", cInfo.NodeID, err.Error())
		}

		go m.requestAssetDelete(cInfo.NodeID, cid)
	}

	return nil
}
