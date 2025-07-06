package main

import (
	"fmt"
	"net/http"

	"github.com/LissaiDev/Delphos/internal/api"
	"github.com/LissaiDev/Delphos/internal/config"
	"github.com/LissaiDev/Delphos/pkg/logger"
)

func main() {
	logger.Log.Info(fmt.Sprintf("\nServer Started\nServer Name: %s\nServer Running on Port: %s", config.Env.Name, config.Env.Port))
	http.HandleFunc("/api/stats", api.GetMonitorInfo)
	http.HandleFunc("/api/stats/sse", api.GetMonitorInfoSSE)
	http.ListenAndServe(config.Env.Port, nil)
}
