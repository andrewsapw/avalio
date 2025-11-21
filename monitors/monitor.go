package monitors

import (
	"log/slog"

	"github.com/andrewsapw/avalio/resources"
	"github.com/andrewsapw/avalio/status"
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
	logger *slog.Logger,
) ResourceCheksFunction {
	errorsCounter := 0
	return func() {
		resourceName := resource.GetName()
		logger.Info("Checking resource", slog.String("resourceName", resourceName))

		ok, details := resource.RunCheck()
		if !ok {
			errorsCounter += 1
			if errorsCounter >= maxRetries {
				checkResult := status.NewCheckResult(resourceName, details, status.StateNotAvailable)
				for _, c := range channels {
					c <- checkResult
				}
			} else {
				logger.Warn(
					"Resource is unavailable",
					"resource_name", resource.GetName(),
					"current_errors_counter", errorsCounter,
					"max_retries", maxRetries,
				)
			}

		} else {
			logger.Info("Resource is available", "resource_name", resource.GetName())
			details := []status.CheckDetails{}

			if errorsCounter >= maxRetries {
				for _, c := range channels {
					c <- status.NewCheckResult(resourceName, details, status.StateRecovered)
				}
			}

			errorsCounter = 0
		}
	}
}
