package status

import (
	"fmt"
	"strings"
)

type CheckDetails struct {
	title       string
	description string
}

type ResourceState int

const (
	StateAvailable         ResourceState = iota // 0
	StateNotAvailable                           // 1
	StateStillNotAvailable                      // 2
	StateRecovered                              // 3
)

func (s ResourceState) String() string {
	switch s {
	case StateAvailable:
		return "available"
	case StateNotAvailable:
		return "not available"
	case StateStillNotAvailable:
		return "still not available"
	case StateRecovered:
		return "recovered"
	default:
		return "unknown"
	}
}

type CheckResult struct {
	ResourceName string
	State        ResourceState
	Details      []CheckDetails
}

func (c CheckResult) ErorrsAsString() string {
	var b strings.Builder
	for _, e := range c.Details {
		b.WriteString(fmt.Sprintf("%s: %s\n", e.title, e.description))
	}
	return b.String()
}

func NewCheckError(title, description string) CheckDetails {
	return CheckDetails{title: title, description: description}
}

func NewCheckResult(resourceName string, details []CheckDetails, state ResourceState) CheckResult {
	return CheckResult{
		ResourceName: resourceName,
		Details:      details,
		State:        state,
	}
}
