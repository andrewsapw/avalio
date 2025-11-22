package notificators

import (
	"log/slog"

	"github.com/andrewsapw/avalio/status"
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
	c.logger.Debug(
		"Got check result for resource",
		"resource", checkResult.ResourceName,
		"details", checkResult.Details,
	)
}

func NewConsoleNotificator(config ConsoleNotificatorConfig, logger *slog.Logger) Notificator {
	return ConsoleNotificator{config: config, logger: logger}
}
