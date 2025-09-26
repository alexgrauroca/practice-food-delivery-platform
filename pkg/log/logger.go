// Package log provides a standardized logging interface. It implements the logging conventions defined in ADR-0004,
// supporting structured, leveled logging with context propagation and consistent field naming.
package log

import (
	"context"
	"log"

	"go.uber.org/zap"
)

// Logger provides methods for structured and leveled logging within the application.
// Sync flushes any buffered log entries.
// WithContext returns a logger instance enriched with context-specific details.
// Debug logs a debug-level message, with optional fields.
// Info logs an information-level message, with optional fields.
// Warn logs a warning-level message, with optional fields.
// Error logs an error-level message and the related error details.
// Fatal logs a fatal-level message, followed by the termination of the application.
type Logger interface {
	Sync() error
	WithContext(ctx context.Context) Logger
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, err error)
	Fatal(msg string, err error)
}

// Field represents a structured key-value pair attached to a log entry.
// Key is the field name used in logs; Value holds the associated data and
// can be any type supported by the underlying logger.
type Field struct {
	Key   string
	Value any
}

// NewProduction constructs a production-ready Logger backed by zap.
// It enables caller annotation and skips this wrapper's frame so that
// the reported caller points to the originating code.
func NewProduction() (Logger, error) {
	logger, err := zap.NewProduction(
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
		return nil, err
	}

	return &zapLogger{logger: logger}, nil
}

// NewTest returns a no-op Logger intended for tests where log output is
// suppressed; its methods are safe to call without producing log entries.
func NewTest() (Logger, error) {
	return &zapLogger{logger: zap.NewNop()}, nil
}

type zapLogger struct {
	logger *zap.Logger
}

func (l *zapLogger) Sync() error {
	err := l.logger.Sync()
	if err != nil && err.Error() != "sync /dev/stderr: invalid argument" {
		return err
	}
	return nil
}

func (l *zapLogger) WithContext(ctx context.Context) Logger {
	return &zapLogger{
		logger: LoggerWithRequestInfo(ctx, l.logger),
	}
}

func (l *zapLogger) Debug(msg string, fields ...Field) {
	l.logger.Debug(msg, toZapFields(fields)...)
}

func (l *zapLogger) Info(msg string, fields ...Field) {
	l.logger.Info(msg, toZapFields(fields)...)
}

func (l *zapLogger) Warn(msg string, fields ...Field) {
	l.logger.Warn(msg, toZapFields(fields)...)
}

func (l *zapLogger) Error(msg string, err error) {
	l.logger.Error(msg, zap.Error(err))
}

func (l *zapLogger) Fatal(msg string, err error) {
	l.logger.Fatal(msg, zap.Error(err))
}

// toZapFields converts Field slice to zap.Field slice
func toZapFields(fields []Field) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for i, field := range fields {
		zapFields[i] = zap.Any(field.Key, field.Value)
	}
	return zapFields
}
