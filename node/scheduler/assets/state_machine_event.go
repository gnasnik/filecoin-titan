package assets

import (
	"golang.org/x/xerrors"
)

type mutator interface {
	apply(state *AssetPullingInfo)
}

// globalMutator is an event which can apply in every state
type globalMutator interface {
	// applyGlobal applies the event to the state. If if returns true,
	//  event processing should be interrupted
	applyGlobal(state *AssetPullingInfo) bool
}

// Ignorable Ignorable
type Ignorable interface {
	Ignore()
}

// Global events

// PullAssetRestart restarts asset pulling
type PullAssetRestart struct{}

func (evt PullAssetRestart) applyGlobal(state *AssetPullingInfo) bool {
	state.RetryCount = 0
	return false
}

// PullAssetFatalError represents a fatal error in asset pulling
type PullAssetFatalError struct{ error }

// FormatError Format error
func (evt PullAssetFatalError) FormatError(xerrors.Printer) (next error) { return evt.error }

func (evt PullAssetFatalError) applyGlobal(state *AssetPullingInfo) bool {
	log.Errorf("Fatal error on asset %s: %+v", state.CID, evt.error)
	return true
}

// AssetForceState forces an asset state
type AssetForceState struct {
	State AssetState
}

func (evt AssetForceState) applyGlobal(state *AssetPullingInfo) bool {
	state.State = evt.State
	return true
}

// InfoUpdate update asset info
type InfoUpdate struct {
	Size   int64
	Blocks int64
}

func (evt InfoUpdate) applyGlobal(state *AssetPullingInfo) bool {
	if state.State == SeedPulling {
		state.Size = evt.Size
		state.Blocks = evt.Blocks
	}

	return true
}

// PulledResult represents the result of node pulling
type PulledResult struct {
	ResultInfo *NodePulledResult
}

func (evt PulledResult) apply(state *AssetPullingInfo) {
	rInfo := evt.ResultInfo
	if rInfo == nil {
		return
	}

	if state.State == SeedPulling {
		state.Size = rInfo.Size
		state.Blocks = rInfo.BlocksCount
	}
}

// PullRequestSent indicates that a pull request has been sent
type PullRequestSent struct {
	CandidateWaitings int64
	EdgeWaitings      int64
}

func (evt PullRequestSent) apply(state *AssetPullingInfo) {
	state.CandidateWaitings = evt.CandidateWaitings
	state.EdgeWaitings = evt.EdgeWaitings
}

// Normal path

// AssetStartPulls start asset pulls
type AssetStartPulls struct{}

func (evt AssetStartPulls) apply(state *AssetPullingInfo) {
}

// AssetRePull re-pull the asset
type AssetRePull struct{}

func (evt AssetRePull) apply(state *AssetPullingInfo) {
	state.RetryCount++
}

// PullSucceed indicates that a node has successfully pulled an asset
type PullSucceed struct{}

func (evt PullSucceed) apply(state *AssetPullingInfo) {
	state.RetryCount = 0

	// Check to node offline while replenishing the temporary replicas
	// After these temporary replicas are pulled, the count should be deleted
	if state.State == EdgesPulling {
		state.ReplenishReplicas = 0
	}
}

// SkipStep skips the current step
type SkipStep struct{}

func (evt SkipStep) apply(state *AssetPullingInfo) {}

// PullFailed indicates that a node has failed to pull an asset
type PullFailed struct{ error }

// FormatError Format error
func (evt PullFailed) FormatError(xerrors.Printer) (next error) { return evt.error }

func (evt PullFailed) apply(state *AssetPullingInfo) {
}

// SelectFailed  indicates that node selection has failed
type SelectFailed struct{ error }

// FormatError Format error
func (evt SelectFailed) FormatError(xerrors.Printer) (next error) { return evt.error }

func (evt SelectFailed) apply(state *AssetPullingInfo) {
}
