package pattern

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type Binding Identifier

func (b Binding) Match(s string, bindings map[string]string) (map[string]string, error) {
	out := bytes.Buffer{}
	err := json.Compact(&out, []byte(s))
	if err != nil {
		return nil, err
	}

	return map[string]string{string(b): out.String()}, nil
}

func (b Binding) Validate(bindings map[string]bool) error {
	if _, exists := bindings[string(b)]; exists {
		return fmt.Errorf("duplicate binding %s", b)
	}

	bindings[string(b)] = true
	return nil
}

func (b Binding) String() string {
	return "<=" + string(b) + ">"
}
