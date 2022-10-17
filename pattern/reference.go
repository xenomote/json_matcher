package pattern

import "fmt"

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

func (r Reference) Validate(bindings map[string]bool) error {
	bind := string(r[0].Identifier)
	exists := bindings[bind]
	bindings[bind] = true

	if !exists {
		return fmt.Errorf("reference to %s before it was bound", r)
	}

	return nil
}

func (r Reference) Match(s string, bindings map[string]string) (map[string]string, error) {
	bind := string(r[0].Identifier)
	o, exists := bindings[bind]
	if !exists {
		return nil, fmt.Errorf("referenced binding %s was not available, was it matched in an optional section?", r)
	}

	if s != o {
		return nil, fmt.Errorf("reference to binding %s did not match expected value: %s != %s", r, s, o)
	}

	return map[string]string{}, nil
}
