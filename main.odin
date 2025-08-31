package main

import "core:fmt"
import "core:os"
import yaml "deps:odin-yaml"
import "src/model"
import "src/parse"

main :: proc() {

	doc, err := yaml.decode_from_file("config/projects.yaml")
	if err != nil {
		fmt.eprintfln("---- {}", yaml.error_string(err, "config/projects.yaml"))
		os.exit(1)
	}
	project_seq := doc.(yaml.Sequence)

	projects := make([dynamic]model.Project)
	model.unmarshal_projects(doc, &projects)
	fmt.println(projects)

	// projects := model.unmarshal_projects(doc.(yaml.Sequence))
	// project := model.make_project()
	// parse.unmarshal(project_seq[0], &project)
	// fmt.println(project)
}

