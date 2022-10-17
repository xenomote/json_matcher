package pattern

import (
	"fmt"
	"strconv"
)

func (Null) Validate(map[string]bool) error {
	return nil
}

func (Null) Match(s string, _ map[string]string) (map[string]string, error) {
	if s != "null" {
		return nil, fmt.Errorf("expected null but matched %s", s)
	}

	return map[string]string{}, nil
}

func (Null) String() string {
	return "null"
}

func (i Number) Index() (int, error) {
	return int(i), nil
}

func (n Number) String() string {
	return fmt.Sprint(float64(n))
}

func (Number) Validate(map[string]bool) error {
	return nil
}

func (n Number) Match(s string, _ map[string]string) (map[string]string, error) {
	m, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil, fmt.Errorf("expected %s but matched value %s could not be interpreted as a number", n, s)
	}

	if float64(n) != m {
		return nil, fmt.Errorf("expected %s but matched value %s", n, s)
	}

	return nil, nil
}

func (k String) Key() (string, error) {
	return string(k), nil
}

func (s String) String() string {
	return string(s)
}

func (String) Validate(map[string]bool) error {
	return nil
}

func (t String) Match(s string, _ map[string]string) (map[string]string, error) {
	if s[0] != '"' || s[len(s)-1] != '"' {
		return nil, fmt.Errorf("expected \"%s\" but matched value \"%s\" could not be interpreted as a string", t, s)
	}

	if string(t) != s[1:len(s)-1] {
		return nil, fmt.Errorf("expected \"%s\" but matched value %s", t, s)
	}

	return nil, nil
}

func (b Boolean) String() string {
	if b {
		return "true"
	} else {
		return "false"
	}
}

func (Boolean) Validate(map[string]bool) error {
	return nil
}

func (b Boolean) Match(s string, _ map[string]string) (map[string]string, error) {
	x, err := strconv.ParseBool(s)
	if err != nil {
		return nil, fmt.Errorf("expected %s but matched value %s could not be interpreted as a boolean", b, s)
	}

	if bool(b) != x {
		return nil, fmt.Errorf("expected %s but matched value %s", b, s)
	}

	return nil, nil
}

func (b BoundLiteral) Match(s string, old_bindings map[string]string) (map[string]string, error) {
	new_bindings := make(map[string]string)

	a_matched, err := b.Assignment.Match(s, old_bindings)
	if err != nil {
		return nil, err
	}

	for k, v := range a_matched {
		new_bindings[k] = v
	}

	b_matched, err := b.Binding.Match(s, old_bindings)
	if err != nil {
		return nil, err
	}

	for k, v := range b_matched {
		new_bindings[k] = v
	}

	return new_bindings, nil
}

func (b BoundLiteral) Validate(bindings map[string]bool) error {
	bind := string(b.Binding)
	if err := b.Assignment.Validate(bindings); err != nil {
		if bindings[bind] {
			return fmt.Errorf("illegal self reference to %s", b.Binding)
		} else {
			return err
		}
	}

	bindings[bind] = true
	return nil
}

func (b BoundLiteral) String() string {
	return b.Binding.String() + b.Assignment.String()
}
