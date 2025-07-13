package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"sync"
)

// --- Configuration ---
const (
	dataFile          = "cv_data.json"
	mainCVSectionsDir = "main_cv/sections"
	bwCVSectionsDir   = "bw_cv/sections"
)

// --- Data Structures for JSON Unmarshaling ---
type PersonalInfo struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	Website        string `json:"website"`
	Linkedin       string `json:"linkedin"`
	LinkedinHandle string `json:"linkedin_handle"`
	Github         string `json:"github"`
	GithubHandle   string `json:"github_handle"`
	ProfilePic     string `json:"profile_pic"`
	Location       string `json:"location"`
}

type Experience struct {
	Role     string   `json:"role"`
	Type     string   `json:"type"`
	Company  string   `json:"company"`
	Location string   `json:"location"`
	Dates    string   `json:"dates"`
	Points   []string `json:"points"`
	Types    []string `json:"types"`
}

type Education struct {
	Institution string   `json:"institution"`
	Degree      string   `json:"degree"`
	Dates       string   `json:"dates"`
	Gpa         string   `json:"gpa"`
	Details     []string `json:"details"`
}

type Award struct {
	Title        string   `json:"title"`
	Organization string   `json:"organization"`
	Date         string   `json:"date"`
	Points       []string `json:"points"`
}

type Project struct {
	Name         string   `json:"name"`
	Github       string   `json:"github"`
	GithubHandle string   `json:"github_handle"`
	Points       []string `json:"points"`
	Types        []string `json:"types"`
}

type Certificate struct {
	Name string `json:"name"`
	Year int    `json:"year"`
}

type SkillsAchievements struct {
	HardSkills           []string      `json:"Hard Skills"`
	SoftSkills           []string      `json:"Soft Skills"`
	ProgrammingLanguages []string      `json:"Programming Languages"`
	DatabaseLanguages    []string      `json:"Database Languages"`
	Misc                 []string      `json:"Misc"`
	Certificates         []Certificate `json:"Certificates"`
}

type CVData struct {
	PersonalInfo       PersonalInfo       `json:"personal_info"`
	Experiences        []Experience       `json:"experiences"`
	Education          []Education        `json:"education"`
	Awards             []Award            `json:"awards"`
	Projects           []Project          `json:"projects"`
	SkillsAchievements SkillsAchievements `json:"skills_achievements"`
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
		fmt.Fprintf(os.Stderr, "Converts Json to Latex Documents\n")
		fmt.Fprintf(os.Stderr, "This is a CLI tool for converting the cv_data.json to the approriate LaTex document based on each sections\n")
		fmt.Fprintf(os.Stderr, "it can convert to different types of CV: main, devops, fullstack\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExample:\n")
		fmt.Fprintf(os.Stderr, "  go run main.go -type=fullstack\n")
	}

	flag.Parse()

	fmt.Printf("==== Generating %s CV ====\n", *cvType)
	fmt.Printf("Loading JSON CV...\n")
	// Load data from JSON file
	byteValue, err := os.ReadFile(dataFile)
	if err != nil {
		log.Fatalf("Error: Data file not found at %s: %v", dataFile, err)
	}

	var cvData CVData
	if err := json.Unmarshal(byteValue, &cvData); err != nil {
		log.Fatalf("Error decoding JSON from %s: %v", dataFile, err)
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
