package resources

import "github.com/andrewsapw/avalio/status"

type Resource interface {
	GetName() string
	GetType() string
	RunCheck() (bool, []status.CheckDetails)
}
