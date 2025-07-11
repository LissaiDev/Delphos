package echo

import (
	"github.com/LissaiDev/Delphos/internal/config"
	"github.com/LissaiDev/Delphos/pkg/hermes"
)

type WebhookData struct {
	url      string
	username string
}

type DiscordHandler struct {
	net         hermes.Fetcher
	webhookData WebhookData
}

func (d *DiscordHandler) Handle(message string) error {
	d.net.Post(hermes.Service("DISCORD"), d.webhookData.url, d.BuildBody(message), nil)
	return nil
}

func (d *DiscordHandler) BuildBody(message string) *map[string]any {
	return &map[string]any{
		"username": d.webhookData.username,
		"content":  message,
	}
}

func NewDiscordHandler(net hermes.Fetcher) Handler {
	return &DiscordHandler{
		net: net,
		webhookData: WebhookData{
			url:      config.Env.WebhookUrl,
			username: config.Env.WebhookUsername,
		},
	}
}
