package resources

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPResource_GetName(t *testing.T) {
	config := HttpResourceConfig{
		Name: "test-resource",
		Url:  "http://example.com",
	}
	resource := NewHTTPResource(config, slog.Default())
	if resource.GetName() != "test-resource" {
		t.Errorf("Expected GetName() to return 'test-resource', got '%s'", resource.GetName())
	}
}

func TestHTTPResource_GetType(t *testing.T) {
	config := HttpResourceConfig{
		Name: "test-resource",
		Url:  "http://example.com",
	}
	resource := NewHTTPResource(config, slog.Default())
	if resource.GetType() != "http" {
		t.Errorf("Expected GetType() to return 'http', got '%s'", resource.GetType())
	}
}

func TestHTTPResource_RunCheck_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := HttpResourceConfig{
		Name:           "test-resource",
		Url:            server.URL,
		ExpectedStatus: http.StatusOK,
	}
	resource := NewHTTPResource(config, slog.Default())

	success, details := resource.RunCheck()
	if !success {
		t.Error("Expected RunCheck() to return true for successful HTTP request")
	}
	if len(details) != 0 {
		t.Errorf("Expected no details for successful check, got %d", len(details))
	}
}

func TestHTTPResource_RunCheck_ConnectionError(t *testing.T) {
	config := HttpResourceConfig{
		Name:           "test-resource",
		Url:            "http://nonexistent-domain",
		ExpectedStatus: http.StatusOK,
		MaxRetries:     1,
	}
	resource := NewHTTPResource(config, slog.Default())

	success, details := resource.RunCheck()
	if success {
		t.Error("Expected RunCheck() to return false for connection error")
	}
	if len(details) == 0 {
		t.Error("Expected details for connection error, got none")
	}
}

func TestHTTPResource_RunCheck_UnexpectedStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	config := HttpResourceConfig{
		Name:           "test-resource",
		Url:            server.URL,
		ExpectedStatus: http.StatusOK,
		MaxRetries:     1,
	}
	resource := NewHTTPResource(config, slog.Default())

	success, details := resource.RunCheck()
	if success {
		t.Error("Expected RunCheck() to return false for unexpected status code")
	}
	if len(details) != 3 {
		t.Errorf("Expected 3 details for unexpected status, got %d", len(details))
	}
}
