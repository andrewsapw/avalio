package monitors

import (
	"time"
)

type Monitor interface {
	GetName() string
	GetResourcesNames() []string
	GetNotificatorsNames() []string
	Next() time.Time
}

type ResourceCheksFunction func()
