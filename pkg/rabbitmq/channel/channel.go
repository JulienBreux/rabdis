package channel

import (
	"errors"

	"github.com/google/uuid"
	"github.com/julienbreux/rabdis/pkg/rabbitmq/bind"
	"github.com/julienbreux/rabdis/pkg/rabbitmq/exchange"
	"github.com/julienbreux/rabdis/pkg/rabbitmq/queue"
	"github.com/streadway/amqp"
)

// Channel represents a RabbitMQ channel
type Channel struct {
	uuid       string
	connection *amqp.Connection
	channel    *amqp.Channel
}

// New creates a new channel
func New(c *amqp.Connection) (*Channel, error) {
	ch := &Channel{
		uuid:       uuid.New().String(),
		connection: c,
	}
	if err := ch.Connect(); err != nil {
		return nil, err
	}
	return ch, nil
}

// Connect connects to the channel
func (ch *Channel) Connect() error {
	var err error
	ch.channel, err = ch.connection.Channel()
	if err != nil {
		return err
	}
	return nil
}

// Close closes the channel
func (ch *Channel) Close() error {
	if ch.channel == nil {
		return errors.New("Channel is nil")
	}
	return ch.channel.Close()
}

// Exchange declares an exchange in the channel
func (ch *Channel) Exchange(e exchange.Exchange) error {
	return ch.channel.ExchangeDeclare(
		e.Name,
		string(e.Type),
		e.Durable,
		e.AutoDelete,
		e.Internal,
		e.NoWait,
		e.Arguments,
	)
}

// Queue declares a queue in the channel
func (ch *Channel) Queue(q queue.Queue) error {
	_, err := ch.channel.QueueDeclare(
		q.Name,
		q.Durable,
		q.AutoDelete,
		q.Exclusive,
		q.NoWait,
		q.Arguments,
	)
	return err
}

// Bind binds a queue to an exchanger
func (ch *Channel) Bind(b bind.Bind) error {
	return ch.channel.QueueBind(
		b.Queue.Name,
		b.RoutingKey,
		b.Exchange.Name,
		b.NoWait,
		b.Arguments,
	)
}

// Consume consumes a queue in the channel
// TODO: create consumer.Consume struct for properties/options
func (ch *Channel) Consume(q queue.Queue, autoAck bool) (<-chan amqp.Delivery, error) {
	return ch.channel.Consume(
		q.Name,
		ch.uuid,
		autoAck, // autoAck
		false,   // exclusive
		false,   // nolocal
		false,   // nowait
		nil,     // args
	)
}
