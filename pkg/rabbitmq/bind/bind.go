package bind

import (
	"github.com/julienbreux/rabdis/pkg/rabbitmq/exchange"
	"github.com/julienbreux/rabdis/pkg/rabbitmq/queue"
)

// Bind represents the bin
type Bind struct {
	Exchange exchange.Exchange `yaml:"-"`
	Queue    queue.Queue       `yaml:"-"`

	RoutingKey string `yaml:"routingKey"`

	NoWait bool `yaml:"noWait"`

	Arguments map[string]interface{} `yaml:"arguments"`
}

// UnmarshalYAML returns an unmarshal YAML implementation
func (b *Bind) UnmarshalYAML(u func(interface{}) error) error {
	type rawBind Bind
	raw := rawBind{
		NoWait: false,

		Arguments: make(map[string]interface{}),
	}
	if err := u(&raw); err != nil {
		return err
	}

	*b = Bind(raw)

	return nil
}
