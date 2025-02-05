package exchange

// Type represents an exchange type
type Type string

const (
	// TypeDirect delivers messages to queues based on
	// the message routing key.
	TypeDirect Type = "direct"

	// TypeFanout routes messages to all of the queues
	// that are bound to it and the routing key is ignored.
	TypeFanout Type = "fanout"

	// TypeHeaders is designed for routing on multiple
	// attributes that are more easily expressed as message
	// headers than a routing key.
	TypeHeaders Type = "headers"

	// TypeTopic route messages to one or many queues
	// based on matching between a message routing key and
	// the pattern that was used to bind a queue to an
	// exchange.
	TypeTopic Type = "topic"
)

// Exchange represents an exchange
type Exchange struct {
	Name string `yaml:"name"`

	Type Type `yaml:"type"`

	Durable    bool `yaml:"durable"`
	AutoDelete bool `yaml:"autoDelete"`
	Internal   bool `yaml:"internal"`
	NoWait     bool `yaml:"noWait"`

	Arguments map[string]interface{} `yaml:"arguments"`
}

// UnmarshalYAML returns an unmarshal YAML implementation
func (e *Exchange) UnmarshalYAML(u func(interface{}) error) error {
	type rawExchange Exchange
	raw := rawExchange{
		Type: TypeTopic,

		Durable:    true,
		AutoDelete: false,
		Internal:   false,
		NoWait:     false,

		Arguments: make(map[string]interface{}),
	}
	if err := u(&raw); err != nil {
		return err
	}

	*e = Exchange(raw)

	return nil
}
