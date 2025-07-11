package logger

// BasicLogger defines the minimal logging interface following ISP
type BasicLogger interface {
	Debug(message string, fields ...map[string]interface{})
	Info(message string, fields ...map[string]interface{})
	Warn(message string, fields ...map[string]interface{})
	Error(message string, fields ...map[string]interface{})
	Fatal(message string, fields ...map[string]interface{})
}

// FieldLogger adds field management capabilities
type FieldLogger interface {
	BasicLogger
	WithFields(fields map[string]interface{}) FieldLogger
}

// ConfigurableLogger adds configuration capabilities
type ConfigurableLogger interface {
	SetLevel(level Level)
	SetFormatter(formatter Formatter)
	AddHandler(handler Handler)
}

// FullLogger combines all logging capabilities
type FullLogger interface {
	BasicLogger
	FieldLogger
	ConfigurableLogger
}

// LevelProvider provides access to log levels
type LevelProvider interface {
	GetLevel() Level
}

// OutputProvider provides access to output destinations
type OutputProvider interface {
	GetHandlers() []Handler
}

// Formatter defines how to format log messages
type Formatter interface {
	Format(level Level, message string, fields map[string]interface{}) string
}

// Handler defines how to process log messages
type Handler interface {
	Handle(level Level, message string, fields map[string]interface{}) error
}

type Logger interface {
	Debug(message string, fields ...map[string]interface{})
	Info(message string, fields ...map[string]interface{})
	Warn(message string, fields ...map[string]interface{})
	Error(message string, fields ...map[string]interface{})
	Fatal(message string, fields ...map[string]interface{})
	WithFields(fields map[string]interface{}) Logger
	SetLevel(level Level)
	SetFormatter(formatter Formatter)
	AddHandler(handler Handler)
}
