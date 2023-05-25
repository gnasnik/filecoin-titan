package types

type CreateAssetReq struct {
	UserID    string
	AssetID   string
	AssetName string
	AssetSize int64
}

type CreateAssetRsp struct {
	CandidateAddr string
	Token         string
}

type AssetProperty struct {
	CID    string
	Name   string
	Status int
	Size   int64
}
