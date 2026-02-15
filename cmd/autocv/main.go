package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	// Helper to create boxes with titles/borders for visualization
	newPrimitive := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text).
			SetBorder(true)
	}

	// config := newPrimitive("Config Files")

	configList := tview.NewList()
	configList.AddItem("Projects", "", 0, nil)
	configList.AddItem("Personal Info", "", 0, nil)
	configList.AddItem("Experience", "", 0, nil)
	configList.SetTitle("Config")
	configList.SetTitleAlign(tview.AlignLeft)
	configList.SetBorder(true)

	configList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'j':
			return tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
		case 'k':
			return tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone)
		}
		return event
	})

	cvTypes := newPrimitive("CV Types")
	content := newPrimitive("Content (Editable)\nCV Content or Config")

	grid := tview.NewGrid()
	grid.SetRows(15, 0).
		SetColumns(35, 0)

	// 1. Config (Top Left): Row 0, Col 0, Spans 1 row, 1 col
	grid.AddItem(configList, 0, 0, 1, 1, 0, 0, true)

	// 2. CV Types (Bottom Left): Row 1, Col 0, Spans 1 row, 1 col
	grid.AddItem(cvTypes, 1, 0, 1, 1, 0, 0, false)

	// 3. Content (Right): Row 0, Col 1, SPANS 2 rows (to cover full height), 1 col
	grid.AddItem(content, 0, 1, 2, 1, 0, 0, false)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'q':
			app.Stop()
			return nil
		}
		return event
	})

	if err := app.SetRoot(grid, true).Run(); err != nil {
		panic(err)
	}
}
