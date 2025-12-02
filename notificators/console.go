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
func (c ConsoleNotificator) Send(checkResult status.CheckResult) error {
	c.logger.Debug(
		"Got check result for resource",
		"state", checkResult.State,
		"resource_name", checkResult.ResourceName,
		"resource_type", checkResult.ResourceType,
		"details", checkResult.ErrorsAsString(),
	)
	return nil
}

func NewConsoleNotificator(config ConsoleNotificatorConfig, logger *slog.Logger) ConsoleNotificator {
	return ConsoleNotificator{config: config, logger: logger}
}
