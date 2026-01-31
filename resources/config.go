package resources

import (
	"errors"
	"fmt"
	"log/slog"
	"net/url"
)

// Error variables for HTTP resource validation
var (
	HTTPResourceNameIsEmptyError     = errors.New("name is required")
	HTTPResourceURLEmptyError        = errors.New("url is required")
	HTTPResourceInvalidURLError      = errors.New("url is invalid")
	HTTPResourceInvalidSchemeError   = errors.New("url must use http or https scheme")
	HTTPResourceMissingHostError     = errors.New("url must include a host")
	HTTPResourceNegativeStatusError  = errors.New("expected_status must be non-negative")
	HTTPResourceInvalidStatusError   = errors.New("expected_status must be a valid HTTP status code (100-599)")
	HTTPResourceNegativeRetriesError = errors.New("max_retries must be non-negative")
	HTTPResourceHighRetriesError     = errors.New("max_retries must not exceed 10")
	HTTPResourceNegativeDelayError   = errors.New("retry_delay must be non-negative")
	HTTPResourceHighDelayError       = errors.New("retry_delay must not exceed 300 seconds (5 minutes)")
	HTTPResourceLongNameError        = errors.New("name must not exceed 255 characters")
	HTTPResourceLongURLError         = errors.New("url must not exceed 2048 characters")
)

// Error variables for Ping resource validation
var (
	PingResourceNameIsEmptyError  = errors.New("name is required")
	PingResourceAddressEmptyError = errors.New("address is required")
	PingResourceLongNameError     = errors.New("name must not exceed 255 characters")
	PingResourceLongAddressError  = errors.New("address must not exceed 255 characters")
	PingResourceZeroTimeoutError  = errors.New("timeout_seconds must be greater than 0")
	PingResourceHighTimeoutError  = errors.New("timeout_seconds must not exceed 300 seconds (5 minutes)")
)

// [[resources.http]]
// name = 'example'
// url = 'https://example.com'
type HttpResourceConfig struct {
	Url            string `toml:"url"`
	Name           string `toml:"name"`
	ExpectedStatus int    `toml:"expected_status"`
	MaxRetries     int    `toml:"max_retries"`
	RetryDelay     int    `toml:"retry_delay"`
}

// Validate checks if the HTTP resource configuration is valid
func (c *HttpResourceConfig) Validate() error {
	// Validate name
	if c.Name == "" {
		return HTTPResourceNameIsEmptyError
	}

	if len(c.Name) > 255 {
		return HTTPResourceLongNameError
	}

	// Validate URL
	if c.Url == "" {
		return HTTPResourceURLEmptyError
	}

	// Check URL length to prevent extremely long URLs
	if len(c.Url) > 2048 {
		return HTTPResourceLongURLError
	}

	parsedUrl, err := url.Parse(c.Url)
	if err != nil {
		return fmt.Errorf("http resource '%s': invalid url '%s': %w", c.Name, c.Url, HTTPResourceInvalidURLError)
	}

	if parsedUrl.Scheme != "http" && parsedUrl.Scheme != "https" {
		return HTTPResourceInvalidSchemeError
	}

	if parsedUrl.Host == "" {
		return HTTPResourceMissingHostError
	}

	// Validate expected status
	if c.ExpectedStatus < 0 {
		return HTTPResourceNegativeStatusError
	}

	// Validate HTTP status code range (100-599)
	if c.ExpectedStatus > 0 && (c.ExpectedStatus < 100 || c.ExpectedStatus > 599) {
		return HTTPResourceInvalidStatusError
	}

	// Validate max retries
	if c.MaxRetries < 0 {
		return HTTPResourceNegativeRetriesError
	}

	// Set reasonable upper limit for max retries to prevent excessive retries
	if c.MaxRetries > 10 {
		return HTTPResourceHighRetriesError
	}

	// Validate retry delay
	if c.RetryDelay < 0 {
		return HTTPResourceNegativeDelayError
	}

	// Set reasonable upper limit for retry delay to prevent extremely long waits
	if c.RetryDelay > 300 { // 5 minutes max
		return HTTPResourceHighDelayError
	}

	return nil
}

// [[resources.ping]]
// name = 'example'
// address = 'https://example.com'
type PingResourceConfig struct {
	Address        string `toml:"address"`
	Name           string `toml:"name"`
	TimeoutSeconds int    `toml:"timeout_seconds"`
}

// Validate checks if the ping resource configuration is valid
func (c *PingResourceConfig) Validate() error {
	if c.Name == "" {
		return PingResourceNameIsEmptyError
	}

	if len(c.Name) > 255 {
		return PingResourceLongNameError
	}

	if c.Address == "" {
		return PingResourceAddressEmptyError
	}

	// Check address length to prevent extremely long addresses
	if len(c.Address) > 255 {
		return PingResourceLongAddressError
	}

	// Validate timeout (must be positive, with reasonable limits)
	if c.TimeoutSeconds <= 0 {
		return PingResourceZeroTimeoutError
	}

	// Set reasonable upper limit for timeout (300 seconds / 5 minutes max)
	if c.TimeoutSeconds > 300 {
		return PingResourceHighTimeoutError
	}

	return nil
}

type ResourcesConfig struct {
	Http []HttpResourceConfig `toml:"http"`
	Ping []PingResourceConfig `toml:"ping"`
}

func BuildResources(config *ResourcesConfig, logger *slog.Logger) ([]Resource, error) {
	var buildedResources []Resource

	for _, httpResourceConfig := range config.Http {
		if err := httpResourceConfig.Validate(); err != nil {
			return nil, fmt.Errorf("invalid http resource configuration: %w", err)
		}

		httpResource := NewHTTPResource(httpResourceConfig, logger)
		logger.Info("Builded resource", "resource_name", httpResource.GetName())
		buildedResources = append(buildedResources, httpResource)
	}

	for _, pingResourceConfig := range config.Ping {
		if err := pingResourceConfig.Validate(); err != nil {
			return nil, fmt.Errorf("invalid ping resource configuration: %w", err)
		}

		pingResource := NewPingResource(pingResourceConfig, logger)
		logger.Info("Builded resource", "resource_name", pingResource.GetName())
		buildedResources = append(buildedResources, pingResource)
	}

	return buildedResources, nil
}
