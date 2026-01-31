package monitors

import (
	"context"
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
	ctx                context.Context
}

func NewMonitorRunner(
	monitor Monitor,
	resource resources.Resource,
	channels []chan status.CheckResult,
	ctx context.Context,
) *MonitorRunner {
	return &MonitorRunner{monitor: monitor, channels: channels, resource: resource, ctx: ctx, isLastMessageError: false}
}

func (m *MonitorRunner) Run() {
	resourceName := m.resource.GetName()
	slog.Info("Starting resource monitor", "monitor_name", m.monitor.GetName(),
		"resource_name", resourceName)

	for {
		slog.Debug("Checking resource", slog.String("resourceName", resourceName))
		checkResult := m.Step()

		err := m.ctx.Err()
		if err != nil {
			// context is closed
			return
		}

		for _, c := range m.channels {
			c <- checkResult
		}

		nextStepAt := m.monitor.Next()
		sleepTime := time.Until(nextStepAt)
		slog.Debug("Check result sent to notificators",
			"state", checkResult.State,
			"next_run", nextStepAt,
			"resource_name", resourceName)

		time.Sleep(sleepTime)
	}
}

func (m *MonitorRunner) Step() status.CheckResult {
	resourceName := m.resource.GetName()
	resourceType := m.resource.GetType()

	ok, details := m.resource.RunCheck()

	var state status.ResourceState
	if !ok {
		if !m.isLastMessageError {
			m.isLastMessageError = true
			state = status.StateNotAvailable
		} else {
			state = status.StateStillNotAvailable
		}
	} else {
		slog.Debug("Resource is available", "resource_name", resourceName)

		if m.isLastMessageError {
			m.isLastMessageError = false
			state = status.StateRecovered
		} else {
			state = status.StateAvailable
		}
	}

	checkResult := status.NewCheckResult(
		resourceName,
		resourceType,
		details,
		state,
	)
	return checkResult
}
