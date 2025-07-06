package api

import (
	"fmt"
	"net/http"

	"github.com/LissaiDev/Delphos/pkg/logger"
)

// CORSHandler wraps HTTP handlers with CORS headers and request logging
// Provides consistent CORS configuration and request/response logging across all endpoints
func CORSHandler(w http.ResponseWriter, r *http.Request, handler http.HandlerFunc) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	// Handle preflight OPTIONS requests
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Log incoming request and delegate to handler
	logger.Log.Info(fmt.Sprintf("<-- Incoming request: %s", r.URL.String()))
	handler(w, r)
}
