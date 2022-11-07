package main

import (
	"bufio"
	"encoding/json"
	"flag"
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

type options struct {
	in, pat io.Reader
	out     io.Writer
}

func main() {
	log.SetFlags(0)

	o := args(os.Args[1:])
	input := o.inOr(os.Stdin)
	output := o.outOr(os.Stdout)

	pat, err := io.ReadAll(o.pat)
	if err != nil {
		log.Fatalln(err)
	}

	p, err := pattern.Parse(string(pat))
	if err != nil {
		log.Fatalln(err)
	}

	lines := bufio.NewScanner(input)
	enc := json.NewEncoder(output)

	for lines.Scan() {
		b, err := p.Interpret(lines.Text())
		if err != nil {
			log.Println(err)
			continue
		}

		err = enc.Encode(b)
		if err != nil {
			log.Fatalln(err)
		}
	}

	if lines.Err() != nil {
		log.Fatalln(lines.Err())
	}
}

func args(args []string) options {
	o := options{}

	f := flag.NewFlagSet("json matcher", flag.ExitOnError)
	f.String(iArg, "", "input `file` to read json structures from")
	f.String(oArg, "", "output `file` to write json bindings to")
	f.String(pArg, "", "string `pattern` to match")
	f.String(fArg, "", "`file` containing pattern to match")
	f.Parse(args)

	f.Visit(func(f *flag.Flag) {
		name := f.Name

		switch name {
		case pArg, fArg:
			if o.pat != nil {
				log.Fatalln(`pattern already specified`)
			}

		case iArg:
			if o.in != nil {
				log.Fatalln(`input already specified`)
			}

		case oArg:
			if o.out != nil {
				log.Fatalln(`output already specified`)
			}
		}

		var err error
		value := f.Value.String()

		switch name {
		case iArg:
			o.in, err = os.Open(value)

		case oArg:
			o.out, err = os.Create(value)

		case pArg:
			o.pat = strings.NewReader(value)
			
		case fArg:
			o.pat, err = os.Open(value)
		}

		if err != nil {
			log.Fatalln(err)
		}
	})

	if o.pat == nil {
		log.Fatalln(`a pattern must be specified, either use -p or -f to set it`)
	}

	return o
}

func (o options) inOr(r io.Reader) io.Reader {
	if o.in != nil {
		return o.in
	}

	return r
}

func (o options) outOr(w io.Writer) io.Writer {
	if o.out != nil {
		return o.out
	}

	return w
}
