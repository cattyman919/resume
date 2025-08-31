package model

import "core:fmt"
import "core:reflect"
import "core:strings"
import yaml "deps:odin-yaml"
import "src:parse"

Project :: struct {
	Name:          string `yaml: "name"`,
	Github:        string `yaml: "github"`,
	Github_Handle: string `yaml: "github_handle"`,
	CV_Types:      [dynamic]string `yaml: "cv_type"`,
	Description:   string `yaml: "description"`,
	Points:        [dynamic]string `yaml: "points"`,
}


make_project :: proc() -> Project {
	return Project{CV_Types = make([dynamic]string), Points = make([dynamic]string)}
}

unmarshal_projects :: proc(projects_yaml: yaml.Value, projects: ^[dynamic]Project) {
	for p_yaml in projects_yaml.(yaml.Sequence) {
		new_project := make_project()
		parse.unmarshal(p_yaml, &new_project)
		append(projects, new_project)
	}
}

