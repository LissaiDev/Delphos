package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/LissaiDev/Delphos/internal/config"
	"github.com/LissaiDev/Delphos/internal/monitor"
	"github.com/LissaiDev/Delphos/pkg/logger"
)

// SystemStatsHandler handles HTTP requests for system monitoring statistics
// Returns JSON response with current system metrics
func SystemStatsHandler(w http.ResponseWriter, r *http.Request) {
	logger.Log.Info("Generating system statistics", map[string]interface{}{
		"endpoint": "/api/stats",
		"format":   "JSON",
	})

	// Generate system statistics
	startTime := time.Now()
	stats, err := monitor.GetSystemStats()
	generationTime := time.Since(startTime)

	if err != nil {
		logger.Log.Error("Failed to generate system statistics", map[string]interface{}{
			"error":           err.Error(),
			"generation_time": generationTime.String(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Log successful statistics generation
	logger.Log.Info("System statistics generated successfully", map[string]interface{}{
		"generation_time":    generationTime.String(),
		"hostname":           stats.Host.Hostname,
		"cpu_cores":          len(stats.CPU),
		"disk_partitions":    len(stats.Disk),
		"network_interfaces": len(stats.Network),
	})

	// Set response headers and encode JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		logger.Log.Error("Failed to encode JSON response", map[string]interface{}{
			"error": err.Error(),
		})
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	logger.Log.Info("JSON response sent successfully", map[string]interface{}{
		"endpoint": r.URL.String(),
		"method":   r.Method,
	})
}

// SystemStatsStreamHandler handles Server-Sent Events (SSE) for real-time system monitoring
// Continuously streams system statistics at configured intervals
func SystemStatsStreamHandler(w http.ResponseWriter, r *http.Request) {
	logger.Log.Info("Starting SSE stream for system statistics", map[string]interface{}{
		"endpoint": "/api/stats/sse",
		"interval": config.Env.Interval,
	})

	// Configure SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Verify streaming support
	flusher, ok := w.(http.Flusher)
	if !ok {
		logger.Log.Error("Streaming not supported by response writer", map[string]interface{}{
			"endpoint": r.URL.String(),
		})
		http.Error(w, "Streaming is not supported", http.StatusInternalServerError)
		return
	}

	// Notify client disconnection
	notifyDesconnection := r.Context().Done()
	ticker := time.NewTicker(time.Duration(config.Env.Interval) * time.Second)
	defer ticker.Stop()

	// Stream system stats at configured intervals
	streamCount := 0
	for {
		select {
		case <-notifyDesconnection:
			logger.Log.Info("Client disconnected", map[string]interface{}{
				"endpoint": r.URL.String(),
			})
			return
		case <-ticker.C:
			streamCount++

			logger.Log.Debug("Generating statistics for SSE stream", map[string]interface{}{
				"stream_count": streamCount,
				"interval":     config.Env.Interval,
			})

			// Generate system statistics
			startTime := time.Now()
			stats, err := monitor.GetSystemStats()
			generationTime := time.Since(startTime)

			if err != nil {
				logger.Log.Error("Failed to generate system statistics for SSE", map[string]interface{}{
					"error":           err.Error(),
					"stream_count":    streamCount,
					"generation_time": generationTime.String(),
				})
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Log successful statistics generation
			logger.Log.Debug("Statistics generated for SSE stream", map[string]interface{}{
				"stream_count":        streamCount,
				"generation_time":     generationTime.String(),
				"hostname":            stats.Host.Hostname,
				"cpu_usage":           stats.CPU[0].Usage,
				"memory_used_percent": (stats.Memory.Used / stats.Memory.Total) * 100,
			})

			// Format and send SSE data
			w.Write([]byte("data: "))
			if err := json.NewEncoder(w).Encode(stats); err != nil {
				logger.Log.Error("Failed to encode SSE data", map[string]interface{}{
					"error":        err.Error(),
					"stream_count": streamCount,
				})
				return
			}
			w.Write([]byte("\n\n"))

			flusher.Flush()

			logger.Log.Debug("SSE data sent successfully", map[string]interface{}{
				"stream_count": streamCount,
				"endpoint":     r.URL.String(),
			})

		}
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
