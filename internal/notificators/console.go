package notificators

import (
	"log/slog"

	"github.com/andrewsapw/avalio/internal/status"
)

type ConsoleNotificator struct {
	config ConsoleNotificatorConfig
	logger *slog.Logger
}

// GetName implements Notificator.
func (c ConsoleNotificator) GetName() string {
	return c.config.Name
}

// Send implements Notificator.
func (c ConsoleNotificator) Send(checkResult status.CheckResult) {
	c.logger.Info("got check result for resource '%s': %s", checkResult.ResourceName, checkResult.Details)
}

func NewConsoleNotificator(config ConsoleNotificatorConfig, logger *slog.Logger) Notificator {
	return ConsoleNotificator{config: config, logger: logger}
}
