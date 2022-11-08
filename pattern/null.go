package pattern

import "fmt"

type Null struct{}

func (Null) Match(s string, _ bindings) (bindings, error) {
	if s != "null" {
		return nil, fmt.Errorf("expected null but matched %s", s)
	}

	return bindings{}, nil
}

func (Null) String() string {
	return "null"
}