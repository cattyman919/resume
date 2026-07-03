package core

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"sync"
	"text/template"
	"time"

	"github.com/cattyman919/autocv/internal/core/config"
	"github.com/cattyman919/autocv/internal/core/domain"
	"github.com/cattyman919/autocv/internal/core/generator"
	"github.com/lmittmann/tint"
)

type App struct {
	CVConfig *config.CVConfig
	Template *template.Template
}

func initLogger() {
	w := os.Stdout

	slog.SetDefault(slog.New(
		tint.NewHandler(w, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.Kitchen,
		}),
	))
}

func NewApp() (*App, error) {
	initLogger()

	if _, err := exec.LookPath("typst"); err != nil {
		slog.Error("Program `tyspt` compiler not found in $PATH")
		fmt.Println("")
		fmt.Println("Please install tyspt compiler first before running the program")
		fmt.Println("https://typst.app/open-source/")
		return nil, err
	}

	cvCfg, err := config.NewConfig()
	if err != nil {
		slog.Error("Error Parsing Config", "Err", err)
		log.Fatalln(err)
	}

	tmpl, err := generator.NewTemplate()
	if err != nil {
		slog.Error("Error Creating Template", "Err", err)
		log.Fatalln(err)
	}

	return &App{
		CVConfig: cvCfg,
		Template: tmpl,
	}, nil
}

func (a *App) Run() {

	var wg sync.WaitGroup

	for _, cvType := range a.CVConfig.CVTypesCfg {
		cvData := domain.CVTypeData{
			General:  &a.CVConfig.GeneralCfg,
			Settings: &a.CVConfig.SettingsCfg,
			CVType:   cvType,
		}

		wg.Go(func() {
			err := generator.GenerateCVType(&cvData, a.Template)
			if err != nil {
				slog.Error("Error generating CV Type", "err", err)
			}
			err = generator.GeneratePDF(&cvData)
			if err != nil {
				slog.Error("Error generating PDF CV Type", "err", err)
			}
		})

	}

	wg.Wait()

}
