package logUtils

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
	msg := fmt.Sprintf("[%s] %s\n", cases.Title(language.English).String(fmt.Sprintf("%s", e.Level)), e.Message)

	// Write the formatted message to the output writer
	_, err := h.w.Write([]byte(msg))
	return err
}

func NewCustomHandler(w io.Writer) *CustomHandler {
	return &CustomHandler{
		w: w,
	}
}

func DefineLoggingLevel(logLevelString string) error {
	// Create and set the custom handler
	customHandler := NewCustomHandler(os.Stdout)
	log.SetHandler(customHandler)
	var err error = nil
	LogLevel, err = StringAsLogLevel(logLevelString)
	log.Info("Log level set to " + LogLevelAsString() + ".")
	return err
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

func StringAsLogLevel(levelStr string) (log.Level, error) {
	switch levelStr {
	case "debug":
		return log.DebugLevel, nil
	case "info":
		return log.InfoLevel, nil
	case "warning":
		return log.WarnLevel, nil
	case "error":
		return log.ErrorLevel, nil
	default:
		return log.InfoLevel, log.ErrInvalidLevel
	}
}
