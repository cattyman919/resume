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
	var cvData = new(model.CVData)
	var err_channel = make(chan error, 3)

	var read_unmarshal_file = func(filePath string, data interface{}) {
		defer wg.Done()
		byteValue, err := os.ReadFile(filePath)
		if err != nil {
			err_channel <- err
		}

		if err := yaml.Unmarshal(byteValue, data); err != nil {
			err_channel <- err
		}
	}

	wg.Add(3)
	go read_unmarshal_file(CV_GENERAL_DATA_PATH, &cvData.General)
	go read_unmarshal_file(CV_PROJECTS_DATA_PATH, &cvData.Projects)
	go read_unmarshal_file(CV_EXPERIENCES_DATA_PATH, &cvData.Experiences)
	wg.Wait()

	close(err_channel)

	for err := range err_channel {
		if err != nil {
			return model.CVData{}, err
		}
	}

	return *cvData, nil
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

	parse.GetTotalTypes(total_types, &cvData.Experiences, &cvData.Projects, &wg, &mu)

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
