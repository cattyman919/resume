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
	CVConfig    *config.CVConfig
	Template    *template.Template
	TypstAvailable bool
	mu         sync.Mutex
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

	typstAvailable := true
	if _, err := exec.LookPath("typst"); err != nil {
		slog.Warn("Typst compiler not found in $PATH — PDF generation will be disabled")
		typstAvailable = false
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
		CVConfig:       cvCfg,
		Template:       tmpl,
		TypstAvailable: typstAvailable,
	}, nil
}

func (a *App) GenerateAllCVs() {
	var wg sync.WaitGroup

	for _, cvType := range a.CVConfig.CVTypesCfg {
		cvData := domain.CVTypeData{
			General:  &a.CVConfig.GeneralCfg,
			Settings: &a.CVConfig.SettingsCfg,
			CVType:   cvType,
		}

		wg.Go(func() {
			err := generator.GenerateCVTypeTemplate(&cvData, a.Template)
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

func (a *App) GenerateCV(cvType *domain.CVType) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	cvData := domain.CVTypeData{
		General:  &a.CVConfig.GeneralCfg,
		Settings: &a.CVConfig.SettingsCfg,
		CVType:   *cvType,
	}
	err := generator.GenerateCVTypeTemplate(&cvData, a.Template)
	if err != nil {
		slog.Error("Error generating CV Type", "err", err)
		return err
	}

	err = generator.GeneratePDF(&cvData)
	if err != nil {
		slog.Error("Error generating PDF CV Type", "err", err)
		return err
	}

	return nil
}

func (a *App) ReloadConfig() error {
	cvCfg, err := config.NewConfig()
	if err != nil {
		return fmt.Errorf("error reloading config: %w", err)
	}
	a.CVConfig = cvCfg
	return nil
}

func (a *App) GetCVType(typeName string) *domain.CVType {
	for i := range a.CVConfig.CVTypesCfg {
		if a.CVConfig.CVTypesCfg[i].TypeName == typeName {
			return &a.CVConfig.CVTypesCfg[i]
		}
	}
	return nil
}

func (a *App) GetCVTypeNames() []string {
	names := make([]string, 0, len(a.CVConfig.CVTypesCfg))
	for _, cvType := range a.CVConfig.CVTypesCfg {
		names = append(names, cvType.TypeName)
	}
	return names
}

func (a *App) GetPDFPath(typeName string) string {
	return fmt.Sprintf("%s - CV (%s).pdf", a.CVConfig.GeneralCfg.PersonalInfo.Name, typeName)
}

func (a *App) SaveGeneralConfig() error {
	return config.SaveGeneralConfig(&a.CVConfig.GeneralCfg)
}

func (a *App) SaveSettingsConfig() error {
	return config.SaveSettingsConfig(&a.CVConfig.SettingsCfg)
}

func (a *App) SaveCVTypeConfig(cvType *domain.CVType) error {
	return config.SaveCVTypeConfig(cvType)
}

func (a *App) DeleteCVType(typeName string) error {
	return config.DeleteCVType(typeName)
}

func (a *App) RenameCVType(oldName, newName string) error {
	return config.RenameCVType(oldName, newName)
}

func (a *App) CreateCVType(typeName string, cloneFrom string) error {
	return config.CreateCVType(typeName, cloneFrom)
}
