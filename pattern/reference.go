package pattern

import (
	"encoding/json"
	"fmt"
)

func (r Reference) String() string {
	s := "<"

	for i, identifier := range r {
		s += string(identifier.Identifier)

		if identifier.Optional {
			s += "?"
		} else if i != len(r)-1 {
			s += "."
		}
	}

	s += ">"

	return s
}

func (r Reference) Validate(s set) error {
	ref := string(r[0].Identifier)

	if !s[ref] {
		return fmt.Errorf("reference to %s before it was bound", r)
	}

	return nil
}

func (r Reference) Match(s []byte, b bindings) (bindings, error) {
	ref := string(r[0].Identifier)

	y, exists := b[ref]
	if !exists {
		return nil, fmt.Errorf("referenced binding %s was not available, was it matched in an optional section?", r)
	}

	var x interface{}
	err := json.Unmarshal(s, &x)
	if err != nil {
		return nil, fmt.Errorf(`could not unmarshal bound value to match: %s`, err)
	}

	if !Matches(x, y) {
		xb, _ := json.Marshal(x)
		yb, _ := json.Marshal(y)
		return nil, fmt.Errorf(`reference to binding '%s' did not match expected value: '%s' != '%s'`, r, string(xb), string(yb))
	}

	return bindings{}, nil
}

func Matches(a, b interface{}) bool {
	if v, ok := a.(json.RawMessage); ok {
		json.Unmarshal(v, &a)
	}

	if v, ok := b.(json.RawMessage); ok {
		json.Unmarshal(v, &b)
	}

	switch a := a.(type) {
	case map[string]interface{}:
		b, ok := b.(map[string]interface{})
		if !ok {
			return false
		}

		for k := range a {
			if _, exists := b[k]; !exists {
				return false
			}

			if !Matches(a[k], b[k]) {
				return false
			}
		}

		for k := range b {
			if _, exists := a[k]; !exists {
				return false
			}
		}

		return true

	case []interface{}:
		b, ok := b.([]interface{})
		if !ok {
			return false
		}

		if len(a) != len(b) {
			return false
		}

		for i := range a {
			if !Matches(a[i], b[i]) {
				return false
			}
		}

		return true
	}

	switch a.(type) {
	case bool, float64, string, nil:
		return a == b

	default:
		panic(fmt.Sprintf(`impossible type %T`, a))
	}
}
