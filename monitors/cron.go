package monitors

import (
	"log/slog"
	"time"

	"github.com/robfig/cron/v3"
)

type CronMonitor struct {
	config   CronMonitorConfig
	logger   *slog.Logger
	schedule cron.Schedule
}

// GetName implements Monitor.
func (c *CronMonitor) GetName() string {
	return c.config.Name
}

// GetNotificatorsNames implements Monitor.
func (c CronMonitor) GetNotificatorsNames() []string {
	return c.config.Notificators
}

// GetResourcesNames implements Monitor.
func (c CronMonitor) GetResourcesNames() []string {
	return c.config.Resources
}

func (c CronMonitor) Next() time.Time {
	now := time.Now()
	return c.schedule.Next(now)
}

func NewCronMonitor(config CronMonitorConfig, logger *slog.Logger) (*CronMonitor, error) {
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	schedule, err := parser.Parse(config.Cron)
	if err != nil {
		return nil, err
	}
	return &CronMonitor{config: config, logger: logger, schedule: schedule}, nil

}
