package app

import (
	"log/slog"
	"os"
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
	Logger       *slog.Logger
}

func NewApplication(
	resources []resources.Resource,
	notificators []notificators.Notificator,
	monitors []monitors.Monitor,
	logger *slog.Logger,
) *Application {
	return &Application{
		Resources:    resources,
		Notificators: notificators,
		Monitors:     monitors,
		Logger:       logger,
	}
}

func (app *Application) Run() {
	var wg sync.WaitGroup

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
		go app.listenNotificator(notificator, channel, app.Logger)
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
				app.Logger.Error("Resource not found", "resource_name", rName)
				os.Exit(1)
			}

		}

		// filter notificaotr channels
		var monitorChannels []chan status.CheckResult
		for _, nName := range m.GetNotificatorsNames() {
			notificator, exists := notificatorsChannels[nName]
			if exists {
				monitorChannels = append(monitorChannels, notificator)
			} else {
				app.Logger.Error("Notificator not found", "notificator_name", nName)
				os.Exit(1)
			}
		}

		for _, r := range monitorResources {
			runner := monitors.NewMonitorRunner(
				m,
				r,
				monitorChannels,
				app.Logger,
			)

			go runner.Run()
		}

	}

	app.Logger.Info("Application started")

	wg.Add(1)
	wg.Wait()
}

func (app Application) listenNotificator(notificator notificators.Notificator, channel <-chan status.CheckResult, logger *slog.Logger) {
	logger.Info("Starting notificator", "notificator_name", notificator.GetName())
	for {
		checkResult := <-channel
		notificator.Send(checkResult)
		time.Sleep(time.Second * 1)
	}
}
