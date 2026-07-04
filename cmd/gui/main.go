package main

import (
	"log"

	"github.com/cattyman919/autocv/internal/gui"
)

func main() {
	guiApp, err := gui.NewApp()
	if err != nil {
		log.Fatalln(err)
	}
	guiApp.Run()
}
