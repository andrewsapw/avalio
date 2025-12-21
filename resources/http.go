package resources

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/andrewsapw/avalio/status"
)

type HTTPResource struct {
	config HttpResourceConfig
	logger *slog.Logger
}

// GetName implements Resource.
func (H HTTPResource) GetName() string {
	return H.config.Name
}

func (H HTTPResource) GetType() string {
	return "http"
}

func (H HTTPResource) RunCheck() (bool, []status.CheckDetails) {
	// Use configured max retries, default to 3 if not set
	maxRetries := H.config.MaxRetries
	if maxRetries <= 0 {
		maxRetries = 3
	}

	// Use configured retry delay, default to 1 second if not set
	retryDelay := time.Duration(H.config.RetryDelay) * time.Second
	if retryDelay <= 0 {
		retryDelay = time.Second
	}

	for i := 0; i < maxRetries; i++ {
		if success, details := H.performCheck(); success {
			return success, details
		}
		if i < maxRetries-1 {
			time.Sleep(retryDelay)
		}
	}
	return H.performCheck()
}

func (h HTTPResource) performCheck() (bool, []status.CheckDetails) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	// Use HEAD to avoid downloading the entire body
	resp, err := client.Head(h.config.Url)
	if err != nil {
		var checkErrors [1]status.CheckDetails
		checkErrors[0] = status.NewCheckError("Причина", "Ошибка соединения")
		return false, checkErrors[:]
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode
	if statusCode != h.config.ExpectedStatus {
		var checkErrors [3]status.CheckDetails
		checkErrors[0] = status.NewCheckError("Причина", "Неожиданный статус ответа")
		checkErrors[1] = status.NewCheckError("Статус ответа", strconv.Itoa(resp.StatusCode))
		checkErrors[2] = status.NewCheckError("Ожидаемый статус ответа", strconv.Itoa(h.config.ExpectedStatus))
		return false, checkErrors[:]
	}

	return true, nil
}

func NewHTTPResource(config HttpResourceConfig, logger *slog.Logger) HTTPResource {
	return HTTPResource{config: config, logger: logger}
}
