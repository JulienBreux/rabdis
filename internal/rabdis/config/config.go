package config

// TODO: Split package

import (
	"io/ioutil"

	"github.com/julienbreux/rabdis/pkg/logger"
	"gopkg.in/yaml.v3"
)

// Config represents the config file structured
type Config struct {
	Version string `yaml:"version"`

	Rules []Rule `yaml:"rules"`
}

// FromFile converts the configuration file to struct
// TODO: Create clean errors on config (most readable)
func FromFile(file string, log logger.Logger) (c *Config, err error) {
	f, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	c, err = Unmarshal(f, log)
	if err != nil {
		return nil, err
	}

	log.Debug(
		"rabdis: configuration file content",
		logger.F("content", c),
	)

	return
}

// Unmarshal decodes configuration
func Unmarshal(in []byte, log logger.Logger) (c *Config, err error) {
	if err := yaml.Unmarshal(in, &c); err != nil {
		return nil, err
	}

	// copy shortcuts
	for i, r := range c.Rules {
		// copy exchange name to final struct
		if r.RabbitMQ.ExchangeName != "" && r.RabbitMQ.Exchange.Name == "" {
			c.Rules[i].RabbitMQ.Exchange.Name = r.RabbitMQ.ExchangeName
		}
		// copy bind routing key to final struct
		if r.RabbitMQ.RoutingKey != "" && r.RabbitMQ.Bind.RoutingKey == "" {
			c.Rules[i].RabbitMQ.Bind.RoutingKey = r.RabbitMQ.RoutingKey
		}
		// copy queue name to final struct
		if r.RabbitMQ.QueueName != "" && r.RabbitMQ.Queue.Name == "" {
			c.Rules[i].RabbitMQ.Queue.Name = r.RabbitMQ.QueueName
		}
		// copy bind
		c.Rules[i].RabbitMQ.Bind.Exchange = c.Rules[i].RabbitMQ.Exchange
		c.Rules[i].RabbitMQ.Bind.Queue = c.Rules[i].RabbitMQ.Queue
	}

	return
}
