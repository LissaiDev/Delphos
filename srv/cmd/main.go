package main

import (
	"github.com/LissaiDev/Delphos/internal/application"
	"github.com/LissaiDev/Delphos/pkg/logger"
)

func main() {
	log := logger.GetInstance()
	app := application.GetInstance()
	if err := app.Start(); err != nil {
		log.Fatal("Application failed to start", map[string]interface{}{
			"error": err.Error(),
		})
	}
}
