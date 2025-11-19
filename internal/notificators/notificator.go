package notificators

import (
	"github.com/andrewsapw/avalio/internal/status"
)

type Notificator interface {
	Send(status.CheckResult)
	GetName() string
}
