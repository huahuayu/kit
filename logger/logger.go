package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

// ILogger defines the logger interface
type ILogger interface {
	Debug(args ...any)
	Debugf(format string, args ...any)
	Info(args ...any)
	Infof(format string, args ...any)
	Warn(args ...any)
	Warnf(format string, args ...any)
	Error(args ...any)
	Errorf(format string, args ...any)
	SetLevel(level Level)
	WithFields(fields Fields) ILogger
}

// Level defines the log level
type Level int

// Log levels
const (
	LevelDebug Level = -4
	LevelInfo  Level = 0
	LevelWarn  Level = 4
	LevelError Level = 8
)

// Fields defines arbitrary log fields which can be logged in a log entry
type Fields map[string]any

type defaultLogger struct {
	logger *log.Logger
	fields Fields
	level  Level
}

var Logger = NewDefaultLogger()

// NewDefaultLogger creates a new instance of the default logger, by default the log level is LevelInfo and log in json format
func NewDefaultLogger() ILogger {
	return &defaultLogger{
		logger: log.New(os.Stdout, "", 0),
		fields: nil,
		level:  LevelInfo,
	}
}

func (l *defaultLogger) log(level Level, msg string) {
	if l.level <= level {
		logEntry := make(map[string]any)
		logEntry["level"] = levelToString(level)
		logEntry["msg"] = msg
		for k, v := range l.fields {
			if k != "time" {
				logEntry[k] = v
			}
		}
		logEntry["time"] = time.Now().UTC().Format(time.RFC3339)

		// Create a custom JSON string with time as the last key
		jsonStr := "{"
		keys := []string{"level", "msg"}
		for k := range l.fields {
			if k != "time" {
				keys = append(keys, k)
			}
		}
		keys = append(keys, "time")
		for i, key := range keys {
			val := logEntry[key]
			jsonStr += fmt.Sprintf(`"%s": "%v"`, key, val)
			if i < len(keys)-1 {
				jsonStr += ", "
			}
		}
		jsonStr += "}"
		l.logger.Println(jsonStr)
	}
}

func levelToString(level Level) string {
	switch level {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

func (l *defaultLogger) Debug(args ...any) {
	l.log(LevelDebug, fmt.Sprint(args...))
}

func (l *defaultLogger) Debugf(format string, args ...any) {
	l.log(LevelDebug, fmt.Sprintf(format, args...))
}

func (l *defaultLogger) Info(args ...any) {
	l.log(LevelInfo, fmt.Sprint(args...))
}

func (l *defaultLogger) Infof(format string, args ...any) {
	l.log(LevelInfo, fmt.Sprintf(format, args...))
}

func (l *defaultLogger) Warn(args ...any) {
	l.log(LevelWarn, fmt.Sprint(args...))
}

func (l *defaultLogger) Warnf(format string, args ...any) {
	l.log(LevelWarn, fmt.Sprintf(format, args...))
}

func (l *defaultLogger) Error(args ...any) {
	l.log(LevelError, fmt.Sprint(args...))
}

func (l *defaultLogger) Errorf(format string, args ...any) {
	l.log(LevelError, fmt.Sprintf(format, args...))
}

func (l *defaultLogger) SetLevel(level Level) {
	l.level = level
}

func (l *defaultLogger) WithFields(fields Fields) ILogger {
	return &defaultLogger{
		logger: l.logger,
		fields: fields,
		level:  l.level,
	}
}
