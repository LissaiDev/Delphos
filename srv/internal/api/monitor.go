package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/LissaiDev/Delphos/internal/monitor"
)

func GetMonitorInfo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Getting monitor info...")
	stats, err := monitor.GetStats()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
