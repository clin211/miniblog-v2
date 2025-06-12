package store

import (
	"context"
)

// Logger defines an interface for logging errors with contextual information.
type Logger interface {
	// Error logs an error message with the associated context.
	Error(ctx context.Context, err error, message string, kvs ...any)
}

// emptyLogger is a no-op logger that implements the Logger interface.
// It does not perform any logging operations.
type emptyLogger struct{}

// NewLogger creates and returns a new instance of emptyLogger.
func NewLogger() *emptyLogger {
	return &emptyLogger{} // Return a new instance of emptyLogger
}

// Error is a no-op method that satisfies the Logger interface.
// It does not log any error messages or context.
func (l *emptyLogger) Error(ctx context.Context, err error, msg string, kvs ...any) {
	// No operation performed for logging errors
}
