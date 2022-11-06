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
	f.String(iArg, "", "sets the `file` to read input from")
	f.String(oArg, "", "sets the `file` to write output to")
	f.String(pArg, "", "sets the `pattern` to use")
	f.String(fArg, "", "sets the `file` to read the pattern from")
	f.Parse(os.Args[1:])

	var err error
	input := os.Stdin
	output := os.Stdout
	var pat io.Reader

	f.Visit(func(f *flag.Flag) {
		switch f.Name {
		case iArg:
			input, err = os.Open(f.Value.String())
		case oArg:
			output, err = os.Create(f.Value.String())
		case pArg:
			pat = strings.NewReader(f.Value.String())
		case fArg:
			pat, err = os.Open(f.Value.String())
		default:
			err = fmt.Errorf(`unhandled flag %s`, f.Name)
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
