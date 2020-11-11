package rabdis

import (
	"github.com/julienbreux/rabdis/pkg/logger"
	"github.com/kelseyhightower/envconfig"
)

// ConfigPrefix defines the configuration prefix
var ConfigPrefix = "RABDIS"

// Option represents a Rabdis option
type Option func(*Options)

// Options represents the Rabdis options
type Options struct {
	Logger logger.Logger `ignored:"true" json:"-"`

	ConfigFile string `default:"rabdis.yaml" envconfig:"CONFIG_FILE"`
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
