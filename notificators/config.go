package notificators

import "fmt"

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
