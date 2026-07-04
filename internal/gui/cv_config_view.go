package gui

import (
	g "github.com/AllenDang/giu"
	"github.com/cattyman919/autocv/internal/core/domain"
)

func (a *app) cv_config_view() g.Widget {
	var selectedCVConfigWidget g.Widget

	if a.selectedCV != nil {
		layoutLabels := mapNewSlice(a.selectedCV.Layouts, func(item domain.Layout) g.Widget {
			return g.Label(string(item))
		})

		cvLayoutWidget := g.ListClipper().Layout(layoutLabels...)

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
