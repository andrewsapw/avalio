package resources

import (
	"log/slog"
	"sync"
)

// [[resources.http]]
// name = 'example'
// url = 'https://example.com'
type HttpResourceConfig struct {
	Url            string `toml:"url"`
	Name           string `toml:"name"`
	ExpectedStatus int    `toml:"expected_status"`
}

// [[resources.ping]]
// name = 'example'
// address = 'https://example.com'
type PingResourceConfig struct {
	Address        string `toml:"address"`
	Name           string `toml:"name"`
	TimeoutSeconds uint   `toml:"timeout_seconds"`
}

type ResourcesConfig struct {
	Http []HttpResourceConfig `toml:"http"`
	Ping []PingResourceConfig `toml:"ping"`
}

func BuildResources(config *ResourcesConfig, logger *slog.Logger) ([]Resource, error) {
	var buildedResources []Resource

	for _, httpResourceConfig := range config.Http {
		httpResource := NewHTTPResource(httpResourceConfig, logger)
		logger.Info("Builded resource", "resource_name", httpResource.GetName())
		buildedResources = append(buildedResources, httpResource)
	}

	pingMutex := sync.Mutex{}
	for _, pingResourceConfig := range config.Ping {
		pingResource := NewPingResource(pingResourceConfig, &pingMutex, logger)
		logger.Info("Builded resource", "resource_name", pingResource.GetName())
		buildedResources = append(buildedResources, pingResource)
	}

	return buildedResources, nil
}
