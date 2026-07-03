package generator

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/ledongthuc/pdf"

	"github.com/cattyman919/autocv/internal/core/domain"
	"github.com/cattyman919/autocv/internal/core/utils"
)

const (
	FOLDER_PERMISSION = 0755
)

func GenerateCVType(cvData *domain.CVTypeData, tmpl *template.Template) error {

	slog.Info("Generating CV", "Type", cvData.Type)

	cvTemplatePath := filepath.Join("src")
	cvOutputPath := filepath.Join("build_cv", cvData.Type)

	utils.CopyDir(cvTemplatePath, cvOutputPath)

	templates := []string{
		"main.typ",
		"sections.typ",
	}

	for _, templateName := range templates {
		outputPath := filepath.Join(cvOutputPath, templateName)
		dirOutputPath := filepath.Dir(outputPath)
		err := os.MkdirAll(dirOutputPath, FOLDER_PERMISSION)
		if err != nil {
			slog.Error("Error creating Directory", "Path", dirOutputPath, "Err", err)
			continue
		}

		file, err := os.Create(outputPath)
		if err != nil {
			slog.Error("Error creating file", "file", outputPath, "Err", err)
			continue
		}

		// We use ExecuteTemplate, passing the base name.
		err = tmpl.ExecuteTemplate(file, templateName, cvData)
		if err != nil {
			slog.Error("Error executing template", "Template", templateName, "Err", err)
			file.Close()
			continue
		}

		file.Close()
	}

	return nil
}

func GeneratePDF(cvData *domain.CVTypeData) error {

	targetPDF := fmt.Sprintf("%s - CV (%s).pdf", cvData.General.PersonalInfo.Name, cvData.Type)

	cvPath := filepath.Join("build_cv", cvData.Type, "main.typ")
	outputPath := filepath.Join("out", targetPDF)

	if err := os.Mkdir("out", 0755); err != nil && !strings.Contains(err.Error(), "exists") {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	program := "typst"
	args := []string{
		"compile",
		cvPath,
		"--root",
		".",
		outputPath,
	}

	cmd := exec.Command(program, args...)

	// DEBUG
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr

	output, err := cmd.CombinedOutput()
	if err != nil {
		// Now we include the actual output from typst in your error message
		return fmt.Errorf("failed to execute '%s': %w\nTypst Output:\n%s", program, err, string(output))
	}

	f, r, err := pdf.Open(outputPath)
	if err != nil {
		slog.Warn("Failed to open PDF")
		return nil
	}
	defer f.Close()

	totalPages := r.NumPage()

	// check if its defined or not first
	if cvData.CVType.MaxPages != nil && totalPages > *cvData.CVType.MaxPages {
		slog.Warn(
			"PDF has more pages than max value",
			"CV",
			cvData.CVType.Type,
			"Total Pages",
			totalPages,
			"Max Pages Value",
			*cvData.CVType.MaxPages,
		)
	}

	slog.Info("Generated CV", "Output", outputPath)

	return nil
}
