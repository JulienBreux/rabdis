package logger

import "github.com/kelseyhightower/envconfig"

// ConfigPrefix defines the configuration prefix
var ConfigPrefix = "LOGGER"

// Option represents an option
type Option func(*Options)

// Options represents the logger options
type Options struct {
	Level  string `default:"info" envconfig:"LEVEL"`
	Format string `default:"json" envconfig:"FORMAT"`

	InstName    string `default:"unknown" envconfig:"INST_NAME"`
	InstVersion string `default:"unknown" envconfig:"INST_VERSION"`

	DefaultFields []Field `ignored:"true" json:"-"`
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

	if opt.Level == "" {
		opt.Level = "info"
	}
	if opt.Format == "" {
		opt.Format = "text"
	}

	return &opt, nil
}

// DefaultFields returns default fields option
func DefaultFields(fields ...Field) Option {
	return func(o *Options) {
		o.DefaultFields = fields
	}
}

// Level returns level option
func Level(level string) Option {
	return func(o *Options) {
		o.Level = level
	}
}

// Format returns format option
func Format(format string) Option {
	return func(o *Options) {
		o.Format = format
	}
}

// InstName returns instance name option
func InstName(instName string) Option {
	return func(o *Options) {
		o.InstName = instName
	}
}

// InstVersion returns instance option
func InstVersion(instVersion string) Option {
	return func(o *Options) {
		o.InstVersion = instVersion
	}
}
