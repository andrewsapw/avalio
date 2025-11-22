package notificators

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/andrewsapw/avalio/status"
)

type TelegramNotificator struct {
	config TelegramNotificatorConfig
	logger *slog.Logger
}

// TelegramResponse represents the structure of Telegram API response
type TelegramResponse struct {
	OK          bool   `json:"ok"`
	Description string `json:"description"`
}

// Send implements Notificator.
func (t TelegramNotificator) Send(checkResult status.CheckResult) {
	isSuccess := len(checkResult.Details) == 0
	if isSuccess {
		return
	}

	var message string
	switch checkResult.State {
	case status.StateNotAvailable:
		message = fmt.Sprintf("❌ Ресурс `%s` недоступен.\n\n%s", checkResult.ResourceName, checkResult.ErorrsAsString())
	case status.StateRecovered:
		message = fmt.Sprintf("✅ Ресурс `%s` снова доступен.", checkResult.ResourceName)
	case status.StateAvailable:
		return
	default:
		message = fmt.Sprintf(
			"Ошибка проверки состояния ресурса '%s'. Код состояния '%d' не поддерживается",
			checkResult.ResourceName,
			checkResult.State,
		)
	}

	t.sendMessage(message)
}

func (t TelegramNotificator) sendMessage(message string) error {
	// Create the request URL
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.config.Token)

	// Create the request body
	requestBody, err := json.Marshal(map[string]any{
		"chat_id":    t.config.ChatID,
		"text":       message,
		"parse_mode": "Markdown",
	})
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %v", err)
	}

	// Send POST request
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	// Parse the response
	var telegramResp TelegramResponse
	if err := json.NewDecoder(resp.Body).Decode(&telegramResp); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	// Check if the message was sent successfully
	if !telegramResp.OK {
		return fmt.Errorf("telegram API error: %s", telegramResp.Description)
	}

	return nil
}

// GetName implements Notificator.
func (t TelegramNotificator) GetName() string {
	return t.config.Name
}

func NewTelegramNotificator(config TelegramNotificatorConfig, logger *slog.Logger) TelegramNotificator {
	return TelegramNotificator{config: config, logger: logger}
}
