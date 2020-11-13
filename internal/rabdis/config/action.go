package config

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/zhouzhuojie/conditions"
)

var (
	regexpVariable = regexp.MustCompile(`\{([\w\[\]\.-]+)\}`)
)

// ActionRedis represents a typed redis action
type ActionRedis string

const (
	// ActionIncrement returns increment action for redis
	ActionIncrement ActionRedis = "increment"
	// ActionDecrement returns decrement action for redis
	ActionDecrement ActionRedis = "decrement"
	// ActionDelete returns delete action for redis
	ActionDelete ActionRedis = "delete"
)

// Action represents a redis action section
type Action struct {
	Key        string      `yaml:"key"`
	Action     ActionRedis `yaml:"action"`
	Conditions []string    `yaml:"conditions"` // TODO: improve parser/lexer engine

	content []byte
}

// FinalKey returns the final key after template
func (a *Action) FinalKey() (string, error) {
	key := a.Key

	varNames := regexpVariable.FindAllStringSubmatch(key, -1)
	if len(varNames) == 0 {
		return key, nil
	}

	for _, varName := range varNames {
		res := gjson.Get(a.contentString(), varName[1])
		if !res.Exists() {
			return "", fmt.Errorf("variable %s not found in action %s", varName[0], a.Action)
		}
		if res.String() == "" {
			return "", fmt.Errorf("variable %s is empty in action %s", varName[0], a.Action)
		}
		fmt.Printf("%s  /  %s\n", key, res.String())
		key = strings.ReplaceAll(key, varName[0], res.String())
	}

	return key, nil
}

// ConditionsCheck checks conditions
func (a *Action) ConditionsCheck() bool {
	if len(a.Conditions) == 0 {
		return true
	}

	check := true
	for _, c := range a.Conditions {
		p := conditions.NewParser(strings.NewReader(c))
		expr, err := p.Parse()
		if err != nil {
			// FIXME: Check errors
			continue
		}
		vars := make(map[string]interface{})
		for _, v := range conditions.Variables(expr) {
			r := gjson.Get(string(a.content), v)
			vars[v] = r.String()
		}

		e, err := conditions.Evaluate(expr, vars)
		if err != nil {
			// FIXME: Check errors
			continue
		}
		check = check && e
	}
	return check
}

// SetContent sets content
func (a *Action) SetContent(content []byte) {
	a.content = content
}

// UnmarshalYAML returns an unmarshal YAML implementation
func (a *Action) UnmarshalYAML(u func(interface{}) error) error {
	type rawAction Action
	raw := rawAction{}
	if err := u(&raw); err != nil {
		return err
	}

	// Key is required
	if raw.Key == "" {
		return errors.New("rules.*.redis.actions.*.key is required")
	}

	*a = Action(raw)

	return nil
}

func (a *Action) contentString() string {
	return string(a.content)
}

// // ActionCondition represents an action condition
// type ActionCondition struct {
// 	Selector string `yaml:"selector"`
// 	Required bool   `yaml:"required"`
// }

// // ConditionsCheck checks conditions
// func (a *ActionCondition) ConditionsCheck(content []byte) bool {
// 	if a.Selector != "" {
// 		r := gjson.Get(string(content), a.Selector)
// 		if a.Required && (!r.Exists() || r.String() == "") {
// 			return false
// 		}
// 	}

// 	return true
// }
