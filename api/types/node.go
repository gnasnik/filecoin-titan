package types

import (
	"time"

	"github.com/Filecoin-Titan/titan/node/modules/dtypes"
)

type Base struct {
	ID        uint      `gorm:"primarykey"`
	CreatedAt time.Time `json:"created_at" gorm:"comment:'创建时间';type:timestamp;"`
	UpdatedAt time.Time `json:"updated_at" gorm:"comment:'更新时间';type:timestamp;"`
}

// NodeInfo contains information about a node.
type NodeInfo struct {
	Base
	NodeID       string    `json:"node_id" form:"nodeId" gorm:"column:node_id;comment:;" db:"node_id"`
	SerialNumber string    `json:"serial_number" form:"serialNumber" gorm:"column:serial_number;comment:;"`
	Type         NodeType  `json:"type"`
	ExternalIP   string    `json:"external_ip" form:"externalIp" gorm:"column:external_ip;comment:;"`
	InternalIP   string    `json:"internal_ip" form:"internalIp" gorm:"column:internal_ip;comment:;"`
	IPLocation   string    `json:"ip_location" form:"ipLocation" gorm:"column:ip_location;comment:;"`
	PkgLossRatio float64   `json:"pkg_loss_ratio" form:"pkgLossRatio" gorm:"column:pkg_loss_ratio;comment:;"`
	Latency      float64   `json:"latency" form:"latency" gorm:"column:latency;comment:;"`
	CPUUsage     float64   `json:"cpu_usage" form:"cpuUsage" gorm:"column:cpu_usage;comment:;"`
	MemoryUsage  float64   `json:"memory_usage" form:"memoryUsage" gorm:"column:memory_usage;comment:;"`
	IsOnline     bool      `json:"is_online" form:"isOnline" gorm:"column:is_online;comment:;"`
	FirstTime    time.Time `db:"first_login_time"`

	DiskUsage       float64         `json:"disk_usage" form:"diskUsage" gorm:"column:disk_usage;comment:;" db:"disk_usage"`
	Blocks          int             `json:"blocks" form:"blockCount" gorm:"column:blocks;comment:;" db:"blocks"`
	BandwidthUp     int64           `json:"bandwidth_up" db:"bandwidth_up"`     // B
	BandwidthDown   int64           `json:"bandwidth_down" db:"bandwidth_down"` // B
	NATType         string          `json:"nat_type" form:"natType" gorm:"column:nat_type;comment:;" db:"nat_type"`
	DiskSpace       float64         `json:"disk_space" form:"diskSpace" gorm:"column:disk_space;comment:;" db:"disk_space"`
	SystemVersion   string          `json:"system_version" form:"systemVersion" gorm:"column:system_version;comment:;" db:"system_version"`
	DiskType        string          `json:"disk_type" form:"diskType" gorm:"column:disk_type;comment:;" db:"disk_type"`
	IoSystem        string          `json:"io_system" form:"ioSystem" gorm:"column:io_system;comment:;" db:"io_system"`
	Latitude        float64         `json:"latitude" db:"latitude"`
	Longitude       float64         `json:"longitude" db:"longitude"`
	NodeName        string          `json:"node_name" form:"nodeName" gorm:"column:node_name;comment:;" db:"node_name"`
	Memory          float64         `json:"memory" form:"memory" gorm:"column:memory;comment:;" db:"memory"`
	CPUCores        int             `json:"cpu_cores" form:"cpuCores" gorm:"column:cpu_cores;comment:;" db:"cpu_cores"`
	ProductType     string          `json:"product_type" form:"productType" gorm:"column:product_type;comment:;" db:"product_type"`
	MacLocation     string          `json:"mac_location" form:"macLocation" gorm:"column:mac_location;comment:;" db:"mac_location"`
	OnlineDuration  int             `json:"online_duration" form:"onlineDuration" db:"online_duration"` // unit:Minute
	Profit          float64         `json:"profit" db:"profit"`
	DownloadTraffic int64           `json:"download_traffic" db:"download_traffic"` // B
	UploadTraffic   int64           `json:"upload_traffic" db:"upload_traffic"`     // B
	DownloadBlocks  int             `json:"download_blocks" form:"downloadCount" gorm:"column:download_blocks;comment:;" db:"download_blocks"`
	PortMapping     string          `db:"port_mapping"`
	LastSeen        time.Time       `db:"last_seen"`
	SchedulerID     dtypes.ServerID `db:"scheduler_sid"`
}

