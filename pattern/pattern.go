package pattern

import "strings"

type Pattern interface {
	Interpret(string) (map[string]string, error)
}

type Validator interface {
	Validate(bindings map[string]bool) error
}

type Value interface {
	Match(string, map[string]string) (map[string]string, error)
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
