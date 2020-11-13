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
	AutoAck  bool
}

func newOptions(opts ...Option) *Options {
	opt := Options{
		AutoAck: false,
	}

	for _, o := range opts {
		o(&opt)
	}

	return &opt
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

// AutoAck returns auto acknowledge option
func AutoAck(opt bool) Option {
	return func(o *Options) {
		o.AutoAck = opt
	}
}
