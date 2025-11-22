package generator

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"resume/internals/resume"
	"slices"
	"strings"
	"text/template"
)

var DebugMode bool = false

const (
	FOLDER_PERMISSION = 0755
)

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()
	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destinationFile.Close()
	_, err = io.Copy(destinationFile, sourceFile)
	return err
}

// copyDir recursively copies a directory from src to dst
func copyDir(src, dst string) error {

	// Create the destination directory with the same permissions
	err := os.MkdirAll(dst, FOLDER_PERMISSION)
	if err != nil {
		return err
	}

	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err // Propagate errors
		}

		// Find the relative path of the item from the source directory
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		// Create the full destination path
		destPath := filepath.Join(dst, relPath)

		if d.IsDir() {
			// Create directory at destination
			return os.MkdirAll(destPath, FOLDER_PERMISSION)
		} else {
			// Copy file
			return copyFile(path, destPath)
		}
	})
}

func Write_CV(cvType string, cvData resume.CVData, tmpl *template.Template) {

	log.Printf("Generating CV (%s)\n", cvType)

	cvTemplatePath := filepath.Join("template_cv", "bw_cv")
	cvOutputPath := filepath.Join("build_cv", cvType, "bw_cv")

	copyDir(cvTemplatePath, cvOutputPath)

	// The filtered functions modifies the cvData
	filterExperiences(&cvType, &cvData.Experiences)
	filterProjects(&cvType, &cvData.Projects)

	filenames := []string{
		"Achivements_Skills.tex",
		"Awards.tex",
		"Education.tex",
		"Experience.tex",
		"Header.tex",
		"Projects.tex",
	}

	for _, templateName := range filenames {
		outputPath := filepath.Join(cvOutputPath, "sections", templateName)
		dirOutputPath := filepath.Dir(outputPath)
		err := os.MkdirAll(dirOutputPath, FOLDER_PERMISSION)
		if err != nil {
			log.Printf("Error creating Directory %s: %v\n", dirOutputPath, err)
			continue
		}

		file, err := os.Create(outputPath)
		if err != nil {
			log.Printf("Error creating file %s: %v", outputPath, err)
			continue
		}

		// We use ExecuteTemplate, passing the base name.
		err = tmpl.ExecuteTemplate(file, templateName, cvData)
		if err != nil {
			log.Printf("Error executing template %s: %v", templateName, err)
			file.Close()
			continue
		}

		file.Close()

		log.Printf("Generated: %s\n", outputPath)
	}

	write_PDF(&cvType, &cvData.General.PersonalInfo.Name)
}

// Runs the 'pdflatex' external command
func write_PDF(cvType *string, name *string) {

	target_pdf := fmt.Sprintf("%s - CV (%s)", *name, *cvType)

	program := "pdflatex"
	args := []string{
		"-output-directory=../../../out",
		"-output-format=pdf",
		"-jobname",
		target_pdf,
		"main.tex",
	}

	cmd := exec.Command(program, args...)
	cmd.Dir = fmt.Sprintf("build_cv/%s/bw_cv", *cvType)

	if DebugMode {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	log.Printf("Generating CV %s.pdf\n", target_pdf)

	if err := cmd.Run(); err != nil {
		log.Fatalf("Command failed with error: %v", err)
	}

	log.Printf("Generated CV %s.pdf\n", target_pdf)
}

func MoveAuxFiles() {
	auxDir := filepath.Join("out", "aux")

	if err := os.MkdirAll(auxDir, FOLDER_PERMISSION); err != nil {
		log.Fatalf("Error creating aux directory: %v", err)
	}

	files, err := os.ReadDir("out")
	if err != nil {
		log.Fatalf("Error reading output directory: %v", err)
	}

	for _, file := range files {
		fileName := file.Name()
		if strings.HasSuffix(fileName, ".log") || strings.HasSuffix(fileName, ".aux") {
			oldPath := filepath.Join("out", fileName)
			newPath := filepath.Join("out", "aux", fileName)
			if err := os.Rename(oldPath, newPath); err != nil {
				log.Printf("Error moving file %s: %v", fileName, err)
			}
		}
	}
}

func filterExperiences(cvType *string, experiences *[]resume.Experience) {
	var filtered []resume.Experience

	for _, exp := range *experiences {
		// if types is empty, dont add it at all to the CV
		if len(exp.CVType) == 0 {
			continue
		}

		if slices.Contains(exp.CVType, *cvType) {
			filtered = append(filtered, exp)
		}
	}

	// If type is not empty in expereinces, the use the filtered experiences. else use the original experiences
	if filtered != nil {
		*experiences = filtered
	}

}

func filterProjects(cvType *string, projects *[]resume.Project) {
	var filtered []resume.Project

	for _, exp := range *projects {
		// if types is empty, dont add it at all to the CV
		if len(exp.CVType) == 0 {
			continue
		}

		if slices.Contains(exp.CVType, *cvType) {
			filtered = append(filtered, exp)
		}
	}

	// If type is not empty in expereinces, the use the filtered experiences. else use the original experiences
	if filtered != nil {
		*projects = filtered
	}

}
