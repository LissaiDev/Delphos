package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/LissaiDev/Delphos/internal/monitor"
	"github.com/LissaiDev/Delphos/pkg/logger"
)

// SystemStatsHandler handles HTTP requests for system monitoring statistics
// Returns JSON response with current system metrics
func SystemStatsHandler(w http.ResponseWriter, r *http.Request) {
	log := logger.GetInstance()
	log.Info("Generating system statistics", map[string]interface{}{
		"endpoint": "/api/stats",
		"format":   "JSON",
	})

	// Generate system statistics
	startTime := time.Now()
	stats, err := monitor.GetSystemStats()
	generationTime := time.Since(startTime)

	if err != nil {
		log.Error("Failed to generate system statistics", map[string]interface{}{
			"error":           err.Error(),
			"generation_time": generationTime.String(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Log successful statistics generation
	log.Info("System statistics generated successfully", map[string]interface{}{
		"generation_time":    generationTime.String(),
		"hostname":           stats.Host.Hostname,
		"cpu_cores":          len(stats.CPU),
		"disk_partitions":    len(stats.Disk),
		"network_interfaces": len(stats.Network),
	})

	// Set response headers and encode JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		log.Error("Failed to encode JSON response", map[string]interface{}{
			"error": err.Error(),
		})
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Info("JSON response sent successfully", map[string]interface{}{
		"endpoint": r.URL.String(),
		"method":   r.Method,
	})
}
