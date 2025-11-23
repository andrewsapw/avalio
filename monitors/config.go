package monitors

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
