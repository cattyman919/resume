package parse

import "core:fmt"
import "core:reflect"
import "core:strings"
import yaml "deps:odin-yaml"

unmarshal :: proc(y: yaml.Value, object: ^$T) {
	unmarshal_recursive(y, object, typeid_of(type_of(object^)))
}

unmarshal_recursive :: proc(y: yaml.Value, object: rawptr, t: typeid) {

	yaml_map, ok := y.(yaml.Mapping)
	if !ok {
		fmt.eprintf(
			"error: expected a yaml.Mapping to unmarshal into a struct, got type %v \n Value: %v \n",
			typeid_of(type_of(y)),
			y,
		)
		return
	}

	struct_tags := reflect.struct_fields_zipped(t)
	object_ptr := uintptr(object)

	for s_field in struct_tags {
		tag_name := get_tag_name(s_field.tag)

		val, exists := yaml_map[tag_name]
		if !exists {
			continue
		}

		field_ptr := object_ptr + s_field.offset

		#partial switch v in val {
		case string:
			(cast(^string)field_ptr)^ = v
		case i64:
			(cast(^i64)field_ptr)^ = v
		case f64:
			(cast(^f64)field_ptr)^ = v
		case bool:
			(cast(^bool)field_ptr)^ = v
		case yaml.Sequence:
			target_array := cast(^[dynamic]string)field_ptr
			clear(target_array)
			for item in v {
				if str_item, is_str := item.(string); is_str {
					append(target_array, str_item)
				}
			}
		case yaml.Mapping:
			if reflect.is_struct(s_field.type) {
				unmarshal_recursive(v, cast(rawptr)field_ptr, s_field.type.id)
			}
		}
	}
}

get_tag_name :: proc(tag: reflect.Struct_Tag) -> string {
	tag_parts := strings.split(string(tag), ":")
	if len(tag_parts) < 2 {
		return "Empty"
	}
	return strings.trim(tag_parts[1], ` "`)
}

