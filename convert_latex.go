package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
	"slices"
	"sync"
)

// --- Configuration ---
const (
	dataFile          = "cv_data.yaml"
	mainCVSectionsDir = "main_cv/sections"
	bwCVSectionsDir   = "bw_cv/sections"
)

// --- Data Structures for YAML Unmarshaling ---
type PersonalInfo struct {
	Name           string `yaml:"name"`
	Email          string `yaml:"email"`
	Phone          string `yaml:"phone"`
	Website        string `yaml:"website"`
	Linkedin       string `yaml:"linkedin"`
	LinkedinHandle string `yaml:"linkedin_handle"`
	Github         string `yaml:"github"`
	GithubHandle   string `yaml:"github_handle"`
	ProfilePic     string `yaml:"profile_pic"`
	Location       string `yaml:"location"`
}

type Experience struct {
	Role     string   `yaml:"role"`
	Type     string   `yaml:"type"`
	Company  string   `yaml:"company"`
	Location string   `yaml:"location"`
	Dates    string   `yaml:"dates"`
	Points   []string `yaml:"points"`
	Types    []string `yaml:"types"`
}

type Education struct {
	Institution string   `yaml:"institution"`
	Degree      string   `yaml:"degree"`
	Dates       string   `yaml:"dates"`
	Gpa         string   `yaml:"gpa"`
	Details     []string `yaml:"details"`
}

type Award struct {
	Title        string   `yaml:"title"`
	Organization string   `yaml:"organization"`
	Date         string   `yaml:"date"`
	Points       []string `yaml:"points"`
}

type Project struct {
	Name         string   `yaml:"name"`
	Github       string   `yaml:"github"`
	GithubHandle string   `yaml:"github_handle"`
	Points       []string `yaml:"points"`
	Types        []string `yaml:"types"`
}

type Certificate struct {
	Name string `yaml:"name"`
	Year int    `yaml:"year"`
}

type SkillsAchievements struct {
	HardSkills           []string      `yaml:"Hard Skills"`
	SoftSkills           []string      `yaml:"Soft Skills"`
	ProgrammingLanguages []string      `yaml:"Programming Languages"`
	DatabaseLanguages    []string      `yaml:"Database Languages"`
	Misc                 []string      `yaml:"Misc"`
	Certificates         []Certificate `yaml:"Certificates"`
}

type CVData struct {
	PersonalInfo       PersonalInfo       `yaml:"personal_info"`
	Experiences        []Experience       `yaml:"experiences"`
	Education          []Education        `yaml:"education"`
	Awards             []Award            `yaml:"awards"`
	Projects           []Project          `yaml:"projects"`
	SkillsAchievements SkillsAchievements `yaml:"skills_achievements"`
}

func writeTexFile(path string, content string, wg *sync.WaitGroup) {
	defer wg.Done()

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatalf("Error creating directory %s: %v", dir, err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		log.Fatalf("Error writing %s: %v", path, err)
	}
}

func main() {
	cvType := flag.String("type", "main", "Type of CV to generate (e.g., main, fullstack, devops)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Converts YAML to Latex Documents\n")
		fmt.Fprintf(os.Stderr, "This is a CLI tool for converting the cv_data.yaml to the approriate LaTex document based on each sections\n")
		fmt.Fprintf(os.Stderr, "it can convert to different types of CV: main, devops, fullstack\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExample:\n")
		fmt.Fprintf(os.Stderr, "  go run main.go -type=fullstack\n")
	}

	flag.Parse()

	fmt.Printf("==== Generating %s CV ====\n", *cvType)

	fmt.Printf("Loading YAML Data...\n")
	// Load data from YAML file
	byteValue, err := os.ReadFile(dataFile)
	if err != nil {
		log.Fatalf("Error: Data file not found at %s: %v", dataFile, err)
	}

	var cvData CVData
	if err := yaml.Unmarshal(byteValue, &cvData); err != nil {
		log.Fatalf("Error decoding YAML from %s: %v", dataFile, err)
	}

	// Filter experiences and projects based on cvType
	var filteredExperiences []Experience
	for _, exp := range cvData.Experiences {
		if len(exp.Types) == 0 {
			filteredExperiences = append(filteredExperiences, exp)
			continue
		}
		if slices.Contains(exp.Types, *cvType) {
			filteredExperiences = append(filteredExperiences, exp)
		}
	}

	cvData.Experiences = filteredExperiences

	var filteredProjects []Project
	for _, proj := range cvData.Projects {
		if len(proj.Types) == 0 {
			filteredProjects = append(filteredProjects, proj)
			continue
		}
		if slices.Contains(proj.Types, *cvType) {
			filteredProjects = append(filteredProjects, proj)
		}
	}
	cvData.Projects = filteredProjects

	type generate_section func(CVData) string
	// maps section to main and bw functions
	sections := map[string][2]generate_section{
		"Header.tex":             {generateHeaderMainCV, generateHeaderBwCV},
		"Experience.tex":         {generateExperienceMainCV, generateExperienceBwCV},
		"Education.tex":          {generateEducationMainCV, generateEducationBwCV},
		"Awards.tex":             {generateAwardsMainCV, generateAwardsBwCV},
		"Projects.tex":           {generateProjectsMainCV, generateProjectsBwCV},
		"Achivements_Skills.tex": {generateSkillsMainCV, generateSkillsBwCV},
	}

	var wg sync.WaitGroup

	// 6 sections * 2 (main + bw)
	wg.Add(len(sections) * 2) // There are total of 12 coroutine function

	fmt.Println("Generating LaTeX sections...")
	for section, function := range sections {
		path := struct {
			main string
			bw   string
		}{
			filepath.Join(mainCVSectionsDir, section),
			filepath.Join(bwCVSectionsDir, section),
		}
		go writeTexFile(path.main, function[0](cvData), &wg)
		go writeTexFile(path.bw, function[1](cvData), &wg)
	}

	wg.Wait()

	fmt.Println("==== LaTeX Section Generation Complete ====")
}
