package resources

import "log/slog"

// [[resources.http]]
// url = 'https://example.com'
// name = 'example'
type HttpResourceConfig struct {
	Url            string `toml:"url"`
	Name           string `toml:"name"`
	ExpectedStatus int    `toml:"expected_status"`
}

type ResourcesConfig struct {
	Http []HttpResourceConfig `toml:"http"`
}

func BuildResources(config *ResourcesConfig, logger *slog.Logger) ([]Resource, error) {
	var buildedResources []Resource
	for _, httpResourceConfig := range config.Http {
		httpResource := NewHTTPResource(httpResourceConfig, logger)
		logger.Info("Builded resource", "resource_name", httpResource.GetName())
		buildedResources = append(buildedResources, httpResource)
	}

	return buildedResources, nil
}
