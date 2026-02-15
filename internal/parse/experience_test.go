package parse

import (
	"path/filepath"
	"testing"
)

func TestParseExperience(t *testing.T) {
	ep, err := NewExperienceParser(filepath.Join(
		"/home",
		"senohebat",
		"repo",
		"resume",
		"src",
		"config",
		"experiences.typ",
	))
	if err != nil {
		t.Error(err)
	}
	defer ep.Close()
	ep.Parse()
}
