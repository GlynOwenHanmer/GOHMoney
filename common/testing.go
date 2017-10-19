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
	if err == nil {
		return
	}
	t.Fatalf("%s: %s", fmt.Sprintf(format, args...), err)
}

func ErrorIfError(t *testing.T, err error, message string) {
	if err == nil {
		return
	}
	t.Errorf("%s: %s", message, err)
}

func ErrorIfErrorf(t *testing.T, err error, format string, args ...interface{}) {
	if err == nil {
		return
	}
	t.Errorf("%s: %s", fmt.Sprintf(format, args...), err)
}
