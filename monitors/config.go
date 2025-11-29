package monitors

import "log/slog"

// [[monitors.cron]]
// name = 'every minute'
// resources = ['moninotr1']
// cron = '* * * * *'
// retries = 3

type MonitorConfig struct {
	Name         string   `toml:"name"`
	Resources    []string `toml:"resources"`
	Notificators []string `toml:"notificators"`
}

type CronMonitorConfig struct {
	MonitorConfig
	Cron string `toml:"cron"`
}

type MonitorsConfig struct {
	Cron []CronMonitorConfig `toml:"cron"`
}

func BuildMonitors(config *MonitorsConfig, logger *slog.Logger) ([]Monitor, error) {
	var buildedMonitors []Monitor
	for _, cronMonitorConfig := range config.Cron {
		cronMonitor, err := NewCronMonitor(cronMonitorConfig, logger)
		if err != nil {
			logger.Error("Error creating monitor", "error", err.Error())
			return nil, err
		}

		logger.Info("Builded monitor", "monitor_name", cronMonitorConfig.Name)
		buildedMonitors = append(buildedMonitors, cronMonitor)
	}

	return buildedMonitors, nil
}
