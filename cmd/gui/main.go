package main

import (
	"log"

	"github.com/cattyman919/autocv/internal/gui"
)

func main() {
	app, err := gui.NewApp()
	if err != nil {
		log.Fatal(err)
	}

	app.Run()
}
