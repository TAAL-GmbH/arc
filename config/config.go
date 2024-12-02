package config

import (
	"time"

	"github.com/ordishs/go-bitcoin"
	"go.opentelemetry.io/otel/attribute"
)

const (
	FreeCache = "freecache"
	Redis     = "redis"
)

type ArcConfig struct {
	LogLevel           string              `mapstructure:"logLevel"`
	LogFormat          string              `mapstructure:"logFormat"`
	ProfilerAddr       string              `mapstructure:"profilerAddr"`
	PrometheusEndpoint string              `mapstructure:"prometheusEndpoint"`
	PrometheusAddr     string              `mapstructure:"prometheusAddr"`
	GrpcMessageSize    int                 `mapstructure:"grpcMessageSize"`
	Network            string              `mapstructure:"network"`
	MessageQueue       *MessageQueueConfig `mapstructure:"messageQueue"`
	Tracing            *TracingConfig      `mapstructure:"tracing"`
	PeerRPC            *PeerRPCConfig      `mapstructure:"peerRpc"`
	Broadcasting       *BroadcastingConfig `mapstructure:"broadcasting"`
	Metamorph          *MetamorphConfig    `mapstructure:"metamorph"`
	Blocktx            *BlocktxConfig      `mapstructure:"blocktx"`
	API                *APIConfig          `mapstructure:"api"`
	K8sWatcher         *K8sWatcherConfig   `mapstructure:"k8sWatcher"`
	Callbacker         *CallbackerConfig   `mapstructure:"callbacker"`
	Cache              *CacheConfig        `mapstructure:"cache"`
}

type BroadcastingConfig struct {
	Mode      string      `mapstructure:"mode"`
	Multicast *Mulsticast `mapstructure:"multicast"`
	Unicast   *Unicast    `mapstructure:"unicast"`
}

type Unicast struct {
	Peers []*PeerConfig `mapstructure:"peers"`
}
type Mulsticast struct {
	Ipv6Enabled     bool      `mapstructure:"ipv6Enabled"`
	MulticastGroups []*string `mapstructure:"multicastGroups"`
	Interfaces      []*string `mapstructure:"interfaces"`
}

type PeerConfig struct {
	Host string          `mapstructure:"host"`
	Port *PeerPortConfig `mapstructure:"port"`
}

type MessageQueueConfig struct {
	URL       string                `mapstructure:"url"`
	Streaming MessageQueueStreaming `mapstructure:"streaming"`
}

type MessageQueueStreaming struct {
	Enabled     bool `mapstructure:"enabled"`
	FileStorage bool `mapstructure:"fileStorage"`
}

type TracingConfig struct {
	Enabled            bool              `mapstructure:"enabled"`
	DialAddr           string            `mapstructure:"dialAddr"`
	Sample             int               `mapstructure:"sample"`
	Attributes         map[string]string `mapstructure:"attributes"`
	KeyValueAttributes []attribute.KeyValue
}

func (a *ArcConfig) IsTracingEnabled() bool {
	return a.Tracing != nil && a.Tracing.IsEnabled()
}

func (t *TracingConfig) IsEnabled() bool {
	return t.Enabled && t.DialAddr != ""
}

