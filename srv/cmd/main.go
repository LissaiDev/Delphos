package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/LissaiDev/Delphos/internal/api"
	"github.com/LissaiDev/Delphos/internal/config"
	"github.com/LissaiDev/Delphos/internal/monitor"
	"github.com/LissaiDev/Delphos/pkg/logger"
)

func main() {
	broker := api.NewBroker()
	broker.Start()
	defer broker.Stop()

	go func() {
		ticker := time.NewTicker(time.Duration(config.Env.Interval) * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			if len(broker.Clients) == 0 {
				continue
			}

			stats, err := monitor.GetSystemStats()
			if err != nil {
				continue
			}
			bytes, err := json.Marshal(stats)
			if err != nil {
				continue
			}

			logger.Log.Info("Broadcasting stats", map[string]interface{}{
				"stats": string(bytes),
			})

			broker.Broadcast(string(bytes))

		}
	}()

	http.HandleFunc("/api/stats", api.WrappedSystemStatsHandler)
	http.Handle("/api/stats/sse", broker)

	if err := http.ListenAndServe(config.Env.Port, nil); err != nil {
		logger.Log.Fatal("Failed to start HTTP server", map[string]interface{}{
			"error": err.Error(),
			"port":  config.Env.Port,
		})
	}
}
