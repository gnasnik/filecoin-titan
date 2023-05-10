package config

// // NOTE: ONLY PUT STRUCT DEFINITIONS IN THIS FILE
// //
// // After making edits here, run 'make cfgdoc-gen' (or 'make gen')

// EdgeCfg edge node config
type EdgeCfg struct {
	// host address and port the edge node api will listen on
	ListenAddress string
	// used when 'ListenAddress' is unspecified. must be a valid duration recognized by golang's time.ParseDuration function
	Timeout string
	// node id
	NodeID string
	// area id
	AreaID string
	// used auth when connect to scheduler
	Secret string
	// metadata path
	MetadataPath string
	// assets path
	AssetsPaths []string
	// upload file bandwidth, unit is B/s
	BandwidthUp int64
	// download file bandwidth, unit is B/s
	BandwidthDown int64
	// if true, get scheduler url from locator
	Locator bool
	// InsecureSkipVerify skip tls verify
	InsecureSkipVerify bool
	// used for http3 server
	// be used if InsecureSkipVerify is true
	CertificatePath string
	// used for http3 server
	// be used if InsecureSkipVerify is true
	PrivateKeyPath string
	// self sign certificate, use for client
	CaCertificatePath string
	// PullBlockTimeout get block timeout
	PullBlockTimeout int
	// PullBlockRetry retry when get block failed
	PullBlockRetry int
	// PullBlockParallel the number of goroutine to pull block
	PullBlockParallel int
	TCPSrvAddr        string
	IPFSAPIURL        string
	// seconds
	ValidateDuration int
}

// CandidateCfg candidate node config
type CandidateCfg struct {
	EdgeCfg
}

// LocatorCfg locator config
type LocatorCfg struct {
	// host address and port the edge node api will listen on
	ListenAddress string
	// used when 'ListenAddress' is unspecified. must be a valid duration recognized by golang's time.ParseDuration function
	Timeout string
	// geodb path
	GeoDBPath string
	// InsecureSkipVerify skip tls verify
	InsecureSkipVerify bool
	// used for http3 server
	// be used if InsecureSkipVerify is false
	CertificatePath string
	// used for http3 server
	// be used if InsecureSkipVerify is false
	PrivateKeyPath string
	// self sign certificate, use for client
	CaCertificatePath string
	// etcd server addresses
	EtcdAddresses []string
	DefaultAreaID string
}

// SchedulerCfg scheduler config
type SchedulerCfg struct {
	// host external address and port
	ExternalURL string
	// host address and port the edge node api will listen on
	ListenAddress string
	// database address
	DatabaseAddress string
	// area id
	AreaID string
	// InsecureSkipVerify skip tls verify
	InsecureSkipVerify bool
	// used for http3 server
	// be used if InsecureSkipVerify is true
	CertificatePath string
	// used for http3 server
	// be used if InsecureSkipVerify is true
	PrivateKeyPath string
	// self sign certificate, use for client
	CaCertificatePath string
	// config to enabled node validation, default: true
	EnableValidation bool
	// etcd server addresses
	EtcdAddresses []string
	// Number of candidate node replicas (does not contain 'seed')
	CandidateReplicas int
	// Proportion of validator in candidate nodes (0 ~ 1)
	ValidatorRatio float64
	// The base downstream bandwidth per validator window (unit : MiB)
	ValidatorBaseBwDn int
	// Increased profit after node validation passes
	ValidationProfit float64
	// Node score level scale
	// The key of map is the rank name, and the value of map is a int array containing two elements,
	// the first element of which is the minimum value of score,
	// and the second element is the maximum value of score. (scores out of 100)
	NodeScoreLevel map[string][]int
	// Node level weight
	// The key of the map is the name of the level, and the value of the map is an int,
	// indicating how many select weight this level can get (the more select weight, the greater the probability of the node being selected)
	LevelSelectWeight map[string]int
}
