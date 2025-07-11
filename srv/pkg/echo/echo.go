package echo

import (
	"sync"
	"time"

	"github.com/LissaiDev/Delphos/pkg/hermes"
)

var (
	echoInstance Notifier
	once         sync.Once
)

type Echo struct {
	Handlers         []Handler
	cooldown         time.Duration
	lastNotification time.Time
	mu               sync.Mutex
}

func (d *Echo) ShouldNotify() bool {
	if time.Since(d.lastNotification) < d.cooldown {
		return false
	}
	return true
}

func (d *Echo) Notify(message string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.ShouldNotify() {
		for _, handler := range d.Handlers {
			handler.Handle(message)
		}
		d.lastNotification = time.Now()
	}

	return nil
}

func (d *Echo) AddHandler(handler Handler) {
	d.Handlers = append(d.Handlers, handler)
}

func New() Notifier {
	net := hermes.GetInstance()

	return &Echo{
		Handlers: []Handler{
			NewDiscordHandler(net),
		},
		cooldown: 30 * time.Second,
	}
}

func GetInstance() Notifier {
	once.Do(func() {
		echoInstance = New()
	})
	return echoInstance
}
