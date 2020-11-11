package config

// Redis represents a Redis section
type Redis struct {
	Actions []Action `yaml:"actions"`
}
