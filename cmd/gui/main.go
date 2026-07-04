package main

import (
	"log"
	"os"

	gioApp "gioui.org/app"
	"github.com/cattyman919/autocv/internal/gui"
)

func main() {
	go func() {
		guiApp, err := gui.NewApp()

		if err != nil {
			log.Fatal(err)
		}

		if err := guiApp.Run(); err != nil {
			log.Fatalln(err)
		}

		os.Exit(0)
	}()

	gioApp.Main()
}
