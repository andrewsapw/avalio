package monitors

import (
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
	checkFunction ResourceCheksFunction,
	channel chan status.CheckResult,
	state status.ResourceState,
) {
	go checkFunction()
	details := <-channel
	if details.State != state {
		t.Errorf("Expected state to be %d", state)
	}
}
func TestMonitorStates(t *testing.T) {
	channel := make(chan status.CheckResult)
	channels := [1]chan status.CheckResult{channel}

	logger := slog.Default()

	toFail := false
	resource := MockedResource{toFail: &toFail}

	checkFunction := runResourceChecks(
		resource,
		channels[:],
		1,
		logger,
	)

	checkAndVerifyState(t, checkFunction, channel, status.StateAvailable)
	toFail = true
	checkAndVerifyState(t, checkFunction, channel, status.StateNotAvailable)
	toFail = false
	checkAndVerifyState(t, checkFunction, channel, status.StateRecovered)
	checkAndVerifyState(t, checkFunction, channel, status.StateAvailable)
}
