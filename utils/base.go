package utils

import (
	"os"
	"strings"
)

type BaseUtil struct{}

func NewBaseUtil() *BaseUtil {
	return &BaseUtil{}
}

func (util *BaseUtil) SafeEnvGet(name string, defaultVal string) string {
	val := strings.TrimSpace(os.Getenv(name))
	if len(val) > 0 {
		return val
	}
	return defaultVal
}
