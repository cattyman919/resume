package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cattyman919/autocv/internal/core/domain"
	yaml "github.com/goccy/go-yaml"
)

const configPath = "config"

func SaveGeneralConfig(general *domain.CVGeneral) error {
	data, err := yaml.Marshal(general)
	if err != nil {
		return fmt.Errorf("failed to marshal general config: %w", err)
	}

	rawContent, err := os.ReadFile(filepath.Join(configPath, "general.yaml"))
	if err != nil {
		return fmt.Errorf("failed to read raw general.yaml: %w", err)
	}

	preserved := RestoreEnvTemplates(string(data), string(rawContent))

	err = os.WriteFile(filepath.Join(configPath, "general.yaml"), []byte(preserved), 0644)
	if err != nil {
		return fmt.Errorf("failed to write general config: %w", err)
	}

	return nil
}

func SaveSettingsConfig(settings *domain.CVSettings) error {
	data, err := yaml.Marshal(settings)
	if err != nil {
		return fmt.Errorf("failed to marshal settings config: %w", err)
	}

	err = os.WriteFile(filepath.Join(configPath, "settings.yaml"), data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write settings config: %w", err)
	}

	return nil
}

func SaveCVTypeConfig(cvType *domain.CVType) error {
	data, err := yaml.Marshal(cvType)
	if err != nil {
		return fmt.Errorf("failed to marshal cv type config: %w", err)
	}

	typePath := filepath.Join(configPath, "types", cvType.TypeName+".yaml")
	err = os.WriteFile(typePath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write cv type config: %w", err)
	}

	return nil
}

func DeleteCVType(typeName string) error {
	typePath := filepath.Join(configPath, "types", typeName+".yaml")
	if _, err := os.Stat(typePath); os.IsNotExist(err) {
		return fmt.Errorf("cv type %s does not exist", typeName)
	}

	err := os.Remove(typePath)
	if err != nil {
		return fmt.Errorf("failed to delete cv type %s: %w", typeName, err)
	}

	return nil
}

func RenameCVType(oldName, newName string) error {
	oldPath := filepath.Join(configPath, "types", oldName+".yaml")
	newPath := filepath.Join(configPath, "types", newName+".yaml")

	if _, err := os.Stat(oldPath); os.IsNotExist(err) {
		return fmt.Errorf("cv type %s does not exist", oldName)
	}

	if _, err := os.Stat(newPath); err == nil {
		return fmt.Errorf("cv type %s already exists", newName)
	}

	data, err := os.ReadFile(oldPath)
	if err != nil {
		return fmt.Errorf("failed to read cv type %s: %w", oldName, err)
	}

	err = os.WriteFile(newPath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write cv type %s: %w", newName, err)
	}

	err = os.Remove(oldPath)
	if err != nil {
		return fmt.Errorf("failed to delete old cv type %s: %w", oldName, err)
	}

	return nil
}

func CreateCVType(typeName string, cloneFrom string) error {
	typeName = strings.TrimSpace(typeName)
	if typeName == "" {
		return fmt.Errorf("type name cannot be empty")
	}

	newPath := filepath.Join(configPath, "types", typeName+".yaml")

	if _, err := os.Stat(newPath); err == nil {
		return fmt.Errorf("cv type %s already exists", typeName)
	}

	var data []byte
	var err error

	if cloneFrom != "" {
		clonePath := filepath.Join(configPath, "types", cloneFrom+".yaml")
		data, err = os.ReadFile(clonePath)
		if err != nil {
			return fmt.Errorf("failed to read clone source %s: %w", cloneFrom, err)
		}
	} else {
		data = []byte("layout:\n  - description\n  - educations\n  - skills\n  - projects\n")
	}

	err = os.WriteFile(newPath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to create cv type %s: %w", typeName, err)
	}

	return nil
}
