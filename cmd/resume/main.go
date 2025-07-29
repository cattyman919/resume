package main

import (
	"flag"
	"fmt"
	"os"
	"resume/internals/builder"
	model "resume/internals/model"
	"resume/internals/parse"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	CV_GENERAL_DATA_PATH     = "config/general.yaml"
	CV_PROJECTS_DATA_PATH    = "config/projects.yaml"
	CV_EXPERIENCES_DATA_PATH = "config/experiences.yaml"
)

func loadYAMLData(wg *sync.WaitGroup) (model.CVData, error) {
	var cvData_general = new(model.CV_General)
	var cvData_projects = new(model.CV_Projects)
	var cvData_experiences = new(model.CV_Experiences)
	var err_channel = make(chan error, 3)

	wg.Add(3)
	go func() {
		defer wg.Done()
		byteValue, err := os.ReadFile(CV_GENERAL_DATA_PATH)
		if err != nil {
			err_channel <- err
			return
		}

		if err := yaml.Unmarshal(byteValue, cvData_general); err != nil {
			err_channel <- err
			return
		}
	}()
	go func() {
		defer wg.Done()
		byteValue, err := os.ReadFile(CV_EXPERIENCES_DATA_PATH)
		if err != nil {
			err_channel <- err
			return
		}

		if err := yaml.Unmarshal(byteValue, cvData_experiences); err != nil {
			err_channel <- err
			return
		}
	}()
	go func() {
		defer wg.Done()
		byteValue, err := os.ReadFile(CV_PROJECTS_DATA_PATH)
		if err != nil {
			err_channel <- err
			return
		}

		if err := yaml.Unmarshal(byteValue, cvData_projects); err != nil {
			err_channel <- err
			return
		}
	}()
	wg.Wait()

	close(err_channel)

	for err := range err_channel {
		if err != nil {
			return model.CVData{}, err
		}
	}

	return model.CVData{
		General:     cvData_general,
		Projects:    *cvData_projects,
		Experiences: *cvData_experiences,
	}, nil
}

func main() {
	debug := flag.Bool("debug", false, "Enable debug output for pdflatex")
	benchmark := flag.Bool("benchmark", false, "Run benchmark for the main process")
	flag.Parse()

	var startTime time.Time
	if *benchmark {
		startTime = time.Now()
	}

	builder.DebugMode = *debug

	fmt.Printf("\n==== Generating All LaTeX CV ====\n")

	fmt.Printf("Loading YAML Data...\n")
	var wg sync.WaitGroup

	cvData, err := loadYAMLData(&wg)
	if err != nil {
		fmt.Printf("Error loading YAML data: %v\n", err)
		return
	}

	fmt.Printf("YAML Data Loaded Successfully!\n")

	// Hash Set to get total types in the YAML config
	total_types := make(map[string]struct{})

	var mu sync.Mutex

	parse.GetTotalTypes(total_types, cvData.Experiences, cvData.Projects, &wg, &mu)

	wg.Add(len(total_types))

	for cvType := range total_types {
		go builder.Write_CV(cvType, cvData, &wg)
	}

	wg.Wait()

	builder.MoveAuxFiles()

	fmt.Printf("==== All LaTeX CV Generation Complete ====\n")

	if *benchmark {
		elapsed := time.Since(startTime)
		fmt.Printf("\n==== Benchmark Result ====\n")
		fmt.Printf("Total execution time: %s\n", elapsed)
	}
}
