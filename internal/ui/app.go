package ui

import (
	"github.com/cattyman919/autocv/internal/ui/views"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func newPrimitive(text string) tview.Primitive {
	return tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText(text).
		SetBorder(true)
}

type App struct {
	app          *tview.Application
	configPanel  *views.ConfigPanel
	cvTypesPanel *views.CVTypesPanel
}

func NewApp() *App {
	tviewApp := tview.NewApplication()
	configPanel := views.NewConfigPanel()
	cvTypesPanel := views.NewCvTypesPanel()

	activeColor := tcell.ColorGreen
	inactiveColor := tcell.ColorDimGray

	content := newPrimitive("Content (Editable)\nCV Content or Config")

	grid := tview.NewGrid()
	grid.SetRows(15, 0).
		SetColumns(35, 0)

	// 1. Config (Top Left): Row 0, Col 0, Spans 1 row, 1 col
	grid.AddItem(configPanel.GetPrimitive(), 0, 0, 1, 1, 0, 0, true)

	// 2. CV Types (Bottom Left): Row 1, Col 0, Spans 1 row, 1 col
	grid.AddItem(cvTypesPanel.GetPrimitive(), 1, 0, 1, 1, 0, 0, false)

	// 3. Content (Right): Row 0, Col 1, SPANS 2 rows (to cover full height), 1 col
	grid.AddItem(content, 0, 1, 2, 1, 0, 0, false)

	tviewApp = tviewApp.SetRoot(grid, true)

	app := &App{
		app:          tviewApp,
		configPanel:  configPanel,
		cvTypesPanel: cvTypesPanel,
	}

	focusable := []tview.Primitive{configPanel.GetPrimitive(), cvTypesPanel.GetPrimitive(), content}
	app.setGlobalKeys(focusable)

	// Apply Styling
	applyFocusStyle(configPanel.GetBox(), activeColor, inactiveColor)
	applyFocusStyle(cvTypesPanel.GetBox(), activeColor, inactiveColor)

	return app
}

func (self *App) Run() error {
	return self.app.Run()
}
