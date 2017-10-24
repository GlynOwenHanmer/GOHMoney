package common

import (
	"fmt"
	"testing"
)

func FatalIfError(t *testing.T, err error, message string) {
	if err == nil {
		return
	}
	t.Fatalf("%s: %s", message, err)
}

func FatalIfErrorf(t *testing.T, err error, format string, args ...interface{}) {
	FatalIfError(t, err, fmt.Sprintf(format, args...))
}

func ErrorIfError(t *testing.T, err error, message string) {
	if err == nil {
		return
	}
	t.Errorf("%s: %s", message, err)
}

func ErrorIfErrorf(t *testing.T, err error, format string, args ...interface{}) {
	ErrorIfError(t, err, fmt.Sprintf(format, args...))
}
