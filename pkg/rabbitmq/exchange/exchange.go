package exchange

// ExchangeType represents an exchange type
type ExchangeType string

const (
	// Type "direct" delivers messages to queues based on
	// the message routing key.
	TypeDirect ExchangeType = "direct"

	// Type "fanout" routes messages to all of the queues
	// that are bound to it and the routing key is ignored.
	TypeFanout ExchangeType = "fanout"

	// Type "headers" is designed for routing on multiple
	// attributes that are more easily expressed as message
	// headers than a routing key.
	TypeHeaders ExchangeType = "headers"

	// Type "topic" route messages to one or many queues
	// based on matching between a message routing key and
	// the pattern that was used to bind a queue to an
	// exchange.
	TypeTopic ExchangeType = "topic"
)

// Exchange represents an exchange
type Exchange struct {
	Name string `yaml:"name"`

	Type ExchangeType `yaml:"type"`

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
