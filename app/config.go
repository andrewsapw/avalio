package app

import (
	"github.com/andrewsapw/avalio/monitors"
	"github.com/andrewsapw/avalio/notificators"
	"github.com/andrewsapw/avalio/resources"
)

type Config struct {
	Resources    resources.ResourcesConfig       `toml:"resources"`
	Notificators notificators.NotificatorsConfig `toml:"notificators"`
	Monitors     monitors.MonitorsConfig         `toml:"monitors"`
}
