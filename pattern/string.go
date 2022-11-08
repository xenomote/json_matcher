package pattern

import "fmt"

type String string

func (k String) Key() (string, error) {
	return string(k), nil
}

func (s String) String() string {
	return string(s)
}

func (t String) Match(s []byte, _ bindings) (bindings, error) {
	if len(s) < 2 || s[0] != '"' || s[len(s)-1] != '"' {
		return nil, fmt.Errorf(`value '%s' could not be interpreted as a string`, s)
	}

	if string(t) != string(s[1:len(s)-1]) {
		return nil, fmt.Errorf(`expected '%s' but matched value '%s'`, t, s)
	}

	return nil, nil
}
