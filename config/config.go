package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	SpotInstrumentService SpotInstrumentService `validate:"required"`
	GRPCServer            GRPCServerConfig      `validate:"required"`
	GRPCClient            GRPCClientConfig      `validate:"required"`
	Infrastructure        InfrastructureConfig  `validate:"required"`
}

type SpotInstrumentService struct {
	LogLevel string `env:"SPOT_INSTRUMENT_SERVICE_LOG_LEVEL" validate:"required,oneof=debug info warn error"`
}

type GRPCServerConfig struct {
	Address              string        `env:"GRPC_SERVER_ADDRESS" validate:"required"`
	MaxRecvMsgSize       int           `env:"GRPC_SERVER_MAX_RECV_MSG_SIZE" validate:"gte=0"`
	MaxSendMsgSize       int           `env:"GRPC_SERVER_MAX_SEND_MSG_SIZE" validate:"gte=0"`
	EnableReflection     bool          `env:"GRPC_SERVER_ENABLE_REFLECTION" validate:"-"`
	TLSCertFile          string        `env:"GRPC_SERVER_TLS_CERT_FILE" validate:"omitempty,file"`
	TLSKeyFile           string        `env:"GRPC_SERVER_TLS_KEY_FILE" validate:"omitempty,file"`
	ReadTimeout          time.Duration `env:"GRPC_SERVER_READ_TIMEOUT" validate:"gte=0"`
	WriteTimeout         time.Duration `env:"GRPC_SERVER_WRITE_TIMEOUT" validate:"gte=0"`
	EnablePrometheus     bool          `env:"GRPC_SERVER_ENABLE_PROMETHEUS" validate:"-"`
	PrometheusListenAddr string        `env:"GRPC_SERVER_PROMETHEUS_LISTEN_ADDR" validate:"required_with=EnablePrometheus,omitempty"`
}

type GRPCClientConfig struct {
	ConnectTimeout    time.Duration `env:"GRPC_CLIENT_CONNECT_TIMEOUT" validate:"gte=0"`
	MaxBackoffDelay   time.Duration `env:"GRPC_CLIENT_MAX_BACKOFF_DELAY" validate:"gte=0"`
	BaseBackoffDelay  time.Duration `env:"GRPC_CLIENT_BASE_BACKOFF_DELAY" validate:"gte=0"`
	BackoffMultiplier float64       `env:"GRPC_CLIENT_BACKOFF_MULTIPLIER" validate:"gte=1"`
	BackoffJitter     float64       `env:"GRPC_CLIENT_BACKOFF_JITTER" validate:"gte=0"`
}

type InfrastructureConfig struct {
	Prometheus    PrometheusConfig    `validate:"required"`
	Opentelemetry OpentelemetryConfig `validate:"required"`
}

type PrometheusConfig struct {
	Address string `env:"PROMETHEUS_ADDRESS" validate:"required"`
}

type OpentelemetryConfig struct {
	ServiceName string `env:"OPENTELEMETRY_SERVICE_NAME" validate:"required"`
	Host        string `env:"OPENTELEMETRY_AGENT_HOST" validate:"required,hostname|ip"`
	Port        string `env:"OPENTELEMETRY_PORT" validate:"required,numeric"`
	LogSpans    bool   `env:"OPENTELEMETRY_LOG_SPANS" validate:"-"`
	Enabled     bool   `env:"OPENTELEMETRY_ENABLED" validate:"-"`
}

func New() (*Config, error) {
	cfg := &Config{}
	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
