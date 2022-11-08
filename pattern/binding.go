package pattern

import (
	"encoding/json"
	"fmt"
)

type Binding Identifier

func (b Binding) Match(s []byte, _ bindings) (bindings, error) {
	var out interface{}
	err := json.Unmarshal(s, &out)
	if err != nil {
		return nil, err
	}

	return bindings{string(b): out}, nil
}

func (b Binding) Validate(s set) error {
	if s[string(b)] {
		return fmt.Errorf("duplicate binding %s", b)
	}

	s[string(b)] = true
	return nil
}

func (b Binding) String() string {
	return "<=" + string(b) + ">"
}
