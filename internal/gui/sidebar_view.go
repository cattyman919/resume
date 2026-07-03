package gui

import (
	"cmp"
	"image/color"
	"slices"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/cattyman919/autocv/internal/core/config"
	"github.com/cattyman919/autocv/internal/core/domain"
)

func sidebarView(cvConfig *config.CVConfig) fyne.CanvasObject {
	cvTypes := cvConfig.CVTypesCfg

	topGroup := container.NewVBox(
		canvas.NewText("Global Config", color.White),
		verticalSpace(20), // 20px gap
		widget.NewSeparator(),
	)

	// sort cv types, so we dont get random order everytime
	slices.SortFunc(cvTypes, func(a, b domain.CVType) int {
		return cmp.Compare(a.TypeName, b.TypeName)
	})

	cvTypeList := widget.NewList(
		func() int {
			return len(cvTypes)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(cvTypes[i].TypeName)
		})

	if len(cvTypes) > 0 {
		cvTypeList.Select(0)
	}

	return container.NewBorder(
		topGroup,
		nil,
		nil,
		nil,
		cvTypeList,
	)

}
