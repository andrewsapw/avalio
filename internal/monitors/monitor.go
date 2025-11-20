package monitors

import (
	"log"

	"github.com/andrewsapw/avalio/internal/resources"
	"github.com/andrewsapw/avalio/internal/status"
)

type Monitor interface {
	Run(resources []resources.Resource, notificationChannels []chan status.CheckResult)
	GetResourcesNames() []string
	GetNotificatorsNames() []string
	GetRetries() int
}

type ResourceCheksFunction func()

func runResourceChecks(resource resources.Resource,
	channels []chan status.CheckResult,
	maxRetries int,
	logger *log.Logger,
) ResourceCheksFunction {
	errorsCounter := 0
	return func() {
		resourceName := resource.GetName()
		logger.Printf("Checking resource %s", resourceName)

		errors := resource.CheckErrors()
		if len(errors) > 0 {
			errorsCounter += 1
			if errorsCounter >= maxRetries {
				checkResult := status.NewCheckResult(resourceName, errors)
				for _, c := range channels {
					c <- checkResult
				}
			} else {
				logger.Printf(
					"Resource %s is unavailable for %d time in a row (max retries %d)",
					resource.GetName(),
					errorsCounter,
					maxRetries,
				)
			}

		} else {
			errorsCounter = 0
			logger.Printf("Resource %s is available", resource.GetName())
		}
	}
}
