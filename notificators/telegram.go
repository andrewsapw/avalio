package notificators

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

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

	var message string
	switch checkResult.State {
	case status.StateNotAvailable:
		message = fmt.Sprintf("❌ Ресурс `%s` недоступен.\n\n%s", checkResult.ResourceName, checkResult.ErorrsAsString())
	case status.StateRecovered:
		message = fmt.Sprintf("✅ Ресурс `%s` снова доступен.", checkResult.ResourceName)
	case status.StateAvailable:
		return
	case status.StateStillNotAvailable:
		return
	default:
		t.logger.Warn(
			"Ошибка проверки состояния ресурса. Код состояния не поддерживается",
			"resource_name", checkResult.ResourceName,
			"state", checkResult.State,
		)
	}

	t.sendMessage(message)
}

func (t TelegramNotificator) sendMessage(message string) error {
	// Create the request URL
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.config.Token)
	client := http.Client{
		Timeout: 10 * time.Second,
	}

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
	resp, err := client.Post(apiURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		t.logger.Error("HTTP request failed", "error", err.Error())
		return fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	// Parse the response
	var telegramResp TelegramResponse
	if err := json.NewDecoder(resp.Body).Decode(&telegramResp); err != nil {
		t.logger.Error("failed to decode response", "error", err.Error())
		return fmt.Errorf("failed to decode response: %v", err)
	}

	// Check if the message was sent successfully
	if !telegramResp.OK {
		t.logger.Error("telegram API error", "error", err.Error())
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
