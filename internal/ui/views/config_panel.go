package views

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ConfigPanel struct {
	configs *tview.List
}

func NewConfigPanel() *ConfigPanel {
	configList := tview.NewList()

	configList.AddItem("Projects", "", 0, nil)
	configList.AddItem("Personal Info", "", 0, nil)
	configList.AddItem("Experience", "", 0, nil)

	configList.SetTitle("Config")
	// configList.SetTitleAlign(tview.AlignLeft)

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

	return &ConfigPanel{
		configs: configList,
	}
}

func (self *ConfigPanel) GetPrimitive() tview.Primitive {
	return self.configs
}

func (self *ConfigPanel) GetBox() *tview.Box {
	return self.configs.Box
}
