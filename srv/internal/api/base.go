package api

import (
	"net/http"
	"time"

	"github.com/LissaiDev/Delphos/pkg/logger"
)

// CORSHandler wraps HTTP handlers with CORS headers and request logging
// Provides consistent CORS configuration and request/response logging across all endpoints
func CORSHandler(w http.ResponseWriter, r *http.Request, handler http.HandlerFunc) {
	log := logger.GetInstance()

	startTime := time.Now()

	// Log incoming request details
	log.Info("Incoming HTTP request", map[string]interface{}{
		"method":      r.Method,
		"url":         r.URL.String(),
		"remote_addr": r.RemoteAddr,
		"user_agent":  r.UserAgent(),
		"headers":     len(r.Header),
	})

	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	// Handle preflight OPTIONS requests
	if r.Method == http.MethodOptions {
		log.Info("Handling CORS preflight request", map[string]interface{}{
			"method": r.Method,
			"url":    r.URL.String(),
		})
		w.WriteHeader(http.StatusOK)
		return
	}

	// Log request processing start
	log.Debug("Processing request", map[string]interface{}{
		"method": r.Method,
		"url":    r.URL.String(),
	})

	// Delegate to handler
	handler(w, r)

	// Log request completion
	duration := time.Since(startTime)
	log.Info("Request completed", map[string]interface{}{
		"method":   r.Method,
		"url":      r.URL.String(),
		"duration": duration.String(),
	})
}
