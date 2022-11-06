package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/xenomote/json_matcher/pattern"
)

const (
	iArg = "i"
	oArg = "o"
	pArg = "p"
	fArg = "f"
)

func main() {
	log.SetFlags(0)

	f := flag.NewFlagSet("json matcher", flag.ExitOnError)
	f.String(iArg, "", "input `file` to read json structures from")
	f.String(oArg, "", "output `file` to write json bindings to")
	f.String(pArg, "", "string `pattern` to match")
	f.String(fArg, "", "`file` containing pattern to match")
	f.Parse(os.Args[1:])

	var err error
	input := os.Stdin
	output := os.Stdout
	var pat io.Reader

	f.Visit(func(f *flag.Flag) {
		name := f.Name
		value := f.Value.String()

		switch name {
		case pArg, fArg:
			if pat != nil {
				log.Fatalln(`pattern already specified`)
			}
		}

		switch name {
		case iArg:
			input, err = os.Open(value)
		case oArg:
			output, err = os.Create(value)
		case pArg:
			pat = strings.NewReader(value)
		case fArg:
			pat, err = os.Open(value)
		default:
			err = fmt.Errorf(`unhandled flag %s`, name)
		}

		if err != nil {
			log.Fatalln(err)
		}
	})

	if pat == nil {
		log.Fatalln(`a pattern must be specified, either use -p or -f to set it`)
	}

	t, err := io.ReadAll(pat)
	if err != nil {
		log.Fatalln(err)
	}

	p, err := pattern.Parse(string(t))
	if err != nil {
		log.Fatalln(err)
	}

	i, err := io.ReadAll(input)
	if err != nil {
		log.Fatalln(err)
	}

	b, err := p.Interpret(string(i))
	if err != nil {
		log.Fatalln(err)
	}

	err = json.NewEncoder(output).Encode(b)
	if err != nil {
		log.Fatalln(err)
	}
}
