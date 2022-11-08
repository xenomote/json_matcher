package pattern

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

//go:generate go run golang.org/x/tools/cmd/goyacc@v0.2.0 -o grammar.go grammar.y

type ValidatedPattern interface {
	Pattern
	Validator
}

func Parse(s string) (ValidatedPattern, error) {
	l := lex{input: []rune(s)}

	if yyParse(&l) != 0 {
		return nil, l.err
	}

	err := l.out.Validate(set{})
	if err != nil {
		return nil, err
	}

	return l.out, nil
}

const EOF = 0

type lex struct {
	input []rune
	i     int

	out ValidatedPattern
	err error
}

func (l *lex) Error(s string) {
	if l.err == nil {
		l.err = fmt.Errorf(s)
	} else {
		l.err = fmt.Errorf("%s: %s", l.err, s)
	}
}

func (l *lex) Lex(lval *yySymType) int {
	for unicode.IsSpace(l.next()) && l.next() != EOF {
		l.take()
	}

	switch l.next() {
	case '[', ']', '{', '}', ':', ',', '<', '>', '=', '?', EOF:
		return int(l.take())
	case '0', '9', '8', '7', '6', '5', '4', '3', '2', '1':
		return l.num(lval)
	case '"':
		return l.str(lval)
	default:
		if l.match("true") {
			return TRUE
		}

		if l.match("false") {
			return FALSE
		}

		if l.match("null") {
			return NULL
		}

		if unicode.IsLetter(l.next()) {
			return l.identifier(lval)
		}

		l.Error(fmt.Sprintf("unrecognised character %c", l.next()))
		return yyErrCode
	}
}

func (l *lex) at(i int) rune {
	if l.i+i < len(l.input) {
		return l.input[l.i+i]
	}

	return EOF
}

func (l *lex) next() rune {
	return l.at(0)
}

func (l *lex) take() rune {
	c := l.next()
	l.i++
	return c
}

func (l *lex) match(s string) bool {
	if len(s) > len(l.input)-l.i {
		return false
	}

	for i, r := range s {
		if r != l.at(i) {
			return false
		}
	}

	l.i += len(s)
	return true
}

func (l *lex) num(lval *yySymType) int {
	var s strings.Builder
	s.WriteRune(l.take())

	for unicode.IsNumber(l.next()) {
		s.WriteRune(l.take())
	}

	if unicode.IsLetter(l.next()) {
		l.Error(fmt.Sprintf("unexpected character %c in number", l.next()))
		return yyErrCode
	}

	n, err := strconv.ParseFloat(s.String(), 64)
	if err != nil {
		l.Error(err.Error())
		return yyErrCode
	}
	lval.num = n

	return NUMBER
}

func (l *lex) str(lval *yySymType) int {
	l.take()
	var s strings.Builder

	for l.next() != '"' && l.next() != EOF {
		s.WriteRune(l.take())
	}

	if l.next() != '"' {
		l.Error("improperly terminated string, reached EOF")
		return yyErrCode
	}
	l.take()
	lval.str = s.String()

	return STRING
}

func (l *lex) identifier(lval *yySymType) int {
	var s strings.Builder
	s.WriteRune(l.take())

	for unicode.IsLetter(l.next()) || unicode.IsNumber(l.next()) {
		s.WriteRune(l.take())
	}

	lval.str = s.String()

	return IDENTIFIER
}
