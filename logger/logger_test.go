package logger

import (
	"testing"
)

func TestDefaultLogger(t *testing.T) {
	// Create a new logger instance
	logger := NewDefaultLogger()

	// Set the log level (e.g., LevelDebug, LevelInfo, LevelWarn, LevelError)
	logger.SetLevel(LevelDebug)

	// Log messages
	logger.Debug("This is a debug message")
	logger.Infof("This is an info message: %s", "Hello, World!")
	logger.Warn("This is a warning message")
	logger.Errorf("This is an error message: %s", "Something went wrong")

	logger.SetLevel(LevelInfo)
	logger.Debugf("This is debug msg that should not be printed: %s", "Something went wrong")
	logger.WithFields(Fields{"foo": "bar"}).Info("This is an info message with fields")
}
