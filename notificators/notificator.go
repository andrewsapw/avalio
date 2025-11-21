package notificators

import (
	"github.com/andrewsapw/avalio/status"
)

type Notificator interface {
	Send(status.CheckResult)
	GetName() string
}