// NodeType node type
type NodeType int

const (
	NodeUnknown NodeType = iota

	NodeEdge
	NodeCandidate
	NodeValidator
	NodeScheduler
	NodeLocator
	NodeUpdater
)

func (n NodeType) String() string {
	switch n {
	case NodeEdge:
		return "edge"
	case NodeCandidate:
		return "candidate"
	case NodeScheduler:
		return "scheduler"
	case NodeValidator:
		return "validator"
	case NodeLocator:
		return "locator"
	}

	return ""
}

// RunningNodeType represents the type of the running node.
var RunningNodeType NodeType

// DownloadHistory represents the record of a node download
type DownloadHistory struct {
	ID           string    `json:"-"`
	NodeID       string    `json:"node_id" db:"node_id"`
	BlockCID     string    `json:"block_cid" db:"block_cid"`
	AssetCID     string    `json:"asset_cid" db:"asset_cid"`
	BlockSize    int       `json:"block_size" db:"block_size"`
	Speed        int64     `json:"speed" db:"speed"`
	Reward       int64     `json:"reward" db:"reward"`
	Status       int       `json:"status" db:"status"`
	FailedReason string    `json:"failed_reason" db:"failed_reason"`
	ClientIP     string    `json:"client_ip" db:"client_ip"`
	CreatedTime  time.Time `json:"created_time" db:"created_time"`
	CompleteTime time.Time `json:"complete_time" db:"complete_time"`
}

// EdgeDownloadInfo represents download information for an edge node
type EdgeDownloadInfo struct {
	Address string
	Tk      *Token
	NodeID  string
	NatType string
}

// EdgeDownloadInfoList represents a list of EdgeDownloadInfo structures along with
// scheduler URL and key
type EdgeDownloadInfoList struct {
	Infos        []*EdgeDownloadInfo
	SchedulerURL string
	SchedulerKey string
}

// CandidateDownloadInfo represents download information for a candidate
type CandidateDownloadInfo struct {
	NodeID  string
	Address string
	Tk      *Token
}

// NodeReplicaStatus represents the status of a node cache
type NodeReplicaStatus struct {
	Hash   string        `db:"hash"`
	Status ReplicaStatus `db:"status"`
}

// NodeReplicaRsp represents the replicas of a node asset
type NodeReplicaRsp struct {
	Replica    []*NodeReplicaStatus
	TotalCount int
}

// NatType represents the type of NAT of a node
type NatType int

const (
	// NatTypeUnknown Unknown NAT type
	NatTypeUnknown NatType = iota
	// NatTypeNo not  nat
	NatTypeNo
	// NatTypeSymmetric Symmetric NAT
	NatTypeSymmetric
	// NatTypeFullCone Full cone NAT
	NatTypeFullCone
	// NatTypeRestricted Restricted NAT
	NatTypeRestricted
	// NatTypePortRestricted Port-restricted NAT
	NatTypePortRestricted
)

func (n NatType) String() string {
	switch n {
	case NatTypeNo:
		return "NoNAT"
	case NatTypeSymmetric:
		return "SymmetricNAT"
	case NatTypeFullCone:
		return "FullConeNAT"
	case NatTypeRestricted:
		return "RestrictedNAT"
	case NatTypePortRestricted:
		return "PortRestrictedNAT"
	}

	return "UnknowNAT"
}

func (n NatType) FromString(natType string) NatType {
	switch natType {
	case "NoNat":
		return NatTypeNo
	case "SymmetricNAT":
		return NatTypeSymmetric
	case "FullConeNAT":
		return NatTypeFullCone
	case "RestrictedNAT":
		return NatTypeRestricted
	case "PortRestrictedNAT":
		return NatTypePortRestricted
	}
	return NatTypeUnknown
}

// ListNodesRsp list node rsp
type ListNodesRsp struct {
	Data  []NodeInfo `json:"data"`
	Total int64      `json:"total"`
}

// ListDownloadRecordRsp download record rsp
type ListDownloadRecordRsp struct {
	Data  []DownloadHistory `json:"data"`
	Total int64             `json:"total"`
}

// ListValidationResultRsp list validated result
type ListValidationResultRsp struct {
	Total                 int                    `json:"total"`
	ValidationResultInfos []ValidationResultInfo `json:"validation_result_infos"`
}

