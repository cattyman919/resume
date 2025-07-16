package main

import (
	"fmt"
	"log"
	"os"
	"resume/internals"
	"sync"

	"gopkg.in/yaml.v3"
)

const (
	dataFile = "cv_data.yaml"
)

func getTotalTypes[T internals.Experience | internals.Project](total_types map[string]struct{}, items []T, wg *sync.WaitGroup, mu *sync.Mutex) {
	defer wg.Done()
	for _, item := range items {
		var types []string

		switch v := any(item).(type) {
		case internals.Experience:
			types = v.Types
		case internals.Project:
			types = v.Types
		}

		if len(types) == 0 {
			continue
		}

		for _, cvType := range types {
			mu.Lock()
			if _, exist := total_types[cvType]; !exist {
				total_types[cvType] = struct{}{}
			}
			mu.Unlock()
		}
	}

}

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

	var wg sync.WaitGroup
	var mu sync.Mutex

	wg.Add(2)
	go getTotalTypes(total_types, cvData.Experiences, &wg, &mu)
	go getTotalTypes(total_types, cvData.Projects, &wg, &mu)
	wg.Wait()

	wg.Add(len(total_types))

	for cvType := range total_types {
		go internals.Write_CV(cvType, cvData, &wg)
	}

	wg.Wait()

	internals.MoveAuxFiles()

	fmt.Printf("==== All LaTeX CV Generation Complete ====\n")
}
