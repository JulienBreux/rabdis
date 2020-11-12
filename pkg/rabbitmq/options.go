package rabbitmq

import (
	"github.com/julienbreux/rabdis/pkg/logger"
	"github.com/kelseyhightower/envconfig"
)

// ConfigPrefix defines the configuration prefix
var ConfigPrefix = "RABBITMQ"

// Option represents a RabbitMQ option
type Option func(*Options)

// Options represents the RabbitMQ options
type Options struct {
	Logger logger.Logger `ignored:"true" json:"-"`

	Host        string `default:"0.0.0.0" envconfig:"HOST"`
	Port        int    `default:"5672" envconfig:"PORT"`
	Username    string `default:"guest" envconfig:"USERNAME"`
	Password    string `default:"guest" envconfig:"PASSWORD"`
	VirtualHost string `default:"/" envconfig:"VHOST"`
	ConnTimeout int    `default:"10" envconfig:"CONN_TIMEOUT"`

	InstName    string `default:"unknown" envconfig:"INST_NAME"`
	InstVersion string `default:"unknown" envconfig:"INST_VERSION"`
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

// Username returns username option
func Username(username string) Option {
	return func(o *Options) {
		o.Username = username
	}
}

// Password returns password option
func Password(password string) Option {
	return func(o *Options) {
		o.Password = password
	}
}

// VirtualHost return virtual host option
func VirtualHost(virtualHost string) Option {
	return func(o *Options) {
		o.VirtualHost = virtualHost
	}
}

// ConnTimeout returns connection timeout option
func ConnTimeout(connTimeout int) Option {
	return func(o *Options) {
		o.ConnTimeout = connTimeout
	}
}

// InstName returns instance name option
func InstName(instName string) Option {
	return func(o *Options) {
		o.InstName = instName
	}
}

// InstVersion returns instance version option
func InstVersion(instVersion string) Option {
	return func(o *Options) {
		o.InstVersion = instVersion
	}
}
