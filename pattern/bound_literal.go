package pattern

import "fmt"

type BoundLiteral struct {
	Name  Binding
	Value Value
}

func (b BoundLiteral) Match(s []byte, bOld bindings) (bindings, error) {
	bNew := bindings{}

	matched, err := b.Name.Match(s, bOld)
	if err != nil {
		return nil, err
	}

	for k, v := range matched {
		bNew[k] = v
	}

	matched, err = b.Value.Match(s, bOld)
	if err != nil {
		return nil, err
	}

	for k, v := range matched {
		bNew[k] = v
	}

	return bNew, nil
}

func (b BoundLiteral) Validate(s set) error {
	ref := b.Name.String()
	if s[ref] {
		return fmt.Errorf("illegal self reference to %s", ref)
	}
	s[ref] = true

	value, ok := b.Value.(Validator)
	if !ok {
		return nil
	}

	err := value.Validate(s)
	if err != nil {
		return err
	}

	return nil
}

func (b BoundLiteral) String() string {
	return b.Name.String() + " " + b.Value.String()
}
