package pattern

import (
	"fmt"
)

type BoundLiteral struct {
	Name  Binding
	Value Value
}

func (b BoundLiteral) Match(s string, old_bindings map[string]string) (map[string]string, error) {
	new_bindings := make(map[string]string)

	a_matched, err := b.Value.Match(s, old_bindings)
	if err != nil {
		return nil, err
	}

	for k, v := range a_matched {
		new_bindings[k] = v
	}

	b_matched, err := b.Value.Match(s, old_bindings)
	if err != nil {
		return nil, err
	}

	for k, v := range b_matched {
		new_bindings[k] = v
	}

	return new_bindings, nil
}

func (b BoundLiteral) Validate(bindings map[string]bool) error {
	name := string(b.Name)
	if bindings[name] {
		return fmt.Errorf("illegal self reference to %s", name)
	}
	bindings[name] = true

	value, ok := b.Value.(Validator)
	if !ok {
		return nil
	}

	err := value.Validate(bindings)
	if err != nil {
		return err
	}

	return nil
}

func (b BoundLiteral) String() string {
	return b.Name.String() + b.Value.String()
}
