package helpers

import (
	"fmt"
	"time"
)

// build file name
func GenFileName(base string) string {
	t := time.Now()

	return fmt.Sprintf("%s-%d_%02d_%02dT%02d_%02d_%02d", base, t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}
