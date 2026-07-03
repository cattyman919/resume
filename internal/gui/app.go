package gui

import (
	"image/color"

	"fyne.io/fyne/v2"
	fyneApp "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"github.com/cattyman919/autocv/internal/core"
)

type App struct {
	coreApp *core.App
	window  fyne.Window
}

func NewApp() (*App, error) {
	coreApp, err := core.NewApp()
	if err != nil {
		return nil, err
	}

	guiApp := fyneApp.New()
	w := guiApp.NewWindow("AutoCV")

	return &App{
		coreApp: coreApp,
		window:  w,
	}, nil
}

func (a *App) Run() {
	w := a.window

	green := color.NRGBA{R: 0, G: 180, B: 0, A: 255}
	yamlConfig := canvas.NewText("CV Type Config", green)
	pdfViewer := canvas.NewText("PDF viewer", green)

	mainContent := container.NewHSplit(yamlConfig, pdfViewer)

	content := container.NewHSplit(sidebarView(a.coreApp.CVConfig), mainContent)
	content.Offset = 0.12

	w.SetContent(content)
	w.ShowAndRun()
}
