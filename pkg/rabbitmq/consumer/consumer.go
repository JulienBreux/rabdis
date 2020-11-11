package consumer

import (
	"github.com/julienbreux/rabdis/pkg/rabbitmq/channel"
	"github.com/julienbreux/rabdis/pkg/rabbitmq/message"
	"github.com/streadway/amqp"
)

const autoAck = false

// Consumer represents the consumer interface
type Consumer interface {
	Start(c *amqp.Connection) error
}

// consumer represents a consumer
type consumer struct {
	o *Options

	h message.OnMessageHandler
}

// New creates a new consumer instance
func New(h message.OnMessageHandler, opts ...Option) (Consumer, error) {
	o, err := newOptions(opts...)
	if err != nil {
		return nil, err
	}

	r := &consumer{
		o: o,

		h: h,
	}

	return r, nil
}

// Start starts the consumer
func (c *consumer) Start(conn *amqp.Connection) error {
	// Create a channel
	ch, err := channel.New(conn)
	if err != nil {
		// TODO: Specific error
		return err
	}

	// Declare exchange
	if err := ch.Exchange(c.o.Exchange); err != nil {
		// TODO: Specific error
		return err
	}

	// Declare queue
	if err := ch.Queue(c.o.Queue); err != nil {
		// TODO: Specific error
		return err
	}

	// Bind queue
	if err := ch.Bind(c.o.Bind); err != nil {
		// TODO: Specific error
		return err
	}

	// Consume
	deliveryChan, err := ch.Consume(c.o.Queue, autoAck)
	if err != nil {
		// TODO: Specific error
		return err
	}

	return c.consume(deliveryChan)
}

func (c *consumer) consume(deliveryChan <-chan amqp.Delivery) error {
	go func() {
		for d := range deliveryChan {
			go func(d amqp.Delivery) {
				_ = c.h(message.NewFromDelivery(d))
			}(d)
		}
	}()
	return nil
}
