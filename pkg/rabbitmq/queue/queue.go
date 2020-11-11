package queue

// Queue represents a RabbitMQ queue
type Queue struct {
	Name string `yaml:"name"`

	Durable    bool `yaml:"durable"`
	Exclusive  bool `yaml:"exclusive"`
	AutoDelete bool `yaml:"autoDelete"`
	NoWait     bool `yaml:"noWait"`

	Arguments map[string]interface{} `yaml:"arguments"`
}

// UnmarshalYAML returns an unmarshal YAML implementation
func (q *Queue) UnmarshalYAML(u func(interface{}) error) error {
	type rawQueue Queue
	raw := rawQueue{
		Durable:    true,
		Exclusive:  false,
		AutoDelete: false,
		NoWait:     false,

		Arguments: make(map[string]interface{}),
	}
	if err := u(&raw); err != nil {
		return err
	}

	*q = Queue(raw)

	return nil
}
