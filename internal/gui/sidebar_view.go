package gui

import (
	g "github.com/AllenDang/giu"
	"github.com/cattyman919/autocv/internal/core/domain"
)

var (
	sidebarWidthPercentage float32 = 0.2
	cvTypeListHeight       float32 = 600
)

var (
	PopupCreateNewCV string = "Create New CV"
)

func (a *app) sidebar_view() g.Widget {
	cvConfig := a.coreApp.CVConfig

	items := mapNewSlice(cvConfig.CVTypesCfg, func(item domain.CVType) string {
		return item.TypeName
	})

	listbox := g.ListBox(items).Size(-1, cvTypeListHeight).OnChange(func(selectedIndex int) {
		a.selectedCV = &cvConfig.CVTypesCfg[selectedIndex]
	})

	w, _ := g.GetAvailableRegion()

	popupNewCV := g.PopupModal(PopupCreateNewCV).Layout(
		g.InputText(toPointer("")).Size(w),
		g.Row(
			g.Button("Cancel").OnClick(func() {
				g.CloseCurrentPopup()
			}),
			g.Button("Create").OnClick(func() {
				g.CloseCurrentPopup()
			}),
		),
	)

	return g.Column(
		g.Label("Config"),
		listbox,
		g.Button("+").Size(-1, 40).OnClick(func() {
			g.OpenPopup(PopupCreateNewCV)
		}),
		popupNewCV,
	)
}
