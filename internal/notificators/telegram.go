package notificators

import (
	"github.com/andrewsapw/avalio/internal/status"
	"log"
)

type TelegraNotificator struct {
	config TelegramNotificatorConfig
	logger *log.Logger
}

// Send implements Notificator.
func (t TelegraNotificator) Send(status.CheckResult) {
	panic("unimplemented")
}

// GetName implements Notificator.
func (t TelegraNotificator) GetName() string {
	return t.config.Name
}

func NewTelegramNotificator(config TelegramNotificatorConfig, logger *log.Logger) Notificator {
	return TelegraNotificator{config: config, logger: logger}
}
