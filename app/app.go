package app

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/andrewsapw/avalio/monitors"
	"github.com/andrewsapw/avalio/notificators"
	"github.com/andrewsapw/avalio/resources"
	"github.com/andrewsapw/avalio/status"
)

type Application struct {
	Resources    []resources.Resource
	Notificators []notificators.Notificator
	Monitors     []monitors.Monitor
}

func NewApplication(
	resources []resources.Resource,
	notificators []notificators.Notificator,
	monitors []monitors.Monitor,
) *Application {
	return &Application{
		Resources:    resources,
		Notificators: notificators,
		Monitors:     monitors,
	}
}

func (app *Application) Run() error {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create channels for each notificator
	notificatorsChannels := make(map[string]chan status.CheckResult)
	for _, n := range app.Notificators {
		nChannel := make(chan status.CheckResult)
		notificatorsChannels[n.GetName()] = nChannel
	}

	// create resource name to objects mapping
	nameToResource := make(map[string]resources.Resource)
	for _, r := range app.Resources {
		nameToResource[r.GetName()] = r
	}

	// start notificators listen
	for _, notificator := range app.Notificators {
		channel := notificatorsChannels[notificator.GetName()]

		wg.Add(1)
		go func() {
			go app.listenNotificator(notificator, channel, ctx)
			wg.Done()
		}()
	}

	// start monitors
	for _, m := range app.Monitors {
		// filter Resources
		var monitorResources []resources.Resource
		for _, rName := range m.GetResourcesNames() {
			resource, exists := nameToResource[rName]
			if exists {
				monitorResources = append(monitorResources, resource)
			} else {
				return fmt.Errorf("Resource '%s' not found", rName)
			}

		}

		// filter notificaotr channels
		var monitorChannels []chan status.CheckResult
		for _, nName := range m.GetNotificatorsNames() {
			notificator, exists := notificatorsChannels[nName]
			if exists {
				monitorChannels = append(monitorChannels, notificator)
			} else {
				return fmt.Errorf("Notificator '%s' not found", nName)
			}
		}

		for _, r := range monitorResources {
			runner := monitors.NewMonitorRunner(
				m,
				r,
				monitorChannels,
				ctx,
			)

			wg.Add(1)
			go func() {
				runner.Run()
				wg.Done()
			}()
		}

	}

	slog.Info("Application started")

	wg.Wait()
	return nil
}

func (app Application) listenNotificator(
	notificator notificators.Notificator,
	channel <-chan status.CheckResult,
	ctx context.Context,
) {
	notificatorName := notificator.GetName()
	slog.Info("Starting notificator", "notificator_name", notificator.GetName())
	for {
		select {
		case <-ctx.Done():
			return
		case checkResult := <-channel:
			if err := notificator.Send(checkResult); err != nil {
				slog.Error(
					"Error sending notification",
					"notificator_name", notificatorName,
					"error", err,
				)
			}
		}
		time.Sleep(time.Second * 1)
	}
}
