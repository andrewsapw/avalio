package app

import (
	"log"
	"sync"

	"github.com/andrewsapw/avalio/internal/monitors"
	"github.com/andrewsapw/avalio/internal/notificators"
	"github.com/andrewsapw/avalio/internal/resources"
	"github.com/andrewsapw/avalio/internal/status"
)

type Application struct {
	Resources    []resources.Resource
	Notificators []notificators.Notificator
	Monitors     []monitors.Monitor
	Logger       *log.Logger
}

func NewApplication(
	resources []resources.Resource,
	notificators []notificators.Notificator,
	monitors []monitors.Monitor,
	logger *log.Logger,
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
		channel, _ := notificatorsChannels[notificator.GetName()]
		app.Logger.Printf("Starting notificator '%s'", notificator.GetName())
		go app.listenNotificator(notificator, channel)
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
				app.Logger.Fatalf("Resource '%s' not found", rName)
			}

		}

		// filter notificaotr channels
		var monitorChannels []chan status.CheckResult
		for _, nName := range m.GetNotificatorsNames() {
			notificator, exists := notificatorsChannels[nName]
			if exists {
				monitorChannels = append(monitorChannels, notificator)
			} else {
				app.Logger.Fatalf("Notificator '%s' not found", nName)
			}
		}

		m.Run(monitorResources, monitorChannels)
	}

	app.Logger.Println("Application started")

	wg.Add(1)
	wg.Wait()
}

func (app Application) listenNotificator(notificator notificators.Notificator, channel <-chan status.CheckResult) {
	for {
		checkResult := <-channel
		notificator.Send(checkResult)
	}
}
