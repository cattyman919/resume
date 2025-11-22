package config

import (
	"log"
	"os"
	"resume/internals/resume"
	"resume/internals/utils"
	"sync"

	"gopkg.in/yaml.v3"
)

const (
	CV_GENERAL_DATA_PATH     = "config/general.yaml"
	CV_PROJECTS_DATA_PATH    = "config/projects.yaml"
	CV_EXPERIENCES_DATA_PATH = "config/experiences.yaml"
)

func LoadYAMLData() (resume.CVData, error) {
	var wg sync.WaitGroup

	var cvData = new(resume.CVData)
	var err_channel = make(chan error, 3)

	var read_unmarshal_file = func(filePath string, dataObject any) {
		log.Printf("Reading and parsing config file %s\n", filePath)
		byteValue, err := os.ReadFile(filePath)
		if err != nil {
			err_channel <- err
		}

		if err := yaml.Unmarshal(byteValue, dataObject); err != nil {
			err_channel <- err
		}
	}

	utils.Go(&wg, func() {
		read_unmarshal_file(CV_GENERAL_DATA_PATH, &cvData.General)
	})
	utils.Go(&wg, func() {
		read_unmarshal_file(CV_PROJECTS_DATA_PATH, &cvData.Projects)
	})
	utils.Go(&wg, func() {
		read_unmarshal_file(CV_EXPERIENCES_DATA_PATH, &cvData.Experiences)
	})

	wg.Wait()

	close(err_channel)

	for err := range err_channel {
		if err != nil {
			return resume.CVData{}, err
		}
	}

	return *cvData, nil
}

func GetTotalTypes(cvData_experiences *[]resume.Experience, cvData_projects *[]resume.Project) map[string]struct{} {
	total_types := make(map[string]struct{})
	var mu sync.Mutex
	var wg sync.WaitGroup

	utils.Go(&wg, func() {
		var experiencesStringsCombined []string
		for _, exp := range *cvData_experiences {
			experiencesStringsCombined = append(experiencesStringsCombined, exp.CVType...)
		}
		getTypes(total_types, &experiencesStringsCombined, &mu)
	})

	utils.Go(&wg, func() {
		var projectsStringsCombined []string
		for _, proj := range *cvData_projects {
			projectsStringsCombined = append(projectsStringsCombined, proj.CVType...)
		}
		getTypes(total_types, &projectsStringsCombined, &mu)
	})

	wg.Wait()
	return total_types
}

func getTypes(total_types map[string]struct{}, totalCvTypes *[]string, mu *sync.Mutex) {
	mu.Lock()
	for _, cvType := range *totalCvTypes {
		if _, exist := total_types[cvType]; !exist {
			total_types[cvType] = struct{}{}
		}
	}
	mu.Unlock()
}
