package api

import "context"

// Candidate is an interface for candidate node
type Candidate interface {
	Common
	Device
	Validation
	DataSync
	Asset
	WaitQuiet(ctx context.Context) error                                                                             //perm:admin
	GetBlocksWithAssetCID(ctx context.Context, assetCID string, randomSeed int64, randomCount int) ([]string, error) //perm:admin
	// GetExternalAddress retrieves the external address of the caller.
	GetExternalAddress(ctx context.Context) (string, error)                        //perm:default
	CheckNetworkConnectivity(ctx context.Context, network, targetURL string) error //perm:default
}

// ValidationResult node Validation result
type ValidationResult struct {
	Validator string
	CID       string
	// verification canceled due to download
	IsCancel  bool
	NodeID    string
	Bandwidth float64
	// seconds duration
	CostTime  int64
	IsTimeout bool

	Msg string
	// if IsCancel is true, Token is valid
	// use for verify edge providing download
	Token string
	// key is random index
	// values is cid
	Cids []string
	// The number of random for validator
	RandomCount int
}
