module github.com/Filecoin-Titan/titan

go 1.17

require (
	contrib.go.opencensus.io/exporter/prometheus v0.4.1
	github.com/BurntSushi/toml v1.2.1
	github.com/acarl005/stripansi v0.0.0-20180116102854-5a71ef0e047d
	github.com/docker/go-units v0.5.0
	github.com/fatih/color v1.13.0
	github.com/filecoin-project/go-jsonrpc v0.2.3
	github.com/filecoin-project/go-statemachine v1.0.3
	github.com/filecoin-project/pubsub v1.0.0
	github.com/gabriel-vasile/mimetype v1.4.2
	github.com/gbrlsnchs/jwt/v3 v3.0.1
	github.com/go-sql-driver/mysql v1.6.0
	github.com/google/uuid v1.3.0
	github.com/gorilla/mux v1.8.0
	github.com/icza/backscanner v0.0.0-20210726202459-ac2ffc679f94
	github.com/ipfs/go-block-format v0.1.1
	github.com/ipfs/go-cid v0.4.0
	github.com/ipfs/go-datastore v0.6.0
	github.com/ipfs/go-ds-leveldb v0.5.0
	github.com/ipfs/go-ds-measure v0.2.0
	github.com/ipfs/go-fetcher v1.6.1
	github.com/ipfs/go-fs-lock v0.0.7
	github.com/ipfs/go-ipfs-http-client v0.5.0
	github.com/ipfs/go-ipld-format v0.4.0
	github.com/ipfs/go-ipld-legacy v0.1.1
	github.com/ipfs/go-libipfs v0.7.0
	github.com/ipfs/go-log/v2 v2.5.1
	github.com/ipfs/go-merkledag v0.9.0
	github.com/ipfs/go-metrics-interface v0.0.1
	github.com/ipfs/go-unixfsnode v1.5.2
	github.com/ipfs/interface-go-ipfs-core v0.11.1
	github.com/ipld/go-codec-dagpb v1.5.0
	github.com/ipld/go-ipld-prime v0.20.0
	github.com/jmoiron/sqlx v1.3.5
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/lib/pq v1.10.7
	github.com/mitchellh/go-homedir v1.1.0
	github.com/multiformats/go-base32 v0.1.0
	github.com/multiformats/go-multiaddr v0.8.0
	github.com/multiformats/go-multihash v0.2.1
	github.com/oschwald/geoip2-golang v1.7.0
	github.com/prometheus/client_golang v1.14.0
	github.com/quic-go/quic-go v0.33.0
	github.com/shirou/gopsutil/v3 v3.22.9
	github.com/syndtr/goleveldb v1.0.1-0.20210819022825-2ae1ddf74ef7
	github.com/urfave/cli/v2 v2.25.3
	github.com/whyrusleeping/cbor-gen v0.0.0-20230126041949-52956bd4c9aa
	go.opencensus.io v0.24.0
	go.uber.org/fx v1.19.2
	golang.org/x/sys v0.8.0
	golang.org/x/time v0.0.0-20220922220347-f3bd1da661af
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2
)

