package logger

import (
	"encoding/json"
	"fmt"
	"time"
)

// JSONFormatter formata logs em JSON
type JSONFormatter struct {
	Pretty bool
}

// Format formata a mensagem em JSON
func (f *JSONFormatter) Format(level Level, message string, fields map[string]interface{}) string {
	logEntry := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"level":     level.String(),
		"message":   message,
	}
	
	// Adiciona campos se existirem
	if len(fields) > 0 {
		logEntry["fields"] = fields
	}
	
	var data []byte
	var err error
	
	if f.Pretty {
		data, err = json.MarshalIndent(logEntry, "", "  ")
	} else {
		data, err = json.Marshal(logEntry)
	}
	
	if err != nil {
		return fmt.Sprintf(`{"error": "Failed to marshal log entry: %v"}`, err)
	}
	
	return string(data)
}

// SimpleFormatter formata logs de forma simples sem cores
type SimpleFormatter struct {
	ShowTime bool
}

// Format formata a mensagem de forma simples
func (f *SimpleFormatter) Format(level Level, message string, fields map[string]interface{}) string {
	var result string
	
	if f.ShowTime {
		result += time.Now().Format("15:04:05") + " "
	}
	
	result += fmt.Sprintf("[%s] %s", level.String(), message)
	
	if len(fields) > 0 {
		result += " | "
		for k, v := range fields {
			result += fmt.Sprintf("%s=%v", k, v)
		}
	}
	
	return result
}

// CustomFormatter permite formatação customizada
type CustomFormatter struct {
	Template string
}

// Format usa um template customizado
func (f *CustomFormatter) Format(level Level, message string, fields map[string]interface{}) string {
	// Template simples com placeholders
	result := f.Template
	result = fmt.Sprintf(result, 
		time.Now().Format("2006-01-02 15:04:05"),
		level.String(),
		message,
	)
	
	if len(fields) > 0 {
		fieldsStr := ""
		for k, v := range fields {
			fieldsStr += fmt.Sprintf("%s=%v ", k, v)
		}
		result += " | " + fieldsStr
	}
	
	return result
} 