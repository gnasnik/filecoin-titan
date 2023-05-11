package storage

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-libipfs/blocks"
	"github.com/ipld/go-car/v2/blockstore"
	"golang.org/x/xerrors"
)

// asset save asset file
type asset struct {
	assetsPaths *assetsPaths
	suffix      string
}

// newAsset initializes a new asset instance.
func newAsset(assetsPaths *assetsPaths, suffix string) (*asset, error) {
	return &asset{assetsPaths: assetsPaths, suffix: suffix}, nil
}

// generateAssetName creates a new asset file name.
func (a *asset) generateAssetName(root cid.Cid) string {
	return root.Hash().String() + a.suffix
}

// storeBlocks stores blocks to the file system.
func (a *asset) storeBlocks(ctx context.Context, root cid.Cid, blks []blocks.Block) error {
	baseDir, err := a.assetsPaths.allocatePath(root, blks)
	if err != nil {
		return err
	}

	assetDir := filepath.Join(baseDir, root.Hash().String())
	err = os.MkdirAll(assetDir, 0o755)
	if err != nil {
		return err
	}

	for _, blk := range blks {
		filePath := filepath.Join(assetDir, blk.Cid().Hash().String())
		if err := os.WriteFile(filePath, blk.RawData(), 0o644); err != nil {
			return err
		}
	}

	return nil
}

// storeAsset stores the asset to the file system.
func (a *asset) storeAsset(ctx context.Context, root cid.Cid) error {
	baseDir, err := a.assetsPaths.findPath(root)
	if err != nil {
		return err
	}

	assetDir := filepath.Join(baseDir, root.Hash().String())
	entries, err := os.ReadDir(assetDir)
	if err != nil {
		return err
	}

	name := a.generateAssetName(root)
	path := filepath.Join(baseDir, name)

	rw, err := blockstore.OpenReadWrite(path, []cid.Cid{root})
	if err != nil {
		return err
	}

	for _, entry := range entries {
		data, err := ioutil.ReadFile(filepath.Join(assetDir, entry.Name()))
		if err != nil {
			return err
		}

		blk := blocks.NewBlock(data)
		if err = rw.Put(ctx, blk); err != nil {
			return err
		}
	}

	if err = rw.Finalize(); err != nil {
		return err
	}

	return os.RemoveAll(assetDir)
}

// get returns a ReadSeekCloser for the given asset root.
// The caller must close the reader.
func (a *asset) get(root cid.Cid) (io.ReadSeekCloser, error) {
	baseDir, err := a.assetsPaths.findPath(root)
	if err != nil {
		return nil, err
	}

	// check if put asset complete
	assetDir := filepath.Join(baseDir, root.Hash().String())
	if _, err := os.Stat(assetDir); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
	} else {
		return nil, xerrors.Errorf("putting asset, not ready")
	}

	name := a.generateAssetName(root)
	filePath := filepath.Join(baseDir, name)
	return os.Open(filePath)
}

// exists checks if the asset exists in the file system.
func (a *asset) exists(root cid.Cid) (bool, error) {
	if ok := a.assetsPaths.exists(root); !ok {
		return false, nil
	}

	baseDir, err := a.assetsPaths.findPath(root)
	if err != nil {
		return false, err
	}

	assetDir := filepath.Join(baseDir, root.Hash().String())
	if _, err := os.Stat(assetDir); err != nil {
		if !os.IsNotExist(err) {
			return false, err
		}
	} else {
		return false, nil
	}

	name := a.generateAssetName(root)
	filePath := filepath.Join(baseDir, name)

	_, err = os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// remove deletes the asset from the file system.
func (a *asset) remove(root cid.Cid) error {
	baseDir, err := a.assetsPaths.findPath(root)
	if err != nil {
		return err
	}

	a.assetsPaths.releasePath(root)

	assetDir := filepath.Join(baseDir, root.Hash().String())
	if err := os.RemoveAll(assetDir); err != nil {
		if e, ok := err.(*os.PathError); !ok {
			return err
		} else if e.Err != syscall.ENOENT {
			return err
		}
	}

	name := a.generateAssetName(root)
	path := filepath.Join(baseDir, name)

	// remove file
	return os.Remove(path)
}

// count returns the number of assets in the file system.
func (a *asset) count() (int, error) {
	count := 0
	for _, baseDir := range a.assetsPaths.baseDirs {
		entries, err := os.ReadDir(baseDir)
		if err != nil {
			return 0, err
		}
		count += len(entries)
	}

	return count, nil
}
