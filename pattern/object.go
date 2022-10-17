package pattern

import (
	"encoding/json"
	"fmt"
)

func (o Object) Validate(bindings map[string]bool) error {
	keys := make(map[string]bool)
	for _, d := range o.Definitions {
		key := d.Key.String()
		if _, exists := keys[key]; exists {
			return fmt.Errorf("duplicate key %s", key)
		}

		keys[key] = true
	}

	for _, d := range o.Definitions {
		if err := d.Validate(bindings); err != nil {
			return fmt.Errorf("at key %s: %s", d.Key, err)
		}
	}

	return nil
}

func (o Object) Interpret(s string) (bindings map[string]string, err error) {
	return o.Match(s, map[string]string{})
}

func (o Object) Match(s string, old_bindings map[string]string) (new_bindings map[string]string, err error) {
	new_bindings = map[string]string{}

	bindings := map[string]string{}
	for k, v := range old_bindings {
		bindings[k] = v
	}

	var input map[string]json.RawMessage
	err = json.Unmarshal([]byte(s), &input)
	if err != nil {
		input := "input"
		if len(s) < 10 {
			input = "\"" + s + "\""
		}

		return nil, fmt.Errorf("%s could not be interpreted as an object: %w", input, err)
	}

	for _, definition := range o.Definitions {
		key, err := definition.Key.Key()
		if err != nil {
			return nil, err
		}

		prefix := ""
		if definition.Key.String() != key {
			prefix = definition.Key.String() + " = "
		}

		value, key_exists := input[key]
		if !key_exists {
			if definition.Optional {
				continue
			} else {
				return nil, fmt.Errorf("object did not contain required field %s\"%s\"", prefix, key)
			}
		}

		matched, err := definition.Assignment.Match(string(value), bindings)
		if err != nil {
			return nil, fmt.Errorf("could not match field %s\"%s\": %s", prefix, key, err)
		}

		for k, v := range matched {
			if _, k_exists := bindings[k]; k_exists {
				return nil, fmt.Errorf("binding for %s already exists and cannot be overwritten", k)
			}

			bindings[k] = v
			new_bindings[k] = v
		}
	}

	return new_bindings, nil
}

func (o Object) String() string {
	s := "{"

	for i, definition := range o.Definitions {
		s += "\n" + indent(definition.String())

		if i < len(o.Definitions)-1 {
			s += ","
		} else {
			s += "\n"
		}
	}

	s += "}"

	return s
}

func (o ObjectDefinition) String() string {
	op := ""
	if o.Optional {
		op = "?"
	}
	return o.Key.String() + op + ": " + o.Assignment.String()
}
