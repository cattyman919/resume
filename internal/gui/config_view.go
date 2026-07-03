package gui

import (
	"fyne.io/fyne/v2"
	"github.com/cattyman919/autocv/internal/core/config"
)

type configViewState string

const (
	configViewStateSetting configViewState = "settings"
	configViewStateGeneral configViewState = "general"
	configViewStateCVType  configViewState = "cvType"
)

func configView(viewState configViewState, cvConfig *config.CVConfig) fyne.CanvasObject {

	var content fyne.CanvasObject

	switch viewState {
	case configViewStateSetting:
	case configViewStateGeneral:
	case configViewStateCVType:
	}

	return content
}
