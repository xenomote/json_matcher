package pattern

import (
	"fmt"
	"strconv"
)

type Number float64

func (i Number) Index() (int, error) {
	return int(i), nil
}

func (n Number) String() string {
	return fmt.Sprint(float64(n))
}

func (n Number) Match(s []byte, _ bindings) (bindings, error) {
	m, err := strconv.ParseFloat(string(s), 64)
	if err != nil {
		return nil, fmt.Errorf("expected %s but matched value %s could not be interpreted as a number", n, s)
	}

	if float64(n) != m {
		return nil, fmt.Errorf("expected %s but matched value %s", n, s)
	}

	return nil, nil
}
