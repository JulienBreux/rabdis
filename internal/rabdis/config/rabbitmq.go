package config

import (
	"errors"

	"github.com/julienbreux/rabdis/pkg/rabbitmq/bind"
	"github.com/julienbreux/rabdis/pkg/rabbitmq/exchange"
	"github.com/julienbreux/rabdis/pkg/rabbitmq/queue"
)

// RabbitMQ represents a RabbitMQ section
type RabbitMQ struct {
	ExchangeName string `yaml:"exchangeName"`
	RoutingKey   string `yaml:"routingKey"`
	QueueName    string `yaml:"queueName"`

	Exchange exchange.Exchange `yaml:"exchange"`
	Bind     bind.Bind         `yaml:"bind"`
	Queue    queue.Queue       `yaml:"queue"`
}

// UnmarshalYAML returns an unmarshal YAML implementation
func (r *RabbitMQ) UnmarshalYAML(u func(interface{}) error) error {
	type rawRabbitMQ RabbitMQ
	raw := rawRabbitMQ{}
	if err := u(&raw); err != nil {
		return err
	}

	// Exchange name is missing
	if raw.ExchangeName == "" && raw.Exchange.Name == "" {
		return errors.New("rules.*.rabbitmq.exchangeName or rules.*.rabbitmq.exchange.name  is required")
	}

	// Bind name is missing
	if raw.RoutingKey == "" && raw.Bind.RoutingKey == "" {
		return errors.New("rules.*.rabbitmq.routingKey or rules.*.rabbitmq.bind.routingKey  is required")
	}

	// Queue name is missing
	if raw.QueueName == "" && raw.Queue.Name == "" {
		return errors.New("rules.*.rabbitmq.queueName or rules.*.rabbitmq.queue.name is required")
	}

	*r = RabbitMQ(raw)

	return nil
}
