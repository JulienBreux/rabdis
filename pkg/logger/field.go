package logger

import "github.com/fatih/structs"

// field represent an internal field
type field struct {
	k string
	v interface{}
}

// F helps to create a logger field
func F(key string, value interface{}) Field {
	return &field{k: key, v: value}
}

// E helps to create an error field
func E(err error) Field {
	return &field{k: "error", v: err.Error()}
}

// S helps to create a struct field
func S(value interface{}) []Field {
	var fs []Field

	m := structs.Map(value)

	for k, v := range m {
		if !structs.IsStruct(v) {
			fs = append(fs, &field{k: k, v: v})
		}
	}

	return fs
}

// G create an anonymous struct to group fields
func G(name string, fields ...Field) Field {
	fs := make(map[string]interface{})
	for _, f := range fields {
		fs[f.Key()] = f.Val()
	}

	return F(name, fs)
}

// Key returns the field key
func (f *field) Key() string {
	return f.k
}

// Val returns the field value
func (f *field) Val() interface{} {
	return f.v
}
