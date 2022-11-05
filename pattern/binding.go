package pattern

import "fmt"

type Binding Identifier

func (b Binding) Match(s string, bindings map[string]string) (map[string]string, error) {
	return map[string]string{string(b): s}, nil
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
