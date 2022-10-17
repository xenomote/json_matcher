package pattern

import "strings"

type Pattern interface {
	Interpret(string) (map[string]string, error)
	Assignment
}

type Object struct {
	Definitions []ObjectDefinition
}

type ObjectDefinition struct {
	Key      Key
	Optional bool
	Assignment
}

type Array struct {
	Definitions []ArrayDefinition
}

type ArrayDefinition struct {
	Index    Index
	Optional bool
	Assignment
}

type Key interface {
	Key() (string, error)
	String() string
}
type Index interface {
	Index() (int, error)
	String() string
}

type Assignment interface {
	Match(string, map[string]string) (map[string]string, error)
	Validate(bindings map[string]bool) error
	String() string
}

type Null struct{}
type Identifier string
type String string
type Number float64
type Boolean bool

type Binding Identifier
type Reference []OptionalIdentifier

type BoundLiteral struct {
	Binding
	Assignment
}

type OptionalIdentifier struct {
	Identifier
	Optional bool
}

const INDENT = "    "

func indent(s string) string {
	return INDENT + strings.ReplaceAll(s, "\n", "\n"+INDENT)
}
