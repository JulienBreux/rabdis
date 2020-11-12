package redis

import (
	"time"

	"github.com/julienbreux/rabdis/pkg/logger"
	"github.com/kelseyhightower/envconfig"
)

// ConfigPrefix define configuration prefix
var ConfigPrefix = "REDIS"

// Option defines a Redis option
type Option func(*Options)

// Options represents Redis options
type Options struct {
	Logger logger.Logger `ignored:"true" json:"-"`

	Host     string `default:"0.0.0.0" envconfig:"HOST"`
	Port     int    `default:"6379" envconfig:"PORT"`
	Password string `default:"" envconfig:"PASSWORD"`
	Database int    `default:"0" envconfig:"DATABASE"`

	PingDelay time.Duration `default:"2s" envconfig:"PING_DELAY"`

	KeyPrefix string `default:"" envconfig:"KEY_PREFIX"`
}

func newOptions(opts ...Option) (*Options, error) {
	opt := Options{}
	err := envconfig.Process(ConfigPrefix, &opt)
	if err != nil {
		return nil, err
	}

	for _, o := range opts {
		o(&opt)
	}

	// Set a default logger
	if opt.Logger == nil {
		if opt.Logger == nil {
			l, err := logger.New()
			if err != nil {
				return nil, err
			}
			opt.Logger = l
		}
	}

	return &opt, nil
}

// Logger returns logger option
func Logger(logger logger.Logger) Option {
	return func(o *Options) {
		o.Logger = logger
	}
}

// Host returns host option
func Host(host string) Option {
	return func(o *Options) {
		o.Host = host
	}
}

// Port returns port option
func Port(port int) Option {
	return func(o *Options) {
		o.Port = port
	}
}

// Password returns password option
func Password(password string) Option {
	return func(o *Options) {
		o.Password = password
	}
}

// Database returns database option
func Database(database int) Option {
	return func(o *Options) {
		o.Database = database
	}
}

// PingDelay returns ping delay option
func PingDelay(pingDelay time.Duration) Option {
	return func(o *Options) {
		o.PingDelay = pingDelay
	}
}

// KeyPrefix returns key prefix option
func KeyPrefix(keyPrefix string) Option {
	return func(o *Options) {
		o.KeyPrefix = keyPrefix
	}
}
