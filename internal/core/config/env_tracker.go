package config

import (
	"os"
	"regexp"
	"strings"
)

var envPattern = regexp.MustCompile(`\{\{\s*env\s+"([^"]+)"\s*\}\}`)

type EnvFieldMap map[string]string

func TrackEnvFields(rawContent string) EnvFieldMap {
	matches := envPattern.FindAllStringSubmatch(rawContent, -1)
	fields := make(EnvFieldMap)
	for _, match := range matches {
		if len(match) >= 2 {
			fields[match[0]] = match[1]
		}
	}
	return fields
}

func HasEnvTemplate(value string) bool {
	return envPattern.MatchString(value)
}

func ExtractEnvVarName(value string) string {
	matches := envPattern.FindStringSubmatch(value)
	if len(matches) >= 2 {
		return matches[1]
	}
	return ""
}

func PreserveEnvInYAML(rawContent string, resolvedContent string) string {
	matches := envPattern.FindAllStringSubmatch(rawContent, -1)
	result := resolvedContent
	for _, match := range matches {
		if len(match) >= 2 {
			envValue := os.Getenv(match[1])
			if envValue != "" {
				result = strings.ReplaceAll(result, envValue, match[0])
			}
		}
	}
	return result
}

func RestoreEnvTemplates(resolvedYAML string, rawYAML string) string {
	matches := envPattern.FindAllStringSubmatch(rawYAML, -1)
	result := resolvedYAML
	for _, match := range matches {
		if len(match) >= 2 {
			envValue := os.Getenv(match[1])
			if envValue != "" {
				result = strings.ReplaceAll(result, envValue, match[0])
			}
		}
	}
	return result
}
