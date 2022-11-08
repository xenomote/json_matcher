package pattern

import (
	"strings"
)

type bindings = map[string]interface{}
type set = map[string]bool

type Pattern interface {
	Interpret(string) (bindings, error)
}

type Validator interface {
	Validate(set) error
}

type Value interface {
	Match(string, bindings) (bindings, error)
	String() string
}

type Reference []OptionalIdentifier
type OptionalIdentifier struct {
	Identifier
	Optional bool
}
type Identifier string

const INDENT = "    "

func indent(s string) string {
	return INDENT + strings.ReplaceAll(s, "\n", "\n"+INDENT)
}
