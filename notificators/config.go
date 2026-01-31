package notificators

import (
	"fmt"
	"log/slog"
	"os"
	"slices"
)

// [[notificator.console]]
// name = 'console'
type ConsoleNotificatorConfig struct {
	Name string `toml:"name"`
}

// [[notificator.telegram]]
// name = 'bot'
// chat_id = '...'
// token = '...'
type TelegramNotificatorConfig struct {
	Name   string `toml:"name"`
	ChatID string `toml:"chat_id"`
	Token  string `toml:"token"`
}

func (c TelegramNotificatorConfig) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("[[notificator.telegram]] - name can't be empty")
	}

	if c.ChatID == "" {
		return fmt.Errorf("[[notificator.telegram]] - chat_id can't be empty")
	}

	if c.Token == "" {
		return fmt.Errorf("[[notificator.telegram]] - token can't be empty")
	}

	return nil
}

type NotificatorsConfig struct {
	Console  []ConsoleNotificatorConfig  `toml:"console"`
	Telegram []TelegramNotificatorConfig `toml:"telegram"`
}

func BuildNotificators(config *NotificatorsConfig) ([]Notificator, error) {
	var buildedNotificators []Notificator
	notificatorsNames := []string{}

	for _, consoleNotificatorConfig := range config.Console {
		consoleNotificator := NewConsoleNotificator(consoleNotificatorConfig)
		if slices.Contains(notificatorsNames, consoleNotificator.GetName()) {
			return nil, fmt.Errorf("Duplicated notificators names: %s", consoleNotificator.GetName())
		}

		slog.Info("Builded notificator", "notificatorName", consoleNotificator.GetName())
		buildedNotificators = append(buildedNotificators, consoleNotificator)

		notificatorsNames = append(notificatorsNames, consoleNotificator.GetName())
	}

	for _, telegramNotificatorConfig := range config.Telegram {
		if err := telegramNotificatorConfig.Validate(); err != nil {
			return nil, err
		}
		telegramNotificator := NewTelegramNotificator(telegramNotificatorConfig)
		if slices.Contains(notificatorsNames, telegramNotificator.GetName()) {
			slog.Error("Duplicated notificators names", "duplicated_names", telegramNotificator.GetName())
			os.Exit(1)
		}
		slog.Info("Builded notificator", "notificator_name", telegramNotificator.GetName())
		buildedNotificators = append(buildedNotificators, telegramNotificator)
	}

	return buildedNotificators, nil
}
