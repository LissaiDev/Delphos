package config

import (
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/LissaiDev/Delphos/pkg/logger"
	"github.com/joho/godotenv"
)

// Service handles configuration management following SRP
type Service struct {
	logger logger.BasicLogger
	env    *Environment
}

var (
	configService *Service
	once          sync.Once
)

// NewService creates a new configuration service
func New() *Service {
	log := logger.GetInstance()
	return &Service{
		logger: log,
		env:    &Environment{},
	}
}

func GetInstance() *Service {
	once.Do(func() {
		configService = New()
	})
	return configService
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
	s.env.CPUThreshold = 80.0
	s.env.MemoryThreshold = 90.0
	s.env.DiskThreshold = 90.0
	s.env.WebhookUrl = ""
	s.env.WebhookUsername = ""
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
	cpuThresholdStr, cpuThresholdExists := os.LookupEnv("CPU_THRESHOLD")
	memoryThresholdStr, memoryThresholdExists := os.LookupEnv("MEMORY_THRESHOLD")
	diskThresholdStr, diskThresholdExists := os.LookupEnv("DISK_THRESHOLD")
	webhookUrl, webhookUrlExists := os.LookupEnv("WEBHOOK_URL")
	webhookUsername, webhookUsernameExists := os.LookupEnv("WEBHOOK_USERNAME")

	s.logger.Debug("Environment variables status", map[string]interface{}{
		"NAME_exists":             nameExists,
		"PORT_exists":             portExists,
		"INTERVAL_exists":         intervalExists,
		"CPU_THRESHOLD_exists":    cpuThresholdExists,
		"MEMORY_THRESHOLD_exists": memoryThresholdExists,
		"DISK_THRESHOLD_exists":   diskThresholdExists,
		"WEBHOOK_URL_exists":      webhookUrlExists,
		"WEBHOOK_USERNAME_exists": webhookUsernameExists,
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

	if cpuThresholdExists {
		if v, err := strconv.ParseFloat(cpuThresholdStr, 64); err == nil {
			s.env.CPUThreshold = v
		} else {
			s.logger.Warn("Failed to parse CPU_THRESHOLD environment variable, using default", map[string]interface{}{
				"value":   cpuThresholdStr,
				"error":   err.Error(),
				"default": s.env.CPUThreshold,
			})
		}
	}

	if memoryThresholdExists {
		if v, err := strconv.ParseFloat(memoryThresholdStr, 64); err == nil {
			s.env.MemoryThreshold = v
		} else {
			s.logger.Warn("Failed to parse MEMORY_THRESHOLD environment variable, using default", map[string]interface{}{
				"value":   memoryThresholdStr,
				"error":   err.Error(),
				"default": s.env.MemoryThreshold,
			})
		}
	}

	if diskThresholdExists {
		if v, err := strconv.ParseFloat(diskThresholdStr, 64); err == nil {
			s.env.DiskThreshold = v
		} else {
			s.logger.Warn("Failed to parse DISK_THRESHOLD environment variable, using default", map[string]interface{}{
				"value":   diskThresholdStr,
				"error":   err.Error(),
				"default": s.env.DiskThreshold,
			})
		}
	}

	if webhookUrlExists {
		s.env.WebhookUrl = webhookUrl
	}
	if webhookUsernameExists {
		s.env.WebhookUsername = webhookUsername
	}

	s.logger.Info("Configuration loaded", map[string]interface{}{
		"name":             s.env.Name,
		"port":             s.env.Port,
		"interval":         s.env.Interval,
		"cpu_threshold":    s.env.CPUThreshold,
		"memory_threshold": s.env.MemoryThreshold,
		"disk_threshold":   s.env.DiskThreshold,
		"webhook_url":      s.env.WebhookUrl,
		"webhook_username": s.env.WebhookUsername,
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

	if s.env.CPUThreshold < 1 || s.env.CPUThreshold > 100 {
		s.logger.Error("CPU_THRESHOLD must be between 1 and 100", map[string]interface{}{
			"cpu_threshold": s.env.CPUThreshold,
		})
		return ErrInvalidCPUThreshold
	}

	if s.env.MemoryThreshold < 1 || s.env.MemoryThreshold > 100 {
		s.logger.Error("MEMORY_THRESHOLD must be between 1 and 100", map[string]interface{}{
			"memory_threshold": s.env.MemoryThreshold,
		})
		return ErrInvalidMemoryThreshold
	}

	if s.env.DiskThreshold < 1 || s.env.DiskThreshold > 100 {
		s.logger.Error("DISK_THRESHOLD must be between 1 and 100", map[string]interface{}{
			"disk_threshold": s.env.DiskThreshold,
		})
		return ErrInvalidDiskThreshold
	}

	return nil
}
