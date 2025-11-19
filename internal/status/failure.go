package status

type CheckError struct {
	title       string
	description string
}

type CheckResult struct {
	ResourceName string
	Errors       []CheckError
}

func NewCheckError(title, description string) CheckError {
	return CheckError{title: title, description: description}
}

func NewCheckResult(resourceName string, errors []CheckError) CheckResult {
	return CheckResult{ResourceName: resourceName, Errors: errors}
}
