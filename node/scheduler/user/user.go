package user

import (
	"context"

	"github.com/Filecoin-Titan/titan/api/types"
)

type User struct {
	ID string
}

// AllocateStorage allocates storage space.
func (u *User) AllocateStorage(ctx context.Context, size int64) error {
	return nil
}

// CreateAPIKey creates a key for the client API.
func (u *User) CreateAPIKey(ctx context.Context, keyName string) (string, error) {
	return "", nil
}

// CreateAsset creates an asset with car CID, car name, and car size.
func (u *User) CreateAsset(ctx context.Context, assetID, assetName string, assetSize int64) (*types.CreateAssetRsp, error) {
	return nil, nil
}

// ListAssets lists the assets of the user.
func (u *User) ListAssets(ctx context.Context) ([]*types.AssetProperty, error) {
	return nil, nil
}

// DeleteAssets deletes the assets of the user.
func (u *User) DeleteAssets(ctx context.Context, assetCIDs []string) error {
	return nil
}

// ShareAssets shares the assets of the user.
func (u *User) ShareAssets(ctx context.Context, assetCIDs []string) ([]string, error) {
	return nil, nil
}
