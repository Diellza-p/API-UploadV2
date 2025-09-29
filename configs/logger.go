package configs

import (
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

// Logger is the global logger instance
var Logger *logrus.Logger

// InitLogger initializes the global logger with proper configuration
func InitLogger() {
	Logger = logrus.New()

	// Set log level based on environment
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		level = logrus.InfoLevel
	}
	Logger.SetLevel(level)

	// Set log format to JSON for better parsing
	Logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "message",
		},
	})

	// Create logs directory if it doesn't exist
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		Logger.Warn("Failed to create logs directory", "error", err)
	}

	// Set up file logging
	logFile := filepath.Join(logDir, "api-uploadv2.log")
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		Logger.Warn("Failed to open log file", "error", err, "file", logFile)
	} else {
		Logger.SetOutput(file)
		// Also log to stdout for development
		Logger.AddHook(&WriterHook{
			Writer: os.Stdout,
			Formatter: &logrus.TextFormatter{
				FullTimestamp:   true,
				TimestampFormat: time.RFC3339,
			},
		})
	}

	Logger.Info("Logger initialized",
		"level", level.String(),
		"log_file", logFile,
		"environment", os.Getenv("ENV"))
}

// WriterHook is a custom hook for writing to multiple outputs
type WriterHook struct {
	Writer    *os.File
	Formatter logrus.Formatter
}

func (hook *WriterHook) Fire(entry *logrus.Entry) error {
	formatted, err := hook.Formatter.Format(entry)
	if err != nil {
		return err
	}
	_, err = hook.Writer.Write(formatted)
	return err
}

func (hook *WriterHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// LogWithContext adds context fields to log entries
func LogWithContext(service, operation string) *logrus.Entry {
	return Logger.WithFields(logrus.Fields{
		"service":   service,
		"operation": operation,
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// LogDatabaseOperation logs database operations with timing
func LogDatabaseOperation(operation, collection string, start time.Time, err error) {
	duration := time.Since(start)

	fields := logrus.Fields{
		"operation":  operation,
		"collection": collection,
		"duration":   duration.String(),
		"service":    "database",
	}

	if err != nil {
		Logger.WithFields(fields).Error("Database operation failed", "error", err)
	} else {
		Logger.WithFields(fields).Info("Database operation completed")
	}
}

// LogHTTPRequest logs HTTP requests
func LogHTTPRequest(method, path string, statusCode int, duration time.Duration, clientIP string) {
	fields := logrus.Fields{
		"method":      method,
		"path":        path,
		"status_code": statusCode,
		"duration":    duration.String(),
		"client_ip":   clientIP,
		"service":     "http",
	}

	if statusCode >= 400 {
		Logger.WithFields(fields).Warn("HTTP request completed with error")
	} else {
		Logger.WithFields(fields).Info("HTTP request completed")
	}
}
