package message

import (
	"github.com/julienbreux/rabdis/pkg/rabbitmq/message/body"
	"github.com/streadway/amqp"
)

// OnMessageHandler represents a handler when RabbitMQ receive a message
type OnMessageHandler func(Message) error

// Message represents the message interface
type Message interface {
	Headers() *map[string]string
	Body() body.Body

	Ack() error
	Nack(requeue bool) error
}

// message represents a message
type message struct {
	headers  *map[string]string
	body     body.Body
	delivery amqp.Delivery
}

// NewFromDelivery creates a new message instance from delivery
func NewFromDelivery(delivery amqp.Delivery) Message {
	headers := make(map[string]string)
	for k, v := range delivery.Headers {
		headers[k], _ = v.(string)
	}

	return &message{
		headers:  &headers,
		body:     body.New(delivery.Body),
		delivery: delivery,
	}
}

// Headers returns the headers of message
func (m *message) Headers() *map[string]string {
	return m.headers
}

// Body returns the body of message
func (m *message) Body() body.Body {
	return m.body
}

// Ack returns the acknowledge
func (m *message) Ack() error {
	return m.delivery.Ack(false)
}

// Nack returns the non acknowledge
func (m *message) Nack(requeue bool) error {
	return m.delivery.Nack(false, requeue)
}