require (
	github.com/alecthomas/units v0.0.0-20210927113745-59d0afb8317a
	github.com/ethereum/go-ethereum v1.11.6
	github.com/ipld/go-car/v2 v2.7.0
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/coreos/go-systemd v0.0.0-20190321100706-95778dfbb74e // indirect
	github.com/coreos/pkg v0.0.0-20160727233714-3ac0863d7acf // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/crackcomm/go-gitignore v0.0.0-20170627025303-887ab5e44cc3 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.2.0 // indirect
	github.com/go-kit/log v0.2.1 // indirect
	github.com/go-logfmt/logfmt v0.5.1 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/go-task/slim-sprig v0.0.0-20210107165309-348f09dbbbc0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/golang/snappy v0.0.5-0.20220116011046-fa5810519dcb // indirect
	github.com/google/pprof v0.0.0-20221203041831-ce31453925ec // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/hashicorp/golang-lru v0.6.0
	github.com/ipfs/bbloom v0.0.4 // indirect
	github.com/ipfs/go-ipfs-blockstore v1.2.0 // indirect
	github.com/ipfs/go-ipfs-cmds v0.8.2 // indirect
	github.com/ipfs/go-ipfs-ds-help v1.1.0 // indirect
	github.com/ipfs/go-ipfs-exchange-interface v0.2.0 // indirect
	github.com/ipfs/go-ipfs-files v0.3.0
	github.com/ipfs/go-ipfs-util v0.0.2 // indirect
	github.com/ipfs/go-log v1.0.5 // indirect
	github.com/ipfs/go-path v0.3.1
	github.com/ipfs/go-unixfs v0.4.3
	github.com/ipfs/go-verifcid v0.0.2 // indirect
	github.com/jbenet/goprocess v0.1.4 // indirect
	github.com/lufia/plan9stats v0.0.0-20211012122336-39d0f177ccd0 // indirect
	github.com/magefile/mage v1.9.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/minio/sha256-simd v1.0.0 // indirect
	github.com/mr-tron/base58 v1.2.0 // indirect
	github.com/multiformats/go-base36 v0.2.0 // indirect
	github.com/multiformats/go-multicodec v0.8.0
	github.com/multiformats/go-varint v0.0.7 // indirect
	github.com/onsi/ginkgo/v2 v2.5.1 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/oschwald/maxminddb-golang v1.9.0 // indirect
	github.com/petar/GoLLRB v0.0.0-20210522233825-ae3b015fd3e9 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/polydawn/refmt v0.89.0 // indirect
	github.com/power-devops/perfstat v0.0.0-20210106213030-5aafc221ea8c // indirect
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/prometheus/common v0.39.0 // indirect
	github.com/prometheus/procfs v0.9.0 // indirect
	github.com/prometheus/statsd_exporter v0.21.0 // indirect
	github.com/quic-go/qpack v0.4.0 // indirect
	github.com/quic-go/qtls-go1-19 v0.2.1 // indirect
	github.com/quic-go/qtls-go1-20 v0.1.1 // indirect
	github.com/rs/cors v1.7.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	github.com/tklauser/go-sysconf v0.3.10 // indirect
	github.com/tklauser/numcpus v0.4.0 // indirect
	github.com/whyrusleeping/cbor v0.0.0-20171005072247-63513f603b11 // indirect
	github.com/xrash/smetrics v0.0.0-20201216005158-039620a65673 // indirect
	github.com/yusufpapurcu/wmi v1.2.2 // indirect
	go.opentelemetry.io/otel v1.12.0 // indirect
	go.opentelemetry.io/otel/trace v1.12.0 // indirect
	go.uber.org/dig v1.16.1 // indirect
	go.uber.org/zap v1.24.0 // indirect
	go4.org v0.0.0-20200411211856-f5505b9728dd // indirect
	golang.org/x/crypto v0.9.0 // indirect
	golang.org/x/exp v0.0.0-20230206171751-46f607a40771 // indirect
	golang.org/x/mod v0.9.0 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	google.golang.org/genproto v0.0.0-20210917145530-b395a37504d4 // indirect
	google.golang.org/grpc v1.45.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

require (
	github.com/StackExchange/wmi v1.2.1 // indirect
	github.com/btcsuite/btcd/btcec/v2 v2.3.2 // indirect
	github.com/deckarep/golang-set/v2 v2.1.0 // indirect
	github.com/filecoin-project/go-cbor-util v0.0.1 // indirect
	github.com/filecoin-project/go-statestore v0.2.0 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/holiman/uint256 v1.2.2 // indirect
	github.com/ipfs/go-bitfield v1.1.0 // indirect
	github.com/ipfs/go-ipld-cbor v0.0.6 // indirect
	github.com/jonboulle/clockwork v0.2.2 // indirect
	github.com/libp2p/go-libp2p v0.26.2 // indirect
	github.com/multiformats/go-multistream v0.4.1 // indirect
	github.com/shirou/gopsutil v3.21.4-0.20210419000835-c7a38de76ee5+incompatible // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.9.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
	gopkg.in/natefinch/npipe.v2 v2.0.0-20160621034901-c1b8fa8bdcce // indirect
	lukechampine.com/blake3 v1.1.7 // indirect
)

require (
	github.com/ipfs/go-blockservice v0.5.0
	github.com/klauspost/cpuid/v2 v2.2.3 // indirect
	github.com/libp2p/go-buffer-pool v0.1.0 // indirect
	github.com/multiformats/go-multibase v0.1.1 // indirect
	go.etcd.io/etcd v0.0.0-20191023171146-3cf2f69b5738
	golang.org/x/tools v0.7.0 // indirect
)

replace github.com/filecoin-project/go-jsonrpc v0.2.3 => github.com/zscboy/go-jsonrpc v0.2.3

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

replace github.com/coreos/bbolt v1.3.4 => go.etcd.io/bbolt v1.3.4
