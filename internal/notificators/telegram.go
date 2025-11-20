package notificators

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/andrewsapw/avalio/internal/status"
)

type TelegraNotificator struct {
	config TelegramNotificatorConfig
	logger *log.Logger
}

// TelegramResponse represents the structure of Telegram API response
type TelegramResponse struct {
	OK          bool   `json:"ok"`
	Description string `json:"description"`
}

// Send implements Notificator.
func (t TelegraNotificator) Send(checkResult status.CheckResult) {
	isSuccess := len(checkResult.Errors) == 0
	if isSuccess {
		return
	}

	message := fmt.Sprintf("❌ Ресурс `%s` недоступен.\n\n%s", checkResult.ResourceName, checkResult.ErorrsAsString())

	t.sendMessage(message)
}

func (t TelegraNotificator) sendMessage(message string) error {
	// Create the request URL
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.config.Token)

	// Create the request body
	requestBody, err := json.Marshal(map[string]interface{}{
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
func (t TelegraNotificator) GetName() string {
	return t.config.Name
}

func NewTelegramNotificator(config TelegramNotificatorConfig, logger *log.Logger) Notificator {
	return TelegraNotificator{config: config, logger: logger}
}
