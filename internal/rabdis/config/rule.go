package config

// Rule represents a Rule section
type Rule struct {
	RabbitMQ RabbitMQ `yaml:"rabbitmq"`
	Redis    Redis    `yaml:"redis"`
}