type PeerRPCConfig struct {
	Password string `mapstructure:"password"`
	User     string `mapstructure:"user"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
}

type PeerPortConfig struct {
	P2P int `mapstructure:"p2p"`
	ZMQ int `mapstructure:"zmq"`
}

type MetamorphConfig struct {
	ListenAddr                              string        `mapstructure:"listenAddr"`
	DialAddr                                string        `mapstructure:"dialAddr"`
	Db                                      *DbConfig     `mapstructure:"db"`
	ProcessorCacheExpiryTime                time.Duration `mapstructure:"processorCacheExpiryTime"`
	UnseenTransactionRebroadcastingInterval time.Duration `mapstructure:"unseenTransactionRebroadcastingInterval"`
	MaxRetries                              int           `mapstructure:"maxRetries"`
	ProcessStatusUpdateInterval             time.Duration `mapstructure:"processStatusUpdateInterval"`
	CheckSeenOnNetworkOlderThan             time.Duration `mapstructure:"checkSeenOnNetworkOlderThan"`
	CheckSeenOnNetworkPeriod                time.Duration `mapstructure:"checkSeenOnNetworkPeriod"`
	MonitorPeers                            bool          `mapstructure:"monitorPeers"`
	CheckUtxos                              bool          `mapstructure:"checkUtxos"`
	Health                                  *HealthConfig `mapstructure:"health"`
	RejectCallbackContaining                []string      `mapstructure:"rejectCallbackContaining"`
	Stats                                   *StatsConfig  `mapstructure:"stats"`
}

type BlocktxConfig struct {
	ListenAddr                    string              `mapstructure:"listenAddr"`
	DialAddr                      string              `mapstructure:"dialAddr"`
	HealthServerDialAddr          string              `mapstructure:"healthServerDialAddr"`
	Db                            *DbConfig           `mapstructure:"db"`
	RecordRetentionDays           int                 `mapstructure:"recordRetentionDays"`
	RequestTxsInterval            time.Duration       `mapstructure:"requestTxsInterval"`
	MonitorPeers                  bool                `mapstructure:"monitorPeers"`
	FillGapsInterval              time.Duration       `mapstructure:"fillGapsInterval"`
	MaxAllowedBlockHeightMismatch int                 `mapstructure:"maxAllowedBlockHeightMismatch"`
	MessageQueue                  *MessageQueueConfig `mapstructure:"mq"`
}

type DbConfig struct {
	Mode     string          `mapstructure:"mode"`
	Postgres *PostgresConfig `mapstructure:"postgres"`
}

type PostgresConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Name         string `mapstructure:"name"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	MaxIdleConns int    `mapstructure:"maxIdleConns"`
	MaxOpenConns int    `mapstructure:"maxOpenConns"`
	SslMode      string `mapstructure:"sslMode"`
}

type CacheConfig struct {
	Engine    string           `mapstructure:"engine"`
	Freecache *FreeCacheConfig `mapstructure:"freecache"`
	Redis     *RedisConfig     `mapstructure:"redis"`
}

type FreeCacheConfig struct {
	Size int `mapstructure:"size"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type HealthConfig struct {
	SeverDialAddr             string `mapstructure:"serverDialAddr"`
	MinimumHealthyConnections int    `mapstructure:"minimumHealthyConnections"`
}

type StatsConfig struct {
	NotSeenTimeLimit  time.Duration `mapstructure:"notSeenTimeLimit"`
	NotFinalTimeLimit time.Duration `mapstructure:"notFinalTimeLimit"`
}

type APIConfig struct {
	Address             string            `mapstructure:"address"`
	WocAPIKey           string            `mapstructure:"wocApiKey"`
	WocMainnet          bool              `mapstructure:"wocMainnet"`
	DefaultPolicy       *bitcoin.Settings `mapstructure:"defaultPolicy"`
	RequestExtendedLogs bool              `mapstructure:"requestExtendedLogs"`
}

type K8sWatcherConfig struct {
	Namespace string `mapstructure:"namespace"`
}

type CallbackerConfig struct {
	ListenAddr              string                      `mapstructure:"listenAddr"`
	DialAddr                string                      `mapstructure:"dialAddr"`
	Health                  *HealthConfig               `mapstructure:"health"`
	Delay                   time.Duration               `mapstructure:"delay"`
	Pause                   time.Duration               `mapstructure:"pause"`
	BatchSendInterval       time.Duration               `mapstructure:"batchSendInterval"`
	Db                      *DbConfig                   `mapstructure:"db"`
	PruneInterval           time.Duration               `mapstructure:"pruneInterval"`
	PruneOlderThan          time.Duration               `mapstructure:"pruneOlderThan"`
	QuarantineCheckInterval time.Duration               `mapstructure:"quarantineCheckInterval"`
	QuarantinePolicy        *CallbackerQuarantinePolicy `mapstructure:"quarantinePolicy"`
}

type CallbackerQuarantinePolicy struct {
	BaseDuration        time.Duration `mapstructure:"baseDuration"`
	PermQuarantineAfter time.Duration `mapstructure:"permQuarantineAfter"`
}
