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

// GetRetries implements Monitor.
func (c CronMonitor) GetRetries() int {
	return c.config.Retries
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
func (c CronMonitor) Run(
	resources []resources.Resource, notificationChannels []chan status.CheckResult) {
	scheduler := cron.New()

	for _, resource := range resources {
		checkFunc := runResourceChecks(resource, notificationChannels, c.GetRetries(), c.logger)
		// checkFunc()

		scheduler.AddFunc(c.config.Cron, checkFunc)
	}

	c.logger.Printf("Starting CRON scheduler")
	scheduler.Start()
}
