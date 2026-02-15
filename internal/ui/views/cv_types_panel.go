package views

import "github.com/rivo/tview"

type CVTypesPanel struct {
	cvTypes *tview.List
}

func NewCvTypesPanel() *CVTypesPanel {
	cvTypes := tview.NewList()
	cvTypes.SetTitle("CV Types")
	cvTypes.SetBorder(true)

	cvTypes.AddItem("Default", "", 0, nil)
	cvTypes.AddItem("Devops", "", 0, nil)
	cvTypes.AddItem("Frontend", "", 0, nil)

	return &CVTypesPanel{
		cvTypes: cvTypes,
	}
}

func (self *CVTypesPanel) GetPrimitive() tview.Primitive {
	return self.cvTypes
}

func (self *CVTypesPanel) GetBox() *tview.Box {
	return self.cvTypes.Box
}