// ListWorkloadRecordRsp list workload result
type ListWorkloadRecordRsp struct {
	Total               int               `json:"total"`
	WorkloadRecordInfos []*WorkloadRecord `json:"workload_result_infos"`
}

// ValidationResultInfo validator result info
type ValidationResultInfo struct {
	ID               int              `db:"id"`
	RoundID          string           `db:"round_id"`
	NodeID           string           `db:"node_id"`
	Cid              string           `db:"cid"`
	ValidatorID      string           `db:"validator_id"`
	BlockNumber      int64            `db:"block_number"` // number of blocks verified
	Status           ValidationStatus `db:"status"`
	Duration         int64            `db:"duration"` // validator duration, microsecond
	Bandwidth        float64          `db:"bandwidth"`
	StartTime        time.Time        `db:"start_time"`
	EndTime          time.Time        `db:"end_time"`
	Profit           float64          `db:"profit"`
	CalculatedProfit bool             `db:"calculated_profit"`
	TokenID          string           `db:"token_id"`
	FileSaved        bool             `db:"file_saved"`

	UploadTraffic float64 `db:"upload_traffic"`
}

// ValidationStatus Validation Status
type ValidationStatus int

const (
	// ValidationStatusCreate  is the initial validation status when the validation process starts.
	ValidationStatusCreate ValidationStatus = iota
	// ValidationStatusSuccess is the validation status when the validation is success.
	ValidationStatusSuccess
	// ValidationStatusCancel is the validation status when the validation is canceled.
	ValidationStatusCancel

	// Node error

	// ValidationStatusNodeTimeOut is the validation status when the node times out.
	ValidationStatusNodeTimeOut
	// ValidationStatusValidateFail is the validation status when the validation fail.
	ValidationStatusValidateFail

	// Validator error

	// ValidationStatusValidatorTimeOut is the validation status when the validator times out.
	ValidationStatusValidatorTimeOut
	// ValidationStatusGetValidatorBlockErr is the validation status when there is an error getting the blocks from validator.
	ValidationStatusGetValidatorBlockErr
	// ValidationStatusValidatorMismatch is the validation status when the validator mismatches.
	ValidationStatusValidatorMismatch

	// Server error

	// ValidationStatusLoadDBErr is the validation status when there is an error loading the database.
	ValidationStatusLoadDBErr
	// ValidationStatusCIDToHashErr is the validation status when there is an error converting a CID to a hash.
	ValidationStatusCIDToHashErr
)

// TokenPayload payload of token
type TokenPayload struct {
	ID         string    `db:"token_id"`
	NodeID     string    `db:"node_id"`
	AssetCID   string    `db:"asset_id"`
	ClientID   string    `db:"client_id"`
	LimitRate  int64     `db:"limit_rate"`
	CreateTime time.Time `db:"create_time"`
	Expiration time.Time `db:"expiration"`
}

// Token access download asset
type Token struct {
	ID string
	// CipherText encrypted TokenPayload by public key
	CipherText string
	// Sign signs CipherText by scheduler private key
	Sign string
}

type Workload struct {
	DownloadSpeed int64
	DownloadSize  int64
	StartTime     int64
	EndTime       int64
	BlockCount    int64
}

// WorkloadStatus Workload Status
type WorkloadStatus int

const (
	// WorkloadStatusCreate is the initial workload status when the workload process starts.
	WorkloadStatusCreate WorkloadStatus = iota
	// WorkloadStatusSucceeded is the workload status when the workload is succeeded.
	WorkloadStatusSucceeded
	// WorkloadStatusFailed is the workload status when the workload is failed.
	WorkloadStatusFailed
)

type WorkloadReport struct {
	TokenID  string
	ClientID string
	NodeID   string
	Workload *Workload
}

// WorkloadReportRecord use to store workloadReport
type WorkloadRecord struct {
	TokenPayload
	Status         WorkloadStatus `db:"status"`
	ClientWorkload []byte         `db:"client_workload"`
	NodeWorkload   []byte         `db:"node_workload"`
}

type NodeWorkloadReport struct {
	// CipherText encrypted []*WorkloadReport by scheduler public key
	CipherText []byte
	// Sign signs CipherText by node private key
	Sign []byte
}

type NatPunchReq struct {
	Tk      *Token
	NodeID  string
	Timeout time.Duration
}

type ConnectOptions struct {
	Token         string
	TcpServerPort int
}
