package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (self *App) setGlobalKeys(focusable []tview.Primitive) {
	self.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'q':
			self.app.Stop()
			return nil
		}

		switch event.Key() {
		case tcell.KeyTab:
			for i, p := range focusable {
				if p.HasFocus() {
					next := (i + 1) % len(focusable)
					self.app.SetFocus(focusable[next])
					return nil
				}
			}
		case tcell.KeyBacktab:
			for i, p := range focusable {
				if p.HasFocus() {
					prev := (i - 1 + len(focusable)) % len(focusable)
					self.app.SetFocus(focusable[prev])
					return nil
				}
			}
		}

		return event
	})
}
