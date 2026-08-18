package utils

import (
	"log"
)

var logger = log.Default()

func Log(a ...any) {
	logger.Println(a...)
}
func LogF(format string, args ...any) {
	logger.Printf(format, args...)
}
