package config

import (
	"os"
	"strconv"
	"time"

	"github.com/LissaiDev/Delphos/pkg/logger"
	"github.com/joho/godotenv"
)

// Service handles configuration management following SRP
type Service struct {
	logger logger.BasicLogger
	env    *Environment
}

// NewService creates a new configuration service
func NewService(log logger.BasicLogger) *Service {
	return &Service{
		logger: log,
		env:    &Environment{},
	}
}

// Load loads configuration from environment variables and .env file
func (s *Service) Load() error {
	s.logger.Info("Loading environment configuration", map[string]interface{}{})

	// Set default values
	s.setDefaults()

	// Try to load .env file (optional)
	if err := s.loadDotEnv(); err != nil {
		s.logger.Debug("No .env file found, using environment variables", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Load from environment variables
	return s.loadFromEnv()
}

// GetConfig returns the current configuration
func (s *Service) GetConfig() *Environment {
	return s.env
}

// GetPort returns the configured port
func (s *Service) GetPort() string {
	return s.env.Port
}

// GetName returns the application name
func (s *Service) GetName() string {
	return s.env.Name
}

// GetInterval returns the monitoring interval as duration
func (s *Service) GetInterval() time.Duration {
	return time.Duration(s.env.Interval) * time.Second
}

// setDefaults sets default configuration values
func (s *Service) setDefaults() {
	s.env.Name = "Delphos Server API"
	s.env.Port = ":8080"
	s.env.Interval = 5
}

// loadDotEnv attempts to load .env file
func (s *Service) loadDotEnv() error {
	err := godotenv.Load()
	if err == nil {
		s.logger.Debug(".env file loaded successfully", map[string]interface{}{})
	}
	return err
}

// loadFromEnv loads configuration from environment variables
func (s *Service) loadFromEnv() error {
	name, nameExists := os.LookupEnv("NAME")
	port, portExists := os.LookupEnv("PORT")
	intervalStr, intervalExists := os.LookupEnv("INTERVAL")

	s.logger.Debug("Environment variables status", map[string]interface{}{
		"NAME_exists":     nameExists,
		"PORT_exists":     portExists,
		"INTERVAL_exists": intervalExists,
	})

	// Load values if they exist
	if nameExists {
		s.env.Name = name
	}

	if portExists {
		s.env.Port = port
	}

	if intervalExists {
		if interval, err := strconv.Atoi(intervalStr); err == nil {
			s.env.Interval = interval
		} else {
			s.logger.Warn("Failed to parse INTERVAL environment variable, using default", map[string]interface{}{
				"value":   intervalStr,
				"error":   err.Error(),
				"default": s.env.Interval,
			})
		}
	}

	s.logger.Info("Configuration loaded", map[string]interface{}{
		"name":     s.env.Name,
		"port":     s.env.Port,
		"interval": s.env.Interval,
	})

	return nil
}

// Validate validates the current configuration
func (s *Service) Validate() error {
	if s.env.Name == "" {
		s.logger.Warn("Application name is empty", map[string]interface{}{})
	}

	if s.env.Port == "" {
		s.logger.Error("Port is required", map[string]interface{}{})
		return ErrInvalidPort
	}

	if s.env.Interval <= 0 {
		s.logger.Error("Interval must be positive", map[string]interface{}{
			"interval": s.env.Interval,
		})
		return ErrInvalidInterval
	}

	return nil
}
