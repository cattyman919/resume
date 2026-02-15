package parse

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/cattyman919/autocv/internal/model"
)

type ExperienceParser struct {
	file    *os.File
	scanner *bufio.Scanner
}

func NewExperienceParser(filePath string) (*ExperienceParser, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)

	return &ExperienceParser{
		file:    file,
		scanner: scanner,
	}, nil
}

func (ep *ExperienceParser) Close() {
	ep.file.Close()
}

func (ep *ExperienceParser) Parse() []model.Experience {

	experiences := make([]model.Experience, 0, 10)
	prefix := "#let "
	suffix := " = make-experience("

	for ep.scanner.Scan() {
		line := strings.TrimSpace(ep.scanner.Text())
		if strings.HasPrefix(line, prefix) && strings.HasSuffix(line, suffix) {
			experience := model.Experience{}
			st := NewStack[rune]()

			// lineExperience := strings.TrimPrefix(line, "#let ")
			// company := strings.TrimSuffix(lineExperience, suffix)
			// experience.Company = company

			st.Push('(')
			for len(st.elements) != 0 && ep.scanner.Scan() {
				line := strings.TrimSpace(ep.scanner.Text())

				for _, ch := range line {
					switch ch {
					case '(':
						st.Push('(')
					case ')':
						st.Pop()
					}
				}

				switch {
				case strings.HasPrefix(line, "company:"):
					company := trimField(line, "company:")
					experience.Company = company
				case strings.HasPrefix(line, "location:"):
					location := trimField(line, "location:")
					experience.Location = location
				case strings.HasPrefix(line, "role:"):
					role := trimField(line, "role:")
					experience.Role = role

				case strings.HasPrefix(line, "job-type:"):
					jobType := trimField(line, "job-type:")
					experience.JobType = jobType

				case strings.HasPrefix(line, "dates:"):
					fromDate, afterDate := getDatesString(line)
					experience.FromDate = fromDate
					experience.AfterDate = afterDate

				case strings.HasPrefix(line, "points:"):
					pointLines := make([]string, 0, 20)

					for ep.scanner.Scan() {
						line := strings.TrimSpace(ep.scanner.Text())
						line = strings.TrimSuffix(line, ",")
						if strings.HasPrefix(line, ")") {
							st.Pop()
							break
						}
						if strings.HasPrefix(line, "[") {
							pointLines = append(pointLines, line)
						}
					}

					parsedPointLines := parseExperiencePoints(pointLines)
					experience.Points = parsedPointLines
				}
			}
			experiences = append(experiences, experience)
		}
	}

	for _, exp := range experiences {
		fmt.Printf("%v\n", exp)
	}

	return experiences
}

func getDatesString(line string) (string, string) {
	line = strings.TrimPrefix(line, "dates: dates-helper")
	line = strings.TrimPrefix(line, "(")
	line = strings.TrimSuffix(line, "),")

	secondPart := false

	fromDate := ""
	afterDate := ""

	runes := []rune(line)

	for i := 0; i < len(runes); i++ {
		ch := runes[i]

		if ch == '"' {
			j := i + 1
			for ; j < len(runes); j++ {
				innerCh := runes[j]
				if innerCh == '"' {
					if !secondPart {
						fromDate = line[i+1 : j]
						secondPart = true
					} else {
						afterDate = line[i+1 : j]
					}
					break
				}
			}
			i = j
		}
	}

	return fromDate, afterDate
}

func trimField(input string, prefix string) string {
	res := strings.TrimPrefix(input, prefix)
	res = strings.TrimSpace(res)
	res = strings.TrimSuffix(res, ",")
	res = strings.TrimPrefix(res, "\"")
	res = strings.TrimSuffix(res, "\"")
	return res
}

func parseExperiencePoints(pointLines []string) []string {
	res := make([]string, len(pointLines))

	for i, line := range pointLines {
		line = strings.TrimSpace(line)
		line = strings.TrimSuffix(line, ",")
		line = strings.TrimPrefix(line, "[")
		line = strings.TrimSuffix(line, "]")
		res[i] = line
	}

	return res
}
