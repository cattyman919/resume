package model

import "base:runtime"
import "core:fmt"
import "core:reflect"
import "core:strings"
import yaml "deps:odin-yaml"

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

unmarshal_projects :: proc(projects_yaml: yaml.Sequence) -> [dynamic]Project {
	projects := make([dynamic]Project)
	struct_tags := reflect.struct_fields_zipped(typeid_of(Project))

	for p in projects_yaml {
		p := p.(yaml.Mapping)
		new_project := make_project()
		ptr := uintptr(&new_project)

		for s_field in struct_tags {
			tag_name := get_tag_name(s_field.tag)
			value, ok := p[tag_name]

			if !ok {
				continue
			}

			field_ptr := uintptr(ptr + s_field.offset)

			#partial switch v in value {
			case string:
				(cast(^string)field_ptr)^ = v
			case yaml.Sequence:
				target_array: ^[dynamic]string
				target_array = cast(^[dynamic]string)field_ptr
				for item in v {
					if str_item, is_str := item.(string); is_str {
						append(target_array, str_item)
					}
				}

			}
		}

		append(&projects, new_project)
	}

	return projects
}

get_tag_name :: proc(tag: reflect.Struct_Tag) -> string {
	tag_parts := strings.split(string(tag), ":")
	if len(tag_parts) < 2 {
		return "Empty"
	}
	return strings.trim(tag_parts[1], ` "`)
}

