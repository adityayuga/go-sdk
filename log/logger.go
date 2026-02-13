package log

import (
	"context"
	"fmt"

	"github.com/adityayuga/go-sdk/debug"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps zap.Logger with context-aware logging
type Logger struct {
	*zap.Logger
}

var global *Logger

// Init initializes the global logger with JSON encoding
func Init(level string) error {
	cfg := zap.NewProductionConfig()
	cfg.Encoding = "json"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Set log level
	switch level {
	case "debug":
		cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		cfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		cfg.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		cfg.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	default:
		cfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}

	logger, err := cfg.Build()
	if err != nil {
		return err
	}

	global = &Logger{Logger: logger}
	return nil
}

// InitDevelopment initializes the logger in development mode with pretty printing
func InitDevelopment() error {
	cfg := zap.NewDevelopmentConfig()
	cfg.Encoding = "json"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := cfg.Build()
	if err != nil {
		return err
	}

	global = &Logger{Logger: logger}
	return nil
}

// NewLogger creates a new logger instance
func NewLogger(level string) (*Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.Encoding = "json"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	switch level {
	case "debug":
		cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		cfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		cfg.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		cfg.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	default:
		cfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}

	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return &Logger{Logger: logger}, nil
}

// GetLogger returns the global logger instance
func GetLogger() *Logger {
	if global == nil {
		Init("info")
	}
	return global
}

// Info logs an info-level message with context
func (l *Logger) Info(ctx context.Context, message string, fields ...interface{}) {
	fields = extractContextFields(ctx, fields)
	l.Logger.Info(message, toZapFields(fields...)...)
}

// Warn logs a warning-level message with context
func (l *Logger) Warn(ctx context.Context, message string, fields ...interface{}) {
	fields = extractContextFields(ctx, fields)
	l.Logger.Warn(message, toZapFields(fields...)...)
}

// Error logs an error-level message with context
func (l *Logger) Error(ctx context.Context, message string, fields ...interface{}) {
	fields = extractContextFields(ctx, fields)
	l.Logger.Error(message, toZapFields(fields...)...)
}

// Debug logs a debug-level message with context
func (l *Logger) Debug(ctx context.Context, message string, fields ...interface{}) {
	fields = extractContextFields(ctx, fields)
	l.Logger.Debug(message, toZapFields(fields...)...)
}

// Fatal logs a fatal-level message and exits
func (l *Logger) Fatal(ctx context.Context, message string, fields ...interface{}) {
	fields = extractContextFields(ctx, fields)
	l.Logger.Fatal(message, toZapFields(fields...)...)
}

// Sync flushes any buffered log entries
func (l *Logger) Sync() error {
	return l.Logger.Sync()
}

// Helper functions

// extractContextFields extracts request ID, trace ID, user ID, and debug ID from context if present
func extractContextFields(ctx context.Context, fields []interface{}) []interface{} {
	if ctx == nil {
		return fields
	}

	if requestID := ctx.Value("request_id"); requestID != nil {
		fields = append(fields, "request_id", fmt.Sprintf("%v", requestID))
	}

	if traceID := ctx.Value("trace_id"); traceID != nil {
		fields = append(fields, "trace_id", fmt.Sprintf("%v", traceID))
	}

	if userID := ctx.Value("user_id"); userID != nil {
		fields = append(fields, "user_id", fmt.Sprintf("%v", userID))
	}

	// Get debug ID from debug package
	debugID, _ := debug.GetDebugIDFromContext(ctx)
	if debugID != "" {
		fields = append(fields, "debug_id", debugID)
	}

	return fields
}

// toZapFields converts interface{} pairs to zap.Field array
func toZapFields(keyvals ...interface{}) []zap.Field {
	fields := make([]zap.Field, 0)

	for i := 0; i < len(keyvals); i += 2 {
		if i+1 >= len(keyvals) {
			break
		}

		key := fmt.Sprintf("%v", keyvals[i])
		value := keyvals[i+1]

		fields = append(fields, zap.Any(key, value))
	}

	return fields
}

// Global function wrappers for convenience

// Info logs an info message using the global logger
func Info(ctx context.Context, message string, fields ...interface{}) {
	GetLogger().Info(ctx, message, fields...)
}

// Warn logs a warning message using the global logger
func Warn(ctx context.Context, message string, fields ...interface{}) {
	GetLogger().Warn(ctx, message, fields...)
}

// Error logs an error message using the global logger
func Error(ctx context.Context, message string, fields ...interface{}) {
	GetLogger().Error(ctx, message, fields...)
}

// Debug logs a debug message using the global logger
func Debug(ctx context.Context, message string, fields ...interface{}) {
	GetLogger().Debug(ctx, message, fields...)
}

// Fatal logs a fatal message using the global logger and exits
func Fatal(ctx context.Context, message string, fields ...interface{}) {
	GetLogger().Fatal(ctx, message, fields...)
}

// Infof logs a formatted info-level message with context
func (l *Logger) Infof(ctx context.Context, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	l.Info(ctx, message)
}

// Warnf logs a formatted warning-level message with context
func (l *Logger) Warnf(ctx context.Context, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	l.Warn(ctx, message)
}

// Errorf logs a formatted error-level message with context
func (l *Logger) Errorf(ctx context.Context, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	l.Error(ctx, message)
}

// Debugf logs a formatted debug-level message with context
func (l *Logger) Debugf(ctx context.Context, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	l.Debug(ctx, message)
}

// Fatalf logs a formatted fatal-level message and exits
func (l *Logger) Fatalf(ctx context.Context, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	l.Fatal(ctx, message)
}

// Global formatted function wrappers for convenience

// Infof logs a formatted info message using the global logger
func Infof(ctx context.Context, format string, args ...interface{}) {
	GetLogger().Infof(ctx, format, args...)
}

// Warnf logs a formatted warning message using the global logger
func Warnf(ctx context.Context, format string, args ...interface{}) {
	GetLogger().Warnf(ctx, format, args...)
}

// Errorf logs a formatted error message using the global logger
func Errorf(ctx context.Context, format string, args ...interface{}) {
	GetLogger().Errorf(ctx, format, args...)
}

// Debugf logs a formatted debug message using the global logger
func Debugf(ctx context.Context, format string, args ...interface{}) {
	GetLogger().Debugf(ctx, format, args...)
}

// Fatalf logs a formatted fatal message using the global logger and exits
func Fatalf(ctx context.Context, format string, args ...interface{}) {
	GetLogger().Fatalf(ctx, format, args...)
}
