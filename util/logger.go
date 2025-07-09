package util

import "fmt"

func Log(format string, args ...any) {
	if AppConfig.Logging {
		fmt.Printf(format, args...)
	}
}
