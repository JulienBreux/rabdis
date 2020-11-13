package config_test

import (
	"testing"

	"github.com/julienbreux/rabdis/internal/rabdis/config"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

var conditionsTests = []struct {
	content    []byte
	conditions []string
	result     bool
}{
	{[]byte(`{"id":"a"}`), []string{}, true},
	{[]byte(`{"id":"b"}`), []string{"{id} == \"b\""}, true},
	{[]byte(`{"id":"c","name":""}`), []string{"{name} != \"\""}, false},
	{[]byte(`{"id":"d","name":""}`), []string{"{mail} != \"\""}, false},
	{[]byte(`{"id":"e","name":""}`), []string{"{name} == \"\""}, true},
}

func TestActionsConditions(t *testing.T) {
	for _, tt := range conditionsTests {
		tt := tt
		t.Run(string(tt.content), func(t *testing.T) {
			a := &config.Action{
				Key:        "",
				Action:     "",
				Conditions: tt.conditions,
			}
			a.SetContent(tt.content)

			assert.Equal(t, tt.result, a.ConditionsCheck())
		})
	}
}

func TestActionUnmarshalYAML(t *testing.T) {
	var a config.Action

	in := []byte(`
key: key
action: delete
conditions:
- '{id} != ""'
`)
	assert.NoError(t, yaml.Unmarshal(in, &a))
	assert.Equal(t, a.Key, "key")
	assert.Equal(t, a.Action, config.ActionDelete)
}

func TestActionUnmarshalYAMLCustomError(t *testing.T) {
	var a config.Action

	in := []byte("key:\naction: delete\n")
	assert.Error(t, yaml.Unmarshal(in, &a))
}

var finalKeyTests = []struct {
	content  []byte
	key      string
	finalKey string
	err      bool
	errMsg   string
}{
	// No errors
	{[]byte(`{"id":"a"}`), "users::all", "users::all", false, ""},
	{[]byte(`{"id":"b"}`), "users::{id}", "users::b", false, ""},
	{[]byte(`{"id":"c","groups":[{"id":"a"}]}`), "users::{id}::{groups.0.id}", "users::c::a", false, ""},
	{[]byte(`{"id":0}`), "users::0", "users::0", false, ""},

	// Errors
	{[]byte(`{"id":"d"}`), "users::{id.instructor}", "", true, "variable {id.instructor} not found in action delete"},
	{[]byte(`{"id":""}`), "users::{id}", "", true, "variable {id} is empty in action delete"},
	{[]byte(`{`), "users::{id}", "", true, "variable {id} not found in action delete"},
}

func TestFinalKey(t *testing.T) {
	for _, tt := range finalKeyTests {
		tt := tt
		t.Run(string(tt.content), func(t *testing.T) {
			a := &config.Action{
				Key:    tt.key,
				Action: "delete",
			}
			a.SetContent(tt.content)

			fn, err := a.FinalKey()
			if tt.err {
				assert.Error(t, err)
				assert.Equal(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.finalKey, fn)
		})
	}
}
