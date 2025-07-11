package echo

import (
	"github.com/LissaiDev/Delphos/pkg/hermes"
	"github.com/LissaiDev/Delphos/pkg/logger"
)

type Echo struct {
	Handlers []Handler
}

func (d *Echo) Notify(message string) error {

	for _, handler := range d.Handlers {
		handler.Handle(message)
	}

	return nil
}

func (d *Echo) AddHandler(handler Handler) {
	d.Handlers = append(d.Handlers, handler)
}

func NewEcho() Notifier {
	log := logger.Log
	net := hermes.NewHermesClient(log, 3, 10)

	return &Echo{
		Handlers: []Handler{
			NewDiscordHandler(net),
		},
	}
}
