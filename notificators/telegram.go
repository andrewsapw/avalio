package notificators

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
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
func (t TelegramNotificator) Send(checkResult status.CheckResult) error {
	var message string

	switch checkResult.State {
	case status.StateNotAvailable:
		messageDetails := []string{
			fmt.Sprintf("Тип проверки: `%s`", checkResult.ResourceType),
			checkResult.ErrorsAsString(),
		}
		message = fmt.Sprintf(
			"❌ Ресурс `%s` недоступен.\n\n%s",
			checkResult.ResourceName,
			strings.Join(messageDetails, "\n"),
		)
	case status.StateRecovered:
		message = fmt.Sprintf("✅ Ресурс `%s` снова доступен.", checkResult.ResourceName)
	case status.StateAvailable:
		return nil
	case status.StateStillNotAvailable:
		return nil
	default:
		t.logger.Warn(
			"Ошибка проверки состояния ресурса. Код состояния не поддерживается",
			"resource_name", checkResult.ResourceName,
			"state", checkResult.State,
		)
	}

	err := t.sendMessage(message)
	return err
}

func (t TelegramNotificator) sendMessage(message string) error {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.config.Token)

	client := http.Client{
		Timeout: 10 * time.Second,
	}

	requestBody, err := json.Marshal(map[string]any{
		"chat_id":    t.config.ChatID,
		"text":       message,
		"parse_mode": "Markdown",
	})
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %v", err)
	}

	resp, err := client.Post(apiURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	var telegramResp TelegramResponse
	if err := json.NewDecoder(resp.Body).Decode(&telegramResp); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

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
