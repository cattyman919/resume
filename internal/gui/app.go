package gui

import (
	"cmp"
	"image/color"
	"slices"

	g "github.com/AllenDang/giu"
	"github.com/cattyman919/autocv/internal/core"
	"github.com/cattyman919/autocv/internal/core/domain"
)

var (
	background = color.NRGBA{R: 0xC0, G: 0xC0, B: 0xC0, A: 0xFF}
	red        = color.NRGBA{R: 0xC0, G: 0x40, B: 0x40, A: 0xFF}
	green      = color.NRGBA{R: 0x40, G: 0xC0, B: 0x40, A: 0xFF}
	blue       = color.NRGBA{R: 0x40, G: 0x40, B: 0xC0, A: 0xFF}
)

type app struct {
	coreApp    *core.App
	window     *g.MasterWindow
	selectedCV *domain.CVType
}

func NewApp() (*app, error) {
	coreApp, err := core.NewApp()
	if err != nil {
		return nil, err
	}

	slices.SortFunc(coreApp.CVConfig.CVTypesCfg, func(a, b domain.CVType) int {
		return cmp.Compare(a.TypeName, b.TypeName)
	})

	wnd := g.NewMasterWindow("AutoCV", 1200, 800, g.MasterWindowFlagsNotResizable)
	wnd.SetTargetFPS(30)

	var selectedCV *domain.CVType
	if len(coreApp.CVConfig.CVTypesCfg) > 0 {
		selectedCV = &coreApp.CVConfig.CVTypesCfg[0]
	}

	return &app{
		coreApp:    coreApp,
		window:     wnd,
		selectedCV: selectedCV,
	}, nil
}

var (
	mainContentSplitRatio float32 = 0.5
)

func (a *app) Run() {
	a.window.Run(func() {
		g.SingleWindow().Layout(
			g.SplitLayout(g.DirectionVertical, &sidebarWidthPercentage,
				a.sidebar_view(),
				g.SplitLayout(g.DirectionVertical, &mainContentSplitRatio,
					a.cv_config_view(),
					g.Label("PDF viewer"),
				).SplitRefType(g.SplitRefProc),
			).SplitRefType(g.SplitRefProc),
		)
	})
}
