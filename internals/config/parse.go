package config

import (
	"os"
	"resume/internals/resume"
	"sync"

	"gopkg.in/yaml.v3"
)

const (
	CV_GENERAL_DATA_PATH     = "config/general.yaml"
	CV_PROJECTS_DATA_PATH    = "config/projects.yaml"
	CV_EXPERIENCES_DATA_PATH = "config/experiences.yaml"
)

func LoadYAMLData(wg *sync.WaitGroup) (resume.CVData, error) {
	var cvData = new(resume.CVData)
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
			return resume.CVData{}, err
		}
	}

	return *cvData, nil
}

func GetTotalTypes(total_types map[string]struct{}, cvData_experiences *[]resume.Experience, cvData_projects *[]resume.Project, wg *sync.WaitGroup, mu *sync.Mutex) {
	wg.Add(2)
	go getTypes(total_types, cvData_experiences, wg, mu)
	go getTypes(total_types, cvData_projects, wg, mu)
	wg.Wait()
}

func getTypes[T resume.Experience | resume.Project](total_types map[string]struct{}, items *[]T, wg *sync.WaitGroup, mu *sync.Mutex) {
	defer wg.Done()
	for _, item := range *items {
		var types []string

		switch v := any(item).(type) {
		case resume.Experience:
			types = v.CVType
		case resume.Project:
			types = v.CVType
		}

		if len(types) == 0 {
			continue
		}

		mu.Lock()
		for _, cvType := range types {
			if _, exist := total_types[cvType]; !exist {
				total_types[cvType] = struct{}{}
			}
		}
		mu.Unlock()
	}

}
