package gui

import (
	g "github.com/AllenDang/giu"
	"github.com/cattyman919/autocv/internal/core/domain"
	"github.com/cattyman919/autocv/internal/gui/widget"
)

func (a *app) cv_config_view() g.Widget {
	var selectedCVConfigWidget g.Widget

	if a.selectedCV != nil {
		layoutLabels := mapNewSlice(a.selectedCV.Layouts, func(item domain.Layout) g.Widget {
			return g.Label(string(item))
		})

		cvLayoutWidget := widget.ReorderableList("cv-layouts", a.selectedCV.Layouts).
			Layout(layoutLabels...).
			OnReorder(func(newLayouts []domain.Layout) {
				a.selectedCV.Layouts = newLayouts
			})

		selectedCVConfigWidget = g.ListClipper().Layout(
			g.Label("Layout"),
			g.Separator(),
			cvLayoutWidget,
			g.Separator(),
			g.Label("Experiences"),
		)
	}

	return g.Column(
		g.Label("CV Config"),
		selectedCVConfigWidget,
	)
}
