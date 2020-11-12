package consumer

import (
	"github.com/julienbreux/rabdis/pkg/rabbitmq/bind"
	"github.com/julienbreux/rabdis/pkg/rabbitmq/exchange"
	"github.com/julienbreux/rabdis/pkg/rabbitmq/queue"
)

// Option represents a Consumer option
type Option func(*Options)

// Options represents Consumer options
type Options struct {
	Exchange exchange.Exchange
	Queue    queue.Queue
	Bind     bind.Bind
}

func newOptions(opts ...Option) (*Options, error) {
	opt := Options{}

	for _, o := range opts {
		o(&opt)
	}

	return &opt, nil
}

// Exchange returns exchange option
func Exchange(opt exchange.Exchange) Option {
	return func(o *Options) {
		o.Exchange = opt
	}
}

// Queue returns queue option
func Queue(opt queue.Queue) Option {
	return func(o *Options) {
		o.Queue = opt
	}
}

// Bind returns bind option
func Bind(opt bind.Bind) Option {
	return func(o *Options) {
		o.Bind = opt
	}
}
