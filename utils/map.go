package utils

import (
	"reflect"
	"strings"
)

func ModelToMap(model any, parsableTag string) map[string]any {
	if model == nil {
		return map[string]any{}
	}
	// Already in desired format.
	if m, ok := model.(map[string]any); ok {
		return m
	}
	result := make(map[string]any)

	t := reflect.TypeOf(model)
	v := reflect.ValueOf(model)

	if t == nil {
		return result
	}

	if t.Kind() == reflect.Ptr {
		if v.IsNil() {
			return result
		}

		t = t.Elem()
		v = v.Elem()
	}

	if t.Kind() != reflect.Struct {
		return result
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		tag := field.Tag.Get(parsableTag)

		if tag == "" || tag == "-" {
			continue
		}

		column := strings.Split(tag, ",")[0]

		if column == "" || column == "-" {
			continue
		}

		result[column] = v.Field(i).Interface()
	}

	return result
}
