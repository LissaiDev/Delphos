package main

import (
	"fmt"
	"net/http"

	"github.com/LissaiDev/Delphos/internal/api"
	"github.com/LissaiDev/Delphos/internal/config"
	"github.com/LissaiDev/Delphos/pkg/logger"
)

func main() {
	// Log server startup information
	logger.Log.Info(fmt.Sprintf("\nServer Started\nServer Name: %s\nServer Running on Port: %s",
		config.Env.Name, config.Env.Port))

	// Register HTTP handlers
	http.HandleFunc("/api/stats", api.WrappedSystemStatsHandler)           // JSON endpoint
	http.HandleFunc("/api/stats/sse", api.WrappedSystemStatsStreamHandler) // SSE endpoint

	// Start HTTP server
	http.ListenAndServe(config.Env.Port, nil)
}
