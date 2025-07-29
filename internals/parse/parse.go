package parse

import (
	model "resume/internals/model"
	"sync"
)

func GetTotalTypes(total_types map[string]struct{}, cvData model.CVData, wg *sync.WaitGroup, mu *sync.Mutex) {
	wg.Add(2)
	go getTypes(total_types, cvData.Experiences, wg, mu)
	go getTypes(total_types, cvData.Projects, wg, mu)
	wg.Wait()
}

func getTypes[T model.Experience | model.Project](total_types map[string]struct{}, items []T, wg *sync.WaitGroup, mu *sync.Mutex) {
	defer wg.Done()
	for _, item := range items {
		var types []string

		switch v := any(item).(type) {
		case model.Experience:
			types = v.CVType
		case model.Project:
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
