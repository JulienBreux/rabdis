package format

import "gopkg.in/yaml.v2"

// ToYAML returns value in YAML
func ToYAML(v interface{}) ([]byte, error) {
	return yaml.Marshal(v)
}
