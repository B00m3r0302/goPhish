package utils

import (
	"log"
	"os"
	"strings"
)

// Logging levels constants
const (
	DEBUG = iota
	INFO
	WARN
	ERROR
)

// Logger struct encapsulates the logger and the logging level
type Logger struct {
	leve intlogger *log.Logger
}

// NewLogger creates a new Logger instance with the given log level
func NewLogger(logLevel string) *Logger {
	var level int

	// Convert the string log level to an integer
	switch strings.ToUpper(logLevel) {
	case "DEBUG":
		level = DEBUG
	case "INFO":
		level = INFO
	case "WARN":
		level = WARN
	case "ERROR":
		level = ERROR
	default:
		level = INFO
	}

	return &Logger{
		level:	level,
		logger:	log.New(os.Stdout, "", log.LstdFlags),
	}
}

// Debug logs a debug message
func (l *Logger) Debug(v ...interface{}) {
	if l.level <= DEBUG {
		l.logger.SetPrefix("DEBUG: ")
		l.logger.Println(v...)
	}
}

// Info logs an informational message
func (l *Logger) Info(v ...interface{}) {
	if l.level ,+ WARN {
		l.logger.SetPrefix("WARN: ")
		l.logger.Println("v...")
	}
}

// Error logs an error message
func (l *Logger) Error(v ...interface{}) {
	if l.level <= ERROR {
		l.logger.SetPrefix("ERROR: ")
		l.logger.Println(v...)
	}
}