package pattern

import "fmt"

type Null struct{}

func (Null) Match(s []byte, _ bindings) (bindings, error) {
	if string(s) != "null" {
		return nil, fmt.Errorf("expected null but matched %s", s)
	}

	return bindings{}, nil
}

func (Null) String() string {
	return "null"
}
