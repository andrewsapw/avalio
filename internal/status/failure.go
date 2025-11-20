package status

import (
	"fmt"
	"strings"
)

type CheckError struct {
	title       string
	description string
}

type CheckResult struct {
	ResourceName string
	Errors       []CheckError
}

func (c CheckResult) ErorrsAsString() string {
	var b strings.Builder
	for _, e := range c.Errors {
		b.WriteString(fmt.Sprintf("%s: %s\n", e.title, e.description))
	}
	return b.String()
}

func NewCheckError(title, description string) CheckError {
	return CheckError{title: title, description: description}
}

func NewCheckResult(resourceName string, errors []CheckError) CheckResult {
	return CheckResult{ResourceName: resourceName, Errors: errors}
}
