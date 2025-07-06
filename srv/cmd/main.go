package main

import (
	"net/http"

	"github.com/LissaiDev/Delphos/internal/api"
	"github.com/LissaiDev/Delphos/internal/config"
	"github.com/LissaiDev/Delphos/pkg/logger"
)

func main() {
	// Log server startup information
	logger.Log.Info("Starting Delphos monitoring server", map[string]interface{}{
		"server_name": config.Env.Name,
		"port":        config.Env.Port,
		"interval":    config.Env.Interval,
	})

	// Register HTTP handlers with detailed logging
	logger.Log.Info("Registering HTTP handlers", map[string]interface{}{
		"endpoint_json": "/api/stats",
		"endpoint_sse":  "/api/stats/sse",
	})

	http.HandleFunc("/api/stats", api.WrappedSystemStatsHandler)           // JSON endpoint
	http.HandleFunc("/api/stats/sse", api.WrappedSystemStatsStreamHandler) // SSE endpoint

	// Start HTTP server
	logger.Log.Info("HTTP server starting", map[string]interface{}{
		"address": config.Env.Port,
	})

	if err := http.ListenAndServe(config.Env.Port, nil); err != nil {
		logger.Log.Fatal("Failed to start HTTP server", map[string]interface{}{
			"error": err.Error(),
			"port":  config.Env.Port,
		})
	}
}
