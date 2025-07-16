package main

import (
	"fmt"
	"log"
	"os"
	"resume/internals"
	"sync"

	"gopkg.in/yaml.v3"
)

// --- Configuration ---
const (
	dataFile = "cv_data.yaml"
)

func main() {

	fmt.Printf("\n==== Generating All LaTeX CV ====\n")

	fmt.Printf("Loading YAML Data...\n")

	// Load data from YAML file
	byteValue, err := os.ReadFile(dataFile)
	if err != nil {
		log.Fatalf("Error: Data file not found at %s: %v", dataFile, err)
	}

	var cvData internals.CVData
	if err := yaml.Unmarshal(byteValue, &cvData); err != nil {
		log.Fatalf("Error decoding YAML from %s: %v", dataFile, err)
	}

	// Hash Set to get total types in the YAML config
	total_types := make(map[string]struct{})

	// O(N^2) searches each experience and iterate through the type
	for _, exp := range cvData.Experiences {
		if len(exp.Types) == 0 {
			continue
		}

		for _, cvType := range exp.Types {
			_, exist := total_types[cvType]
			if !exist {
				total_types[cvType] = struct{}{}
			}

		}
	}

	// O(N^2) searches each projects and iterate through the type
	// TODO : use concurrency
	for _, proj := range cvData.Projects {
		if len(proj.Types) == 0 {
			continue
		}

		for _, cvType := range proj.Types {
			_, exist := total_types[cvType]
			if !exist {
				total_types[cvType] = struct{}{}
			}

		}
	}
	var wg sync.WaitGroup
	wg.Add(len(total_types))

	for cvType := range total_types {
		go internals.Write_CV(cvType, cvData, &wg)
	}

	wg.Wait()

	fmt.Printf("==== All LaTeX CV Generation Complete ====\n")
}
