package main

import (
	"flag"
	"fmt"
	"os"
	"resume/internals/config"
	"resume/internals/generator"
	"sync"
	"time"
)

func main() {
	debug := flag.Bool("debug", false, "Enable debug output for pdflatex")
	benchmark := flag.Bool("benchmark", false, "Run benchmark for the main process")
	flag.Parse()

	var startTime time.Time
	if *benchmark {
		startTime = time.Now()
	}

	generator.DebugMode = *debug

	fmt.Printf("\n==== Generating All LaTeX CV ====\n")

	var wg sync.WaitGroup

	fmt.Printf("Loading YAML Data...\n")

	cvData, err := config.LoadYAMLData(&wg)
	if err != nil {
		fmt.Printf("Error loading YAML data: %v\n", err)
		return
	}

	fmt.Printf("YAML Data Loaded Successfully!\n")

	// Hash Set to get total types in the YAML config
	total_types := make(map[string]struct{})

	var mu sync.Mutex

	config.GetTotalTypes(total_types, &cvData.Experiences, &cvData.Projects, &wg, &mu)

	wg.Add(len(total_types))

	err = os.MkdirAll("out", generator.FolderPermission)
	if err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
	}

	// Generate each CV type concurrently
	for cvType := range total_types {
		go generator.Write_CV(cvType, cvData, &wg)
	}

	wg.Wait()

	generator.MoveAuxFiles()

	fmt.Printf("==== All LaTeX CV Generation Complete ====\n")

	if *benchmark {
		elapsed := time.Since(startTime)
		fmt.Printf("\n==== Benchmark Result ====\n")
		fmt.Printf("Total execution time: %s\n", elapsed)
	}
}
