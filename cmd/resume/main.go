package main

import (
	"flag"
	"log"
	"os"
	"resume/internals/config"
	"resume/internals/generator"
	"resume/internals/utils"
	"sync"
)

func main() {
	debug := flag.Bool("debug", false, "Enable debug output for pdflatex")
	flag.Parse()

	generator.DebugMode = *debug

	log.Printf("\n==== Generating All LaTeX CV ====\n")

	log.Printf("Loading YAML Data...\n")

	cvData, err := config.LoadYAMLData()
	if err != nil {
		log.Printf("Error loading YAML data: %v\n", err)
		return
	}

	log.Printf("YAML Data Loaded Successfully!\n")

	// hash set containing all the unique CV Types
	total_types := config.GetTotalTypes(&cvData.Experiences, &cvData.Projects)

	err = os.MkdirAll("out", generator.FOLDER_PERMISSION)
	if err != nil {
		log.Printf("Error creating output directory: %v\n", err)
	}

	var wg sync.WaitGroup

	tmpl, err := generator.GetTemplate()
	if err != nil {
		log.Fatalln(err)
	}
	// Generate each CV type concurrently
	for cvType := range total_types {
		// Helper function to do wg.Add(1) and wg.Done()| like wg.Go in the newer Go version
		utils.Go(
			&wg,
			func() { generator.Write_CV(cvType, cvData, tmpl) },
		)
	}

	wg.Wait()

	generator.MoveAuxFiles()

	log.Printf("==== All LaTeX CV Generation Complete ====\n")
}
