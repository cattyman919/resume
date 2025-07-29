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
	"sync"
)

var DebugMode bool = false

const (
	name              = "Seno Pamungkas Rahman"
	mainCVSectionsDir = "main_cv/sections"
	bwCVSectionsDir   = "bw_cv/sections"
	folderPermission  = 0755
	filePermission    = 0644
)

func writeTexFile(path string, content string, wg *sync.WaitGroup, outputChan chan string) {
	defer wg.Done()

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, folderPermission); err != nil {
		log.Fatalf("Error creating directory %s: %v", dir, err)
	}
	if err := os.WriteFile(path, []byte(content), filePermission); err != nil {
		log.Fatalf("Error writing %s: %v", path, err)
	}

	outputChan <- fmt.Sprintf("Written %s", path)
}

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
func copyDir(src, dst string, wg *sync.WaitGroup) error {
	defer wg.Done()

	// Create the destination directory with the same permissions
	err := os.MkdirAll(dst, folderPermission)
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
			return os.MkdirAll(destPath, folderPermission)
		} else {
			// Copy file
			return copyFile(path, destPath)
		}
	})
}

func Write_CV(cvType string, cvData resume.CVData, wgParent *sync.WaitGroup) {
	defer wgParent.Done()

	fmt.Printf("Generating CV (%s)\n", cvType)

	mainPath := "cv/" + cvType + "/main_cv"
	bwPath := "cv/" + cvType + "/bw_cv"

	var wg sync.WaitGroup
	wg.Add(2)

	go copyDir("template_cv/main_cv", mainPath, &wg)
	go copyDir("template_cv/bw_cv", bwPath, &wg)

	// Filter experiences and projects based on cvType
	var filteredExperiences []resume.Experience
	for _, exp := range cvData.Experiences {

		// if types is empty, dont add it at all to the CV
		if len(exp.JobType) == 0 {
			continue
		}

		if slices.Contains(exp.CVType, cvType) {
			filteredExperiences = append(filteredExperiences, exp)
		}
	}

	// If type is not empty in expereinces, the use the filtered experiences. else use the original experiences
	if filteredExperiences != nil {
		cvData.Experiences = filteredExperiences
	}

	var filteredProjects []resume.Project
	for _, proj := range cvData.Projects {

		// if types is empty, dont add it at all to the CV
		if len(proj.CVType) == 0 {
			continue
		}

		if slices.Contains(proj.CVType, cvType) {
			filteredProjects = append(filteredProjects, proj)
		}
	}

	// If type is not empty in projects, the use the filtered projects. else use the original projects
	if filteredProjects != nil {
		cvData.Projects = filteredProjects
	}

	type generate_section func(*resume.CVData) string

	// maps section to main and bw functions
	sections := map[string][2]generate_section{
		"Header.tex":             {GenerateHeaderMainCV, GenerateHeaderBwCV},
		"Experience.tex":         {GenerateExperienceMainCV, GenerateExperienceBwCV},
		"Education.tex":          {GenerateEducationMainCV, GenerateEducationBwCV},
		"Awards.tex":             {GenerateAwardsMainCV, GenerateAwardsBwCV},
		"Projects.tex":           {GenerateProjectsMainCV, GenerateProjectsBwCV},
		"Achivements_Skills.tex": {GenerateSkillsMainCV, GenerateSkillsBwCV},
	}

	// Waits for the copying files from template_cv to the specific cv type directory
	// preamble.sty and main.tex
	wg.Wait()

	// 6 sections * 2 (main + bw)
	total_io_ops := len(sections) * 2 // total of 12 goroutine function for creating all sections files

	wg.Add(total_io_ops)

	outputChan := make(chan string, total_io_ops)

	for section, function := range sections {
		mainPath := filepath.Join("cv", cvType, mainCVSectionsDir, section)
		bwPath := filepath.Join("cv", cvType, bwCVSectionsDir, section)

		go writeTexFile(mainPath, function[0](&cvData), &wg, outputChan)
		go writeTexFile(bwPath, function[1](&cvData), &wg, outputChan)
	}

	go func() {
		wg.Wait()
		close(outputChan)
	}()

	// loops until finally chanel is closed after all goroutines finished
	for s := range outputChan {
		fmt.Println(s)
	}

	wg.Add(2)
	go write_PDF(&cvType, &cvData.General.PersonalInfo.Name, "main", &wg)
	go write_PDF(&cvType, &cvData.General.PersonalInfo.Name, "bw", &wg)
	wg.Wait()
}

func write_PDF(cvType *string, name *string, type_bw_main string, wg *sync.WaitGroup) {
	defer wg.Done()

	target_pdf := fmt.Sprintf("%s - CV (%s) (%s)", *name, *cvType, strings.ToUpper(type_bw_main))

	program := "pdflatex"
	args := []string{
		"-output-directory=../../../out",
		"-output-format=pdf",
		"-jobname",
		target_pdf,
		"main.tex",
	}

	cmd := exec.Command(program, args...)

	// Generate either the main or black and white CV PDF
	switch type_bw_main {
	case "main":
		cmd.Dir = fmt.Sprintf("cv/%s/main_cv", *cvType)
	case "bw":
		cmd.Dir = fmt.Sprintf("cv/%s/bw_cv", *cvType)
	default:
		log.Fatalf("Error Handling Writing PDF for type %s", type_bw_main)
	}

	if DebugMode {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	// Blocking function
	err := cmd.Run()

	if err != nil {
		log.Fatalf("Command failed with error: %v", err)
	}

	fmt.Printf("Generated %s.pdf\n", target_pdf)
}

func MoveAuxFiles() {
	auxDir := filepath.Join("out", "aux")

	if err := os.MkdirAll(auxDir, folderPermission); err != nil {
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
