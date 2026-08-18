package utils

import (
	"encoding/json"
	"fmt"

	"github.com/jinzhu/inflection"
)

func Coalesce(vals ...interface{}) interface{} {
	for _, v := range vals {
		if IsTruthy(v) {
			return v
		}
	}
	return nil
}

func ToJSON(data any) string {
	bytes, err := json.Marshal(data)
	if err != nil {
		fmt.Println("unable to transform to JSON", err.Error())
		return ""
	}
	return string(bytes)
}

func Singular(plural string) string {
	return inflection.Singular(plural)
}

func ToString(v any) string {
	return fmt.Sprint(v)
}
