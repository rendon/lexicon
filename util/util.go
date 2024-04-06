package util

import (
	"fmt"
	"time"
)

// Min returns the minimum of two integers.
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// FormatTime returns Time as a string with this format yyyy-mm-dd hh:mm:ss.
func FormatTime(t time.Time) string {
	return fmt.Sprintf(
		"%04d-%02d-%02d %02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
}
