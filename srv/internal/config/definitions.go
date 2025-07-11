package config

import "errors"

// Environment holds the application configuration settings
// Loaded from environment variables with sensible defaults
type Environment struct {
	Name            string  // Application name
	Port            string  // Server port (e.g., ":8080")
	Interval        int     // Monitoring interval in seconds
	CPUThreshold    float64 // Threshold de uso de CPU em porcentagem
	MemoryThreshold float64 // Threshold de uso de memória em porcentagem
	DiskThreshold   float64 // Threshold de uso de disco em porcentagem
}

// Configuration errors
var (
	ErrInvalidPort     = errors.New("invalid port configuration")
	ErrInvalidInterval = errors.New("invalid interval configuration")
	ErrInvalidName     = errors.New("invalid name configuration")
)
