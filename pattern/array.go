package pattern

import (
	"encoding/json"
	"fmt"
)

func (a Array) Interpret(s string) (bindings map[string]string, err error) {
	return a.Match(s, map[string]string{})
}

func (a Array) Match(s string, old_bindings map[string]string) (new_bindings map[string]string, err error) {
	new_bindings = map[string]string{}

	bindings := map[string]string{}
	for k, v := range old_bindings {
		bindings[k] = v
	}

	var input []json.RawMessage
	err = json.Unmarshal([]byte(s), &input)
	if err != nil {
		input := "input"
		if len(s) < 10 {
			input = "\"" + s + "\""
		}

		return nil, fmt.Errorf("%s could not be interpreted as an array: %w", input, err)
	}

	for _, definition := range a.Definitions {
		index, err := definition.Index.Index()
		if err != nil {
			return nil, err
		}

		prefix := ""
		if definition.Index.String() != fmt.Sprint(index) {
			prefix = definition.Index.String() + " = "
		}

		if index >= len(input) {
			if definition.Optional {
				continue
			} else {
				return nil, fmt.Errorf("array was not long enough to contain required index %s%d", prefix, index)
			}
		}

		value := input[index]
		matched_bindings, err := definition.Assignment.Match(string(value), bindings)
		if err != nil {
			return nil, fmt.Errorf("could not match index %s%d: %s", prefix, index, err)
		}

		for k, v := range matched_bindings {
			if _, k_exists := bindings[k]; k_exists {
				return nil, fmt.Errorf("binding for %s already exists and cannot be overwritten", k)
			}

			bindings[k] = v
			new_bindings[k] = v
		}
	}

	return new_bindings, nil
}

func (a Array) Validate(bindings map[string]bool) error {
	indices := make(map[string]bool)
	for _, d := range a.Definitions {
		if _, exists := indices[d.Index.String()]; exists {
			return fmt.Errorf("duplicate index %s", d.Index.String())
		}

		indices[d.Index.String()] = true
	}

	for _, d := range a.Definitions {
		if err := d.Validate(bindings); err != nil {
			return fmt.Errorf("at index %s: %s", d.Index, err)
		}
	}

	return nil
}

func (a Array) String() string {
	s := "["

	for i, definition := range a.Definitions {
		s += "\n" + indent(definition.String())

		if i < len(a.Definitions)-1 {
			s += ","
		} else {
			s += "\n"
		}
	}

	s += "]"

	return s
}

func (o ArrayDefinition) String() string {
	op := ""
	if o.Optional {
		op = "?"
	}
	return o.Index.String() + op + ": " + o.Assignment.String()
}
