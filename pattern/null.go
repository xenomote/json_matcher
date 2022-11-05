package pattern

import "fmt"

type Null struct{}

func (Null) Match(s string, _ map[string]string) (map[string]string, error) {
	if s != "null" {
		return nil, fmt.Errorf("expected null but matched %s", s)
	}

	return map[string]string{}, nil
}

func (Null) String() string {
	return "null"
}