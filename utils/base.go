package utils

import (
	"os"
	"reflect"
	"strings"
)

func SafeEnvGet(name string, defaultVal any) any {
	val := strings.TrimSpace(os.Getenv(name))
	if len(val) > 0 {
		return val
	}
	return defaultVal
}

func IsTruthy(v interface{}) bool {
	if v == nil {
		return false
	}

	switch val := v.(type) {
	case string:
		return val != ""
	case bool:
		return val
	case int:
		return val != 0
	case int8, int16, int32, int64:
		return reflect.ValueOf(v).Int() != 0
	case uint, uint8, uint16, uint32, uint64:
		return reflect.ValueOf(v).Uint() != 0
	case float32, float64:
		return reflect.ValueOf(v).Float() != 0
	case []interface{}:
		return len(val) > 0
	case []string:
		return len(val) > 0
	default:
		rv := reflect.ValueOf(v)
		switch rv.Kind() {
		case reflect.Slice, reflect.Array, reflect.Map:
			return rv.Len() > 0
		case reflect.Ptr, reflect.Interface:
			return !rv.IsNil()
		}
	}

	return true
}
