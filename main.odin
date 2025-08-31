package main

import "core:fmt"
import "core:os"
import yaml "deps:odin-yaml"
import "src/model"

main :: proc() {

	doc, err := yaml.decode_from_file("config/projects.yaml")
	if err != nil {
		fmt.eprintfln("---- {}", yaml.error_string(err, "config/projects.yaml"))
		os.exit(1)
	}
	projects_sequence := doc.(yaml.Sequence)


	// fmt.println(projects_slice[0])

	projects := model.unmarshal_projects(projects_sequence)

	// project.Name = projects_slice[0].(yaml.Mapping)["name"].(string)
}

