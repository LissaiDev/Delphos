package config

import (
	"github.com/LissaiDev/Delphos/pkg/logger"
)

// Global environment configuration instance
var Env Environment

// LoadEnvironment loads configuration from environment variables
// Falls back to default values if environment variables are not set
// Maintains backward compatibility
func LoadEnvironment() {
	log := logger.GetInstance()
	configService := GetInstance()
	if err := configService.Load(); err != nil {
		log.Error("Failed to load configuration", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	if err := configService.Validate(); err != nil {
		log.Error("Configuration validation failed", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// Update global Env variable for backward compatibility
	Env = *configService.GetConfig()
}

// GetConfigService returns the configuration service instance

// Initialize configuration on package import
func init() {
	LoadEnvironment()
}
