package cmd

import (
	"flag"
	"log/slog"
	"os"
	"slices"

	"github.com/BurntSushi/toml"
	"github.com/andrewsapw/avalio/app"
	"github.com/andrewsapw/avalio/monitors"
	"github.com/andrewsapw/avalio/notificators"
	"github.com/andrewsapw/avalio/resources"
)

func StartAvalio() {
	configPath := flag.String("config", "", "config path")

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	logger := slog.New(handler)

	flag.Parse()

	logger.Info("loading configuration file", "config_path", *configPath)

	config, err := parseConfig(*configPath)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	resources, _ := buildResources(config, logger)
	notificators := buildNotificators(config, logger)
	monitors := buildMonitors(config, logger)

	application := app.NewApplication(resources, notificators, monitors, logger)
	application.Run()
}

func parseConfig(configPath string) (*app.Config, error) {
	var config app.Config

	_, err := toml.DecodeFile(configPath, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func buildResources(config *app.Config, logger *slog.Logger) ([]resources.Resource, error) {
	var buildedResources []resources.Resource
	for _, httpResourceConfig := range config.Resources.Http {
		httpResource := resources.NewHTTPResource(httpResourceConfig, logger)
		logger.Info("Builded resource", "resource_name", httpResource.GetName())
		buildedResources = append(buildedResources, httpResource)
	}

	return buildedResources, nil
}

func buildNotificators(config *app.Config, logger *slog.Logger) []notificators.Notificator {
	var buildedNotificators []notificators.Notificator
	notificatorsNames := []string{}

	for _, consoleNotificatorConfig := range config.Notificators.Console {
		consoleNotificator := notificators.NewConsoleNotificator(consoleNotificatorConfig, logger)
		if slices.Contains(notificatorsNames, consoleNotificator.GetName()) {
			logger.Error("Duplicated notificators names", "duplicated_names", consoleNotificator.GetName())
			os.Exit(1)
		}

		logger.Info("Builded notificator", "notificatorName", consoleNotificator.GetName())
		buildedNotificators = append(buildedNotificators, consoleNotificator)

		notificatorsNames = append(notificatorsNames, consoleNotificator.GetName())
	}

	for _, telegramNotificatorConfig := range config.Notificators.Telegram {
		if err := telegramNotificatorConfig.Validate(); err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
		telegramNotificator := notificators.NewTelegramNotificator(telegramNotificatorConfig, logger)
		if slices.Contains(notificatorsNames, telegramNotificator.GetName()) {
			logger.Error("Duplicated notificators names", "duplicated_names", telegramNotificator.GetName())
			os.Exit(1)
		}
		logger.Info("Builded notificator", "notificator_name", telegramNotificator.GetName())
		buildedNotificators = append(buildedNotificators, telegramNotificator)
	}

	return buildedNotificators
}

func buildMonitors(config *app.Config, logger *slog.Logger) []monitors.Monitor {
	var buildedMonitors []monitors.Monitor
	for _, cronMonitorConfig := range config.Monitors.Cron {
		cronMonitor := monitors.NewCronMonitor(cronMonitorConfig, logger)
		logger.Info("Builded monitor", "monitor_name", cronMonitorConfig.Name)
		buildedMonitors = append(buildedMonitors, cronMonitor)
	}

	return buildedMonitors
}
