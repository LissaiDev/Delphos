package echo

import (
	"sync"
	"time"

	"github.com/LissaiDev/Delphos/pkg/hermes"
	"github.com/LissaiDev/Delphos/pkg/logger"
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
	logger           logger.BasicLogger
}

func (d *Echo) ShouldNotify() bool {
	timeSince := time.Since(d.lastNotification)
	shouldNotify := timeSince >= d.cooldown

	d.logger.Debug("Checking if should notify", map[string]interface{}{
		"time_since_last": timeSince.String(),
		"cooldown":        d.cooldown.String(),
		"should_notify":   shouldNotify,
	})

	return shouldNotify
}

func (d *Echo) Notify(message string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.logger.Info("Notification request received", map[string]interface{}{
		"message_length": len(message),
		"handlers_count": len(d.Handlers),
	})

	if d.ShouldNotify() {
		d.logger.Info("Processing notification", map[string]interface{}{
			"message": message,
		})

		for i, handler := range d.Handlers {
			d.logger.Debug("Sending to handler", map[string]interface{}{
				"handler_index": i,
				"handler_type":  getHandlerType(handler),
			})

			if err := handler.Handle(message); err != nil {
				d.logger.Error("Handler failed to process notification", map[string]interface{}{
					"handler_index": i,
					"handler_type":  getHandlerType(handler),
					"error":         err.Error(),
				})
			} else {
				d.logger.Debug("Handler processed notification successfully", map[string]interface{}{
					"handler_index": i,
					"handler_type":  getHandlerType(handler),
				})
			}
		}

		d.lastNotification = time.Now()
		d.logger.Info("Notification processed successfully", map[string]interface{}{
			"last_notification": d.lastNotification.Format(time.RFC3339),
		})
	} else {
		timeUntilNext := d.cooldown - time.Since(d.lastNotification)
		d.logger.Debug("Notification skipped due to cooldown", map[string]interface{}{
			"time_until_next":   timeUntilNext.String(),
			"last_notification": d.lastNotification.Format(time.RFC3339),
		})
	}

	return nil
}

func (d *Echo) AddHandler(handler Handler) {
	d.logger.Info("Adding new handler", map[string]interface{}{
		"handler_type":   getHandlerType(handler),
		"total_handlers": len(d.Handlers) + 1,
	})

	d.Handlers = append(d.Handlers, handler)
}

func New() Notifier {
	net := hermes.GetInstance()
	log := logger.GetInstance()

	log.Info("Initializing Echo notification system", map[string]interface{}{
		"cooldown": "30s",
	})

	return &Echo{
		Handlers: []Handler{
			NewDiscordHandler(net),
		},
		cooldown: 30 * time.Second,
		logger:   log,
	}
}

func GetInstance() Notifier {
	once.Do(func() {
		echoInstance = New()
	})
	return echoInstance
}

// getHandlerType returns a string representation of the handler type
func getHandlerType(handler Handler) string {
	switch handler.(type) {
	case *DiscordHandler:
		return "DiscordHandler"
	default:
		return "UnknownHandler"
	}
}
