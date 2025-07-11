package echo

import (
	"fmt"

	"github.com/LissaiDev/Delphos/internal/config"
	"github.com/LissaiDev/Delphos/pkg/hermes"
	"github.com/LissaiDev/Delphos/pkg/logger"
)

type WebhookData struct {
	url      string
	username string
}

type DiscordHandler struct {
	net         hermes.Fetcher
	webhookData WebhookData
	logger      logger.BasicLogger
}

func (d *DiscordHandler) Handle(message string) error {
	d.logger.Debug("Discord handler processing message", map[string]interface{}{
		"message_length": len(message),
		"webhook_url":    d.webhookData.url,
		"username":       d.webhookData.username,
	})

	if d.webhookData.url == "" {
		d.logger.Warn("Discord webhook URL not configured, skipping notification", map[string]interface{}{
			"message": message,
		})
		return nil
	}

	body := d.BuildBody(message)

	d.logger.Debug("Sending Discord webhook", map[string]interface{}{
		"url":      d.webhookData.url,
		"username": d.webhookData.username,
		"content":  message,
	})

	response := d.net.Post(hermes.Service("DISCORD"), d.webhookData.url, body, nil)
	if !response.Success {
		d.logger.Error("Failed to send Discord webhook", map[string]interface{}{
			"status_code": response.Code,
			"url":         d.webhookData.url,
			"message":     message,
		})
		return fmt.Errorf("discord webhook failed with status code: %d", response.Code)
	}

	d.logger.Info("Discord webhook sent successfully", map[string]interface{}{
		"url":         d.webhookData.url,
		"username":    d.webhookData.username,
		"content":     message,
		"status_code": response.Code,
	})

	return nil
}

func (d *DiscordHandler) BuildBody(message string) *map[string]any {
	return &map[string]any{
		"username": d.webhookData.username,
		"content":  message,
	}
}

func NewDiscordHandler(net hermes.Fetcher) Handler {
	log := logger.GetInstance()

	log.Info("Creating Discord handler", map[string]interface{}{
		"webhook_url":  config.Env.WebhookUrl,
		"webhook_user": config.Env.WebhookUsername,
	})

	return &DiscordHandler{
		net: net,
		webhookData: WebhookData{
			url:      config.Env.WebhookUrl,
			username: config.Env.WebhookUsername,
		},
		logger: log,
	}
}
