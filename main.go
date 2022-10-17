package main

import (
	"fmt"
	"log"

	"github.com/xenomote/object_language/pattern"
)

func main() {
	p, err := pattern.Parse(`
	{
		"type": "event.source",
		"data": {
			"metadata": <=x> {
				"set": null
			}
		}
	}
	`)
	if err != nil {
		log.Fatalln(err)
	}

	b, err := p.Interpret(`
	{
		"type": "event.source",
		"data": {
			"metadata": {
				"set": "null"
			}
		}
	}
	`)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(b)
}
