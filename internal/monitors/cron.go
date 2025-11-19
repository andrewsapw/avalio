package monitors

import (
	"log"

	"github.com/andrewsapw/avalio/internal/resources"
	"github.com/andrewsapw/avalio/internal/status"
	"github.com/robfig/cron/v3"
)

type CronMonitor struct {
	config CronMonitorConfig
	logger *log.Logger
}

// GetNotificatorsNames implements Monitor.
func (c CronMonitor) GetNotificatorsNames() []string {
	return c.config.Notificators
}

// GetResourcesNames implements Monitor.
func (c CronMonitor) GetResourcesNames() []string {
	return c.config.Resources
}

func NewCronMonitor(config CronMonitorConfig, logger *log.Logger) Monitor {
	return CronMonitor{config: config, logger: logger}
}

// Run implements Monitor.
func (c CronMonitor) Run(resources []resources.Resource, notificationChannels []chan status.CheckResult) {
	scheduler := cron.New()
	for _, resource := range resources {
		scheduler.AddFunc(c.config.Cron, func() { c.runResourceChecks(resource, notificationChannels) })
		go c.runResourceChecks(resource, notificationChannels)
	}

	c.logger.Printf("Starting CRON scheduler")
	scheduler.Start()
}

func (c CronMonitor) runResourceChecks(resource resources.Resource, notificationChannels []chan status.CheckResult) {
	resourceName := resource.GetName()
	c.logger.Printf("Checking resource %s", resourceName)
	errors := resource.CheckErrors()
	if errors != nil || len(errors) > 0 {
		checkResult := status.NewCheckResult(resourceName, errors)
		for _, c := range notificationChannels {
			c <- checkResult
		}
	} else {
		c.logger.Printf("Resource %s is available", resource.GetName())
	}
}
