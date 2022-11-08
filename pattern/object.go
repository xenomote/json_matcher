package pattern

import (
	"encoding/json"
	"fmt"
)

type Object struct {
	Fields []Field
}

type Field struct {
	Key      Key
	Value    Value
	Optional bool
}

type Key interface {
	Key() (string, error)
	String() string
}

func (o Object) Validate(s set) error {
	sObj := set{}
	for _, f := range o.Fields {
		ref := f.Key.String()
		if sObj[ref] {
			return fmt.Errorf("duplicate key %s", ref)
		}
		sObj[ref] = true
	}

	for _, f := range o.Fields {
		value, ok := f.Value.(Validator)
		if !ok {
			continue
		}

		if err := value.Validate(s); err != nil {
			return fmt.Errorf("at key %s: %s", f.Key, err)
		}
	}

	return nil
}

func (o Object) Interpret(s string) (bindings, error) {
	return o.Match(s, bindings{})
}

func (o Object) Match(s string, bOld bindings) (bindings, error) {
	bCopy := bindings{}
	for k, v := range bOld {
		bCopy[k] = v
	}

	var input map[string]json.RawMessage
	err := json.Unmarshal([]byte(s), &input)
	if err != nil {
		input := "input"
		if len(s) < 10 {
			input = "'" + s + "'"
		}

		return nil, fmt.Errorf("%s could not be interpreted as an object", input)
	}

	bNew := bindings{}
	for _, definition := range o.Fields {
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
			}

			return nil, fmt.Errorf("object did not contain required field %s\"%s\"", prefix, key)
		}

		matched, err := definition.Value.Match(string(value), bCopy)
		if err != nil {
			return nil, fmt.Errorf("could not match field %s\"%s\": %s", prefix, key, err)
		}

		for k, v := range matched {
			if _, k_exists := bCopy[k]; k_exists {
				return nil, fmt.Errorf("binding for %s already exists and cannot be overwritten", k)
			}

			bCopy[k] = v
			bNew[k] = v
		}
	}

	return bNew, nil
}

func (o Object) String() string {
	s := "{"

	for i, definition := range o.Fields {
		s += "\n" + indent(definition.String())

		if i < len(o.Fields)-1 {
			s += ","
		} else {
			s += "\n"
		}
	}

	s += "}"

	return s
}

func (f Field) String() string {
	op := ""
	if f.Optional {
		op = "?"
	}

	return f.Key.String() + op + ": " + f.Value.String()
}
