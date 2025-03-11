package log

import (
	"log/slog"
	"os"
	"strings"
)

var logger *slog.Logger

func SetLogger(level string) {
	opts := &slog.HandlerOptions{
		Level: parseLogLevel(level),
	}
	textHandler := slog.NewTextHandler(os.Stdout, opts) // Use nil for default options
	// jsonHandler := slog.NewJSONHandler(os.Stdout, nil) // Uncomment for JSON logs

	// Create a logger instance with the selected handler
	logger = slog.New(textHandler)
	logger.Debug("Setting up logger", "level", level)
}

// parseLogLevel converts a string to slog.Level
func parseLogLevel(levelStr string) slog.Level {
	switch strings.ToLower(levelStr) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// Info logs an informational message
func Info(msg string, args ...any) {
	logger.Info(msg, args...)
}

// Error logs an error message
func Error(msg string, args ...any) {
	logger.Error(msg, args...)
}

// Debug logs a debug message
func Debug(msg string, args ...any) {
	logger.Debug(msg, args...)
}

// Warn logs a warning message
func Warn(msg string, args ...any) {
	logger.Warn(msg, args...)
}

func Fatal(msg string, args ...any) {
	logger.Error(msg, args...)
	os.Exit(1)
}
