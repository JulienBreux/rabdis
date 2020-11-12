package metrics

import (
	"github.com/julienbreux/rabdis/pkg/logger"
	"github.com/kelseyhightower/envconfig"
)

// ConfigPrefix defines the configuration prefix
var ConfigPrefix = "METRICS"

// Option represents a Metrics option
type Option func(*Options)

// Options represents the Metrics options
type Options struct {
	Logger logger.Logger `ignored:"true" json:"-"`

	Port  int    `default:"9090" envconfig:"PORT"`
	Route string `default:"/metrics" envconfig:"ROUTE"`
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

// Port returns port option
func Port(port int) Option {
	return func(o *Options) {
		o.Port = port
	}
}

// Route returns route option
func Route(route string) Option {
	return func(o *Options) {
		o.Route = route
	}
}
