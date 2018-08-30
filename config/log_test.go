package config

import (
	"testing"
)

func TestLogDefault(t *testing.T) {
	assertViperDefault(t, map[string]string{
		LogFormat: "text",
		LogLevel:  "info",
	})
}
