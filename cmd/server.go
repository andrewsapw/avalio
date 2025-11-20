package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"
	"slices"

	"github.com/BurntSushi/toml"
	"github.com/andrewsapw/avalio/internal/app"
	"github.com/andrewsapw/avalio/internal/monitors"
	"github.com/andrewsapw/avalio/internal/notificators"
	"github.com/andrewsapw/avalio/internal/resources"
)

func StartServer() {
	configPath := flag.String("config", "", "config path")

	flag.Parse()

	fmt.Println("configPath:", *configPath)

	logger := log.New(os.Stdout, "[app] ", log.LstdFlags|log.Lshortfile)

	config, err := parseConfig(*configPath)
	if err != nil {
		logger.Fatal(err)
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

func buildResources(config *app.Config, logger *log.Logger) ([]resources.Resource, error) {
	var buildedResources []resources.Resource
	for _, httpResourceConfig := range config.Resources.Http {
		httpResource := resources.NewHTTPResource(httpResourceConfig, logger)
		logger.Printf("Builded resource %s", httpResource.GetName())
		buildedResources = append(buildedResources, httpResource)
	}

	return buildedResources, nil
}

func buildNotificators(config *app.Config, logger *log.Logger) []notificators.Notificator {
	var buildedNotificators []notificators.Notificator
	notificatorsNames := []string{}

	for _, consoleNotificatorConfig := range config.Notificators.Console {
		consoleNotificator := notificators.NewConsoleNotificator(consoleNotificatorConfig, logger)
		if slices.Contains(notificatorsNames, consoleNotificator.GetName()) {
			logger.Fatalf("Duplicated notificators names: %s", consoleNotificator.GetName())
		}

		logger.Printf("Builded notificator %s", consoleNotificator.GetName())
		buildedNotificators = append(buildedNotificators, consoleNotificator)

		notificatorsNames = append(notificatorsNames, consoleNotificator.GetName())
	}

	for _, telegramNotificatorConfig := range config.Notificators.Telegram {
		if err := telegramNotificatorConfig.Validate(); err != nil {
			logger.Fatal(err.Error())
		}
		telegramNotificator := notificators.NewTelegramNotificator(telegramNotificatorConfig, logger)
		if slices.Contains(notificatorsNames, telegramNotificator.GetName()) {
			logger.Fatalf("Duplicated notificators names: %s", telegramNotificator.GetName())
		}
		logger.Printf("Builded notificator %s", telegramNotificator.GetName())
		buildedNotificators = append(buildedNotificators, telegramNotificator)
	}

	return buildedNotificators
}

func buildMonitors(config *app.Config, logger *log.Logger) []monitors.Monitor {
	var buildedMonitors []monitors.Monitor
	for _, cronMonitorConfig := range config.Monitors.Cron {
		cronMonitor := monitors.NewCronMonitor(cronMonitorConfig, logger)
		logger.Printf("Builded monitor %s", cronMonitorConfig.Name)
		buildedMonitors = append(buildedMonitors, cronMonitor)
	}

	return buildedMonitors
}
