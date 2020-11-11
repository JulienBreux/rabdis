package consumer

import (
	"github.com/julienbreux/rabdis/pkg/rabbitmq/bind"
	"github.com/julienbreux/rabdis/pkg/rabbitmq/exchange"
	"github.com/julienbreux/rabdis/pkg/rabbitmq/queue"
)

// Options  represents message options
type Options struct {
	Exchange exchange.Exchange
	Queue    queue.Queue
	Bind     bind.Bind
}

// Option represents a message option
type Option func(*Options)

func newOptions(opts ...Option) (*Options, error) {
	opt := Options{}

	for _, o := range opts {
		o(&opt)
	}

	return &opt, nil
}

// Exchange sets exchange option
func Exchange(opt exchange.Exchange) Option {
	return func(o *Options) {
		o.Exchange = opt
	}
}

// Queue sets queue option
func Queue(opt queue.Queue) Option {
	return func(o *Options) {
		o.Queue = opt
	}
}

// Bind sets bind option
func Bind(opt bind.Bind) Option {
	return func(o *Options) {
		o.Bind = opt
	}
}
