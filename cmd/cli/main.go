package main

import (
	"log"

	"github.com/cattyman919/autocv/internal/core"
)

func main() {
	app, err := core.NewApp()
	if err != nil {
		log.Fatalln(err)
	}

	app.Run()
}
