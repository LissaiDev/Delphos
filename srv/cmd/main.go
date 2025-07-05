package main

import (
	"fmt"
	"net/http"

	"github.com/LissaiDev/Delphos/internal/api"
	"github.com/LissaiDev/Delphos/internal/config"
	"github.com/LissaiDev/Delphos/pkg/logger"
)

func main() {
	logger.Log.Info(fmt.Sprintf("Server started on port %s", config.Env.Port))
	http.HandleFunc("/api/stats", api.GetMonitorInfo)
	http.ListenAndServe(config.Env.Port, nil)
}
