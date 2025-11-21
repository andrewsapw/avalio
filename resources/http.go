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
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	// Use HEAD to avoid downloading the entire body
	resp, err := client.Head(H.config.Url)
	if err != nil {
		var checkErrors [2]status.CheckDetails
		checkErrors[0] = status.NewCheckError("Причина", "Ошибка соединения")
		checkErrors[1] = status.NewCheckError("Ошибка", err.Error())
		return false, checkErrors[:]
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode
	if statusCode != H.config.ExpectedStatus {
		var checkErrors [3]status.CheckDetails
		checkErrors[0] = status.NewCheckError("Причина", "Неожиданный статус ответа")
		checkErrors[1] = status.NewCheckError("Статус ответа", strconv.Itoa(resp.StatusCode))
		checkErrors[2] = status.NewCheckError("Ожидаемый статус ответа", strconv.Itoa(H.config.ExpectedStatus))
		return false, checkErrors[:]
	}

	return true, nil
}

func NewHTTPResource(config HttpResourceConfig, logger *slog.Logger) HTTPResource {
	return HTTPResource{config: config, logger: logger}
}
