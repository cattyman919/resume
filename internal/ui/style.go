package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// applyFocusStyle adds focus/blur color toggling to any tview primitive.
// We pass *tview.Box because all standard tview widgets embed it.
func applyFocusStyle(box *tview.Box, activeColor, inactiveColor tcell.Color) {
	if box.HasFocus() {
		box.SetBorderColor(activeColor)
		box.SetTitleColor(activeColor)
	} else {
		box.SetBorderColor(inactiveColor)
		box.SetTitleColor(inactiveColor)
	}
	// 1. Set the initial state (unfocused)

	// 2. What happens when the panel gets focus
	box.SetFocusFunc(func() {
		box.SetBorderColor(activeColor)
		box.SetTitleColor(activeColor)
	})

	// 3. What happens when the panel loses focus
	box.SetBlurFunc(func() {
		box.SetBorderColor(inactiveColor)
		box.SetTitleColor(inactiveColor)
	})
}
