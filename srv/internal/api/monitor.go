package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/LissaiDev/Delphos/internal/monitor"
	"github.com/LissaiDev/Delphos/pkg/logger"
)

func GetMonitorInfo(w http.ResponseWriter, r *http.Request) {
	logger.Log.Info(fmt.Sprintf("<-- Incoming request: %s", r.URL.String()))
	stats, err := monitor.GetStats()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
	logger.Log.Info(fmt.Sprintf("--> Response sent: %s", r.URL.String()))
}

func GetMonitorInfoSSE(w http.ResponseWriter, r *http.Request) {
	logger.Log.Info(fmt.Sprintf("<-- Incoming request: %s", r.URL.String()))
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("cache-control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming is not supported", http.StatusInternalServerError)
		return
	}

	for {
		stats, err := monitor.GetStats()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(stats)

		flusher.Flush()
		time.Sleep(time.Second * 5)
		logger.Log.Info(fmt.Sprintf("--> Response sent: %s", r.URL.String()))
	}
}
