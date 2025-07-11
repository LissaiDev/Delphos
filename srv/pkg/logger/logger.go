package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

// Instatiating a singleton
var (
	LogInstance Logger
	once        sync.Once
)

// Level represents the log level
type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	FATAL
)

// String returns the string representation of the log level
func (l Level) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Color returns the ANSI color for the log level
func (l Level) Color() string {
	switch l {
	case DEBUG:
		return "\033[36m" // Cyan
	case INFO:
		return "\033[32m" // Green
	case WARN:
		return "\033[33m" // Yellow
	case ERROR:
		return "\033[31m" // Red
	case FATAL:
		return "\033[35m" // Magenta
	default:
		return "\033[0m" // Reset
	}
}

// Logger is the main logger interface

// defaultFormatter is the default implementation of Formatter
type defaultFormatter struct {
	useColors bool
	showTime  bool
}

// Format formats the log message
func (f *defaultFormatter) Format(level Level, message string, fields map[string]interface{}) string {
	var result string

	if f.showTime {
		result += time.Now().Format("2006-01-02 15:04:05") + " "
	}

	if f.useColors {
		result += level.Color() + level.String() + "\033[0m "
	} else {
		result += level.String() + " "
	}

	result += message

	if len(fields) > 0 {
		result += " | "
		for k, v := range fields {
			result += fmt.Sprintf("%s=%v ", k, v)
		}
	}

	return result + "\n---"
}

// defaultHandler is the default implementation of Handler
type defaultHandler struct {
	writer io.Writer
	mutex  sync.Mutex
}

// Handle processes the log message
func (h *defaultHandler) Handle(level Level, message string, fields map[string]interface{}) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	_, err := fmt.Fprintln(h.writer, message)
	return err
}

// logger is the main implementation of Logger
type logger struct {
	level     Level
	formatter Formatter
	handlers  []Handler
	fields    map[string]interface{}
	mutex     sync.RWMutex
}

// New creates a new logger
func New() Logger {
	return &logger{
		level:     INFO,
		formatter: &defaultFormatter{useColors: true, showTime: true},
		handlers:  []Handler{&defaultHandler{writer: os.Stdout}},
		fields:    make(map[string]interface{}),
	}
}

// NewWithLevel creates a new logger with a specific level
func NewWithLevel(level Level) Logger {
	l := New()
	l.SetLevel(level)
	return l
}

// log records a message at the specified level
func (l *logger) log(level Level, message string, fields map[string]interface{}) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	if level < l.level {
		return
	}

	// Mescla campos globais com campos específicos
	allFields := make(map[string]interface{})
	for k, v := range l.fields {
		allFields[k] = v
	}
	for k, v := range fields {
		allFields[k] = v
	}

	formatted := l.formatter.Format(level, message, allFields)

	for _, handler := range l.handlers {
		if err := handler.Handle(level, formatted, allFields); err != nil {
			// Fallback para log padrão em caso de erro
			log.Printf("Logger error: %v", err)
		}
	}

	// Para FATAL, sempre sair
	if level == FATAL {
		os.Exit(1)
	}
}

func (l *logger) Debug(message string, fields ...map[string]interface{}) {
	l.log(DEBUG, message, l.getFields(fields...))
}

func (l *logger) Info(message string, fields ...map[string]interface{}) {
	l.log(INFO, message, l.getFields(fields...))
}

func (l *logger) Warn(message string, fields ...map[string]interface{}) {
	l.log(WARN, message, l.getFields(fields...))
}

func (l *logger) Error(message string, fields ...map[string]interface{}) {
	l.log(ERROR, message, l.getFields(fields...))
}

func (l *logger) Fatal(message string, fields ...map[string]interface{}) {
	l.log(FATAL, message, l.getFields(fields...))
}

func (l *logger) WithFields(fields map[string]interface{}) Logger {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	newLogger := &logger{
		level:     l.level,
		formatter: l.formatter,
		handlers:  l.handlers,
		fields:    make(map[string]interface{}),
	}

	// Copia campos existentes
	for k, v := range l.fields {
		newLogger.fields[k] = v
	}

	// Adiciona novos campos
	for k, v := range fields {
		newLogger.fields[k] = v
	}

	return newLogger
}

func (l *logger) SetLevel(level Level) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.level = level
}

func (l *logger) SetFormatter(formatter Formatter) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.formatter = formatter
}

func (l *logger) AddHandler(handler Handler) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.handlers = append(l.handlers, handler)
}

// getFields extracts fields from a slice of maps
func (l *logger) getFields(fields ...map[string]interface{}) map[string]interface{} {
	if len(fields) == 0 {
		return make(map[string]interface{})
	}
	return fields[0]
}

// GetLevel returns the current log level
func (l *logger) GetLevel() Level {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return l.level
}

// GetHandlers returns the current handlers
func (l *logger) GetHandlers() []Handler {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return l.handlers
}

func GetInstance() Logger {
	once.Do(func() {
		LogInstance = New()
	})

	return LogInstance
}
