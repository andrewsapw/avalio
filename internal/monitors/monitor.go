package monitors

import (
	"github.com/andrewsapw/avalio/internal/resources"
	"github.com/andrewsapw/avalio/internal/status"
)

type Monitor interface {
	Run(resources []resources.Resource, notificationChannels []chan status.CheckResult)
	GetResourcesNames() []string
	GetNotificatorsNames() []string
}

