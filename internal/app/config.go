package app

import (
	"github.com/andrewsapw/avalio/internal/monitors"
	"github.com/andrewsapw/avalio/internal/notificators"
	"github.com/andrewsapw/avalio/internal/resources"
)

type Config struct {
	Resources    resources.ResourcesConfig       `toml:"resources"`
	Notificators notificators.NotificatorsConfig `toml:"notificators"`
	Monitors     monitors.MonitorsConfig         `toml:"monitors"`
}
