package config

import (
	"os"
	"strconv"

	"github.com/LissaiDev/Delphos/pkg/logger"
	"github.com/joho/godotenv"
)

// Global environment configuration instance
var Env Environment

// LoadEnvironment loads configuration from environment variables
// Falls back to default values if environment variables are not set
func LoadEnvironment() {
	logger.Log.Info("Loading environment configuration")

	defaultConfig := Environment{
		Name:     "Delphos Server API",
		Port:     ":8080",
		Interval: 5,
	}

	// Try to load .env file (optional)
	if err := godotenv.Load(); err != nil {
		logger.Log.Debug("No .env file found, using environment variables", map[string]interface{}{
			"error": err.Error(),
		})
	} else {
		logger.Log.Debug(".env file loaded successfully")
	}

	// Load configuration from environment variables
	name, nameExists := os.LookupEnv("NAME")
	port, portExists := os.LookupEnv("PORT")
	intervalStr, intervalExists := os.LookupEnv("INTERVAL")

	logger.Log.Debug("Environment variables status", map[string]interface{}{
		"NAME_exists":     nameExists,
		"PORT_exists":     portExists,
		"INTERVAL_exists": intervalExists,
	})

	// Use environment variables if all are present, otherwise use defaults
	if nameExists && portExists && intervalExists {
		interval, err := strconv.Atoi(intervalStr)
		if err != nil {
			logger.Log.Warn("Failed to parse INTERVAL environment variable, using default", map[string]interface{}{
				"value":   intervalStr,
				"error":   err.Error(),
				"default": defaultConfig.Interval,
			})
			// If interval parsing fails, use default
			Env = defaultConfig
			logger.Log.Info("Using default configuration", map[string]interface{}{
				"name":     defaultConfig.Name,
				"port":     defaultConfig.Port,
				"interval": defaultConfig.Interval,
			})
			return
		}

		Env = Environment{
			Name:     name,
			Port:     port,
			Interval: interval,
		}

		logger.Log.Info("Configuration loaded from environment variables", map[string]interface{}{
			"name":     name,
			"port":     port,
			"interval": interval,
		})
	} else {
		Env = defaultConfig
		logger.Log.Info("Using default configuration (environment variables not set)", map[string]interface{}{
			"name":     defaultConfig.Name,
			"port":     defaultConfig.Port,
			"interval": defaultConfig.Interval,
		})
	}
}

// Initialize configuration on package import
func init() {
	LoadEnvironment()
}
