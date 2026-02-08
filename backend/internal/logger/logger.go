package logger

import (
	"log"
	"log/slog"
	"os"
)

// Logger wraps slog.Logger for structured logging
type Logger struct {
	*slog.Logger
}

// New creates a new logger instance
// In production, you'd configure JSON output and log levels via env vars
func New() *Logger {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
		// AddCaller shows file:line in logs (useful for debugging)
		AddCaller: true,
	}

	var handler slog.Handler
	// Use JSON handler in production, text handler for development
	if os.Getenv("LOG_FORMAT") == "json" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return &Logger{
		Logger: slog.New(handler),
	}
}

// Default logger instance
var defaultLogger *Logger

func init() {
	defaultLogger = New()
}

// Default returns the default logger instance
func Default() *Logger {
	return defaultLogger
}

// Logging middleware helper - wraps standard log for compatibility
func LogError(err error, msg string) {
	if err != nil {
		defaultLogger.Error(msg, "error", err)
	} else {
		defaultLogger.Info(msg)
	}
}

// Fallback for code that still uses standard log
func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
