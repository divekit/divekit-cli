package logUtils

import (
	"github.com/apex/log"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLogLevelAsString(t *testing.T) {
	testCases := []struct {
		name     string
		logLevel log.Level // input
		string   string    // expected
	}{
		{"DebugLevel should be debug", log.DebugLevel, "debug"},
		{"InfoLevel should be info", log.InfoLevel, "info"},
		{"WarningLevel should be warning", log.WarnLevel, "warning"},
		{"ErrorLevel should be containsError", log.ErrorLevel, "error"},
		{"FatalLevel should be info", log.FatalLevel, "info"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			LogLevel = testCase.logLevel
			str := LogLevelAsString()
			assert.Equal(t, testCase.string, str, "input"+testCase.logLevel.String())
		})
	}
}

func TestStringAsLogLevel(t *testing.T) {
	testCases := []struct {
		name     string
		string   string    // input
		logLevel log.Level // expected
		error    error     // expected
	}{
		{"debug should be DebugLevel", "debug", log.DebugLevel, nil},
		{"info should be InfoLevel", "info", log.InfoLevel, nil},
		{"warning should be WarnLevel", "warning", log.WarnLevel, nil},
		{"containsError should be ErrorLevel", "error", log.ErrorLevel, nil},
		{"invalid should be InfoLevel and contain an InvalidLogLevelError", "invalid",
			log.InfoLevel, log.ErrInvalidLevel},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			logLevel, err := StringAsLogLevel(testCase.string)
			assert.Equal(t, testCase.logLevel, logLevel, "input"+testCase.string)
			assert.IsType(t, testCase.error, err, "invalid log level string: %v", testCase.string)
		})
	}
}
