package main

import (
	"github.com/cattyman919/autocv/internal/ui"
)

func main() {
	app := ui.NewApp()

	if err := app.Run(); err != nil {
		panic(err)
	}
}
