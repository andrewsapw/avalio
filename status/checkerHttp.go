package status

type HTTPChecker struct {
	ResourceName string
}

func NewHTTPChecker(resourceName string) *HTTPChecker {
	return &HTTPChecker{ResourceName: resourceName}
}

func (H HTTPChecker) CheckStatus() Status {
	return NewStatus(false, H.ResourceName)
}
