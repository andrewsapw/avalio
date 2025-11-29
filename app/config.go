package app

import (
	"github.com/BurntSushi/toml"
	"github.com/andrewsapw/avalio/monitors"
	"github.com/andrewsapw/avalio/notificators"
	"github.com/andrewsapw/avalio/resources"
)

type Config struct {
	LogLevel     string                          `toml:"log_level"`
	Resources    resources.ResourcesConfig       `toml:"resources"`
	Notificators notificators.NotificatorsConfig `toml:"notificators"`
	Monitors     monitors.MonitorsConfig         `toml:"monitors"`
}

func ParseConfig(configPath string) (*Config, error) {
	var config Config

	_, err := toml.DecodeFile(configPath, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
