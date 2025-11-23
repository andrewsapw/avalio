package monitors

import (
	"log/slog"
	"time"

	"github.com/andrewsapw/avalio/resources"
	"github.com/andrewsapw/avalio/status"
)

type MonitorRunner struct {
	monitor            Monitor
	resource           resources.Resource
	isLastMessageError bool
	channels           []chan status.CheckResult
	logger             *slog.Logger
}

func NewMonitorRunner(monitor Monitor, resource resources.Resource, channels []chan status.CheckResult, logger *slog.Logger) *MonitorRunner {
	return &MonitorRunner{monitor: monitor, resource: resource, logger: logger, isLastMessageError: false}
}

func (m *MonitorRunner) Run() {
	resourceName := m.resource.GetName()
	m.logger.Info("Starting resource monitor", "monitor_name", m.monitor.GetName(),
		"resource_name", resourceName)

	for {
		m.logger.Debug("Checking resource", slog.String("resourceName", resourceName))
		checkResult := m.Step()

		for _, c := range m.channels {
			c <- checkResult
		}

		nextStepAt := m.monitor.Next()
		sleepTime := time.Until(nextStepAt)
		m.logger.Debug("Check result sent to notificators",
			"state", checkResult.State,
			"next_run", nextStepAt)
		time.Sleep(sleepTime)
	}
}

func (m *MonitorRunner) Step() status.CheckResult {
	resourceName := m.resource.GetName()

	ok, details := m.resource.RunCheck()
	if !ok {
		if !m.isLastMessageError {
			m.isLastMessageError = true
			checkResult := status.NewCheckResult(resourceName, details, status.StateNotAvailable)
			return checkResult
		} else {
			checkResult := status.NewCheckResult(resourceName, details, status.StateStillNotAvailable)
			return checkResult
		}
	} else {
		m.logger.Debug("Resource is available", "resource_name", resourceName)
		details := []status.CheckDetails{}

		if m.isLastMessageError {
			m.isLastMessageError = false
			return status.NewCheckResult(resourceName, details, status.StateRecovered)
		} else {
			return status.NewCheckResult(resourceName, details, status.StateAvailable)
		}
	}
}
