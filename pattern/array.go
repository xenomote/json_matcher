package pattern

import (
	"encoding/json"
	"fmt"
)

type Array struct {
	Elements []Element
}

type Element struct {
	Index    Index
	Value    Value
	Optional bool
}

type Index interface {
	Index() (int, error)
	String() string
}

func (a Array) Interpret(s string) (bindings, error) {
	return a.Match(s, bindings{})
}

func (a Array) Match(s string, bOld bindings) (bindings, error) {
	bCopy := bindings{}
	for k, v := range bOld {
		bCopy[k] = v
	}

	var input []json.RawMessage
	err := json.Unmarshal([]byte(s), &input)
	if err != nil {
		input := "input"
		if len(s) < 10 {
			input = "'" + s + "'"
		}

		return nil, fmt.Errorf("%s could not be interpreted as an array", input)
	}

	bNew := bindings{}
	for _, definition := range a.Elements {
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
			}

			return nil, fmt.Errorf("array was not long enough to contain required index %s%d", prefix, index)
		}

		value := input[index]
		matched_bindings, err := definition.Value.Match(string(value), bCopy)
		if err != nil {
			return nil, fmt.Errorf("could not match index %s%d: %s", prefix, index, err)
		}

		for k, v := range matched_bindings {
			if _, k_exists := bCopy[k]; k_exists {
				return nil, fmt.Errorf("binding for %s already exists and cannot be overwritten", k)
			}

			bCopy[k] = v
			bNew[k] = v
		}
	}

	return bNew, nil
}

func (a Array) Validate(s set) error {
	indices := set{}
	for _, e := range a.Elements {
		index := e.Index.String()
		if _, exists := indices[index]; exists {
			return fmt.Errorf("duplicate index %s", index)
		}

		indices[index] = true
	}

	for _, e := range a.Elements {
		value, ok := e.Value.(Validator)
		if !ok {
			continue
		}

		if err := value.Validate(s); err != nil {
			return fmt.Errorf("at index %s: %s", e.Index, err)
		}
	}

	return nil
}

func (a Array) String() string {
	s := "["

	for i, definition := range a.Elements {
		s += "\n" + indent(definition.String())

		if i < len(a.Elements)-1 {
			s += ","
		} else {
			s += "\n"
		}
	}

	s += "]"

	return s
}

func (e Element) String() string {
	op := ""
	if e.Optional {
		op = "?"
	}

	return e.Index.String() + op + ": " + e.Value.String()
}
