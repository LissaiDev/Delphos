package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/LissaiDev/Delphos/internal/monitor"
)

func GetMonitorInfo(w http.ResponseWriter, r *http.Request) {
	stats, err := monitor.GetStats()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func GetMonitorInfoSSE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("cache-control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streming is not supported", http.StatusInternalServerError)
		return
	}

	for {
		stats, err := monitor.GetStats()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write([]byte("{data: "))
		json.NewEncoder(w).Encode(stats)
		w.Write([]byte("}\n\n"))

		flusher.Flush()
		time.Sleep(time.Second * 5)
	}
}
