package utils

import "reflect"

func SafeGet(row map[string]any, key string, defaultVal any) any {
	if row == nil {
		return defaultVal
	}
	if val, ok := row[key]; ok && val != nil {
		return val
	}
	return defaultVal
}

func MergeMap(base interface{}, values ...interface{}) (map[string]interface{}, error) {
	out := make(map[string]interface{})

	rv := reflect.ValueOf(base)
	if rv.Kind() == reflect.Struct {
		rt := rv.Type()
		for i := 0; i < rv.NumField(); i++ {
			field := rt.Field(i)
			if field.PkgPath != "" {
				continue // unexported
			}
			out[field.Name] = rv.Field(i).Interface()
		}
	} else if m, ok := base.(map[string]interface{}); ok {
		for k, v := range m {
			out[k] = v
		}
	} else {
		out["Data"] = base
	}

	for i := 0; i < len(values); i += 2 {
		k := values[i].(string)
		out[k] = values[i+1]
	}

	return out, nil
}
