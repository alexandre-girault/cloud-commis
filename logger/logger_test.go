package logger

import (
	"testing"
)

func TestSetLogLevel(t *testing.T) {
	SetLogLevel("info")

	if LogLevel.String() != "LevelVar(INFO)" {
		t.Error(
			"SetLogLevel() doesn't work !" + LogLevel.String())
	}
}
