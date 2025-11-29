package monitors

import (
	"context"
	"log/slog"
	"testing"

	"github.com/andrewsapw/avalio/status"
)

type MockedResource struct {
	toFail *bool
}

// GetName implements resources.Resource.
func (m MockedResource) GetName() string {
	return "Mocked Resource"
}

// GetType implements resources.Resource.
func (m MockedResource) GetType() string {
	panic("unimplemented")
}

// RunCheck implements resources.Resource.
func (m MockedResource) RunCheck() (bool, []status.CheckDetails) {
	if *m.toFail {
		return false, nil
	} else {
		return true, nil
	}
}

func checkAndVerifyState(
	t *testing.T,
	runner *MonitorRunner,
	state status.ResourceState,
) {
	checkResult := runner.Step()
	if checkResult.State != state {
		t.Errorf("Expected state to be %d", state)
	}
}

func TestMonitorStates(t *testing.T) {
	channel := make(chan status.CheckResult)
	channels := [1]chan status.CheckResult{channel}

	cronConfig := CronMonitorConfig{
		Cron: "* * * *",
	}

	toFail := false
	resource := MockedResource{toFail: &toFail}
	logger := slog.Default()
	monitor, _ := NewCronMonitor(cronConfig, logger)

	runner := NewMonitorRunner(
		monitor,
		resource,
		channels[:],
		context.Background(),
		logger,
	)

	checkAndVerifyState(t, runner, status.StateAvailable)
	toFail = true
	checkAndVerifyState(t, runner, status.StateNotAvailable)
	checkAndVerifyState(t, runner, status.StateStillNotAvailable)
	toFail = false
	checkAndVerifyState(t, runner, status.StateRecovered)
	checkAndVerifyState(t, runner, status.StateAvailable)
}
