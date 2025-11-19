package notificators

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

type NotificatorsConfig struct {
	Console  []ConsoleNotificatorConfig  `toml:"console"`
	Telegram []TelegramNotificatorConfig `toml:"telegram"`
}
