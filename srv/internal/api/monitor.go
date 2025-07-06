package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/LissaiDev/Delphos/internal/config"
	"github.com/LissaiDev/Delphos/internal/monitor"
	"github.com/LissaiDev/Delphos/pkg/logger"
)

// SystemStatsHandler handles HTTP requests for system monitoring statistics
// Returns JSON response with current system metrics
func SystemStatsHandler(w http.ResponseWriter, r *http.Request) {
	logger.Log.Info(fmt.Sprintf("--> Response sent: %s", r.URL.String()))

	stats, err := monitor.GetSystemStats()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// SystemStatsStreamHandler handles Server-Sent Events (SSE) for real-time system monitoring
// Continuously streams system statistics at configured intervals
func SystemStatsStreamHandler(w http.ResponseWriter, r *http.Request) {
	// Configure SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Verify streaming support
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming is not supported", http.StatusInternalServerError)
		return
	}

	// Stream system stats at configured intervals
	for {
		stats, err := monitor.GetSystemStats()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Format SSE data
		w.Write([]byte("data: "))
		json.NewEncoder(w).Encode(stats)
		w.Write([]byte("\n\n"))

		flusher.Flush()
		time.Sleep(time.Second * time.Duration(config.Env.Interval))
		logger.Log.Info(fmt.Sprintf("--> Response sent: %s", r.URL.String()))
	}
}

// WrappedSystemStatsHandler provides CORS-enabled access to system statistics
// Uses CORSHandler wrapper for consistent CORS handling
func WrappedSystemStatsHandler(w http.ResponseWriter, r *http.Request) {
	CORSHandler(w, r, SystemStatsHandler)
}

// WrappedSystemStatsStreamHandler provides CORS-enabled access to real-time system statistics
// Uses CORSHandler wrapper for consistent CORS handling
func WrappedSystemStatsStreamHandler(w http.ResponseWriter, r *http.Request) {
	CORSHandler(w, r, SystemStatsStreamHandler)
}
