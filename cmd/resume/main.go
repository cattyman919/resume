package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"resume/internals/builder"
	model "resume/internals/model"
	"resume/internals/parse"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	dataFile = "cv_data.yaml"
)

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

	// Load data from YAML file
	byteValue, err := os.ReadFile(dataFile)
	if err != nil {
		log.Fatalf("Error: Data file not found at %s: %v", dataFile, err)
	}

	var cvData model.CVData
	if err := yaml.Unmarshal(byteValue, &cvData); err != nil {
		log.Fatalf("Error decoding YAML from %s: %v", dataFile, err)
	}

	// Hash Set to get total types in the YAML config
	total_types := make(map[string]struct{})

	var wg sync.WaitGroup
	var mu sync.Mutex

	parse.GetTotalTypes(total_types, cvData, &wg, &mu)

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
