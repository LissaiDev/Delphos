package config

import (
	"github.com/LissaiDev/Delphos/pkg/logger"
)

// Global environment configuration instance
var Env Environment

// Global configuration service instance
var configService *Service

// LoadEnvironment loads configuration from environment variables
// Falls back to default values if environment variables are not set
// Maintains backward compatibility
func LoadEnvironment() {
	configService = NewService(logger.Log)
	if err := configService.Load(); err != nil {
		logger.Log.Error("Failed to load configuration", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	if err := configService.Validate(); err != nil {
		logger.Log.Error("Configuration validation failed", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// Update global Env variable for backward compatibility
	Env = *configService.GetConfig()
}

// GetConfigService returns the configuration service instance
func GetConfigService() *Service {
	if configService == nil {
		LoadEnvironment()
	}
	return configService
}

// Initialize configuration on package import
func init() {
	LoadEnvironment()
}
