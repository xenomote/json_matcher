package main

import (
	"log"

	"github.com/xenomote/json_matcher/pattern"
)

func main() {
	log.SetFlags(0)

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
				"set": null
			}
		}
	}
	`)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(b)
}
