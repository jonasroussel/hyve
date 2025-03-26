package main

import (
	"log"

	"github.com/jonasroussel/hyve"
)

func main() {
	app := hyve.New()

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
