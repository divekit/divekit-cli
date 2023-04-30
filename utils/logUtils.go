package utils

import (
	"fmt"
	"github.com/apex/log"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"io"
	"os"
	"sync"
)

type Chattyness int

const (
	Silent Chattyness = iota
	Verbose
	VeryVerbose
)

var (
	LogLevel = log.InfoLevel
)

type CustomHandler struct {
	mu sync.Mutex
	w  io.Writer
}

func (h *CustomHandler) HandleLog(e *log.Entry) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Format the log message
	msg := fmt.Sprintf("[%s] %s\n", cases.Title(language.Dutch).String(fmt.Sprintf("%s", e.Level)), e.Message)

	// Write the formatted message to the output writer
	_, err := h.w.Write([]byte(msg))
	return err
}

func NewCustomHandler(w io.Writer) *CustomHandler {
	return &CustomHandler{
		w: w,
	}
}

func DefineLoggingConfig(verboseFlag bool, debugFlag bool) {
	// Create and set the custom handler
	customHandler := NewCustomHandler(os.Stdout)
	log.SetHandler(customHandler)
	if verboseFlag {
		log.SetLevel(log.InfoLevel)
		LogLevel = log.InfoLevel
	} else {
		log.SetLevel(log.WarnLevel)
		LogLevel = log.WarnLevel
	}
	if debugFlag {
		log.SetLevel(log.DebugLevel)
		LogLevel = log.DebugLevel
	}
}

func LogLevelAsString() string {
	switch LogLevel {
	case log.DebugLevel:
		return "debug"
	case log.InfoLevel:
		return "info"
	case log.WarnLevel:
		return "warning"
	case log.ErrorLevel:
		return "error"
	default:
		return "info"
	}
}
