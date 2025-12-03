package cmd

import (
	"flag"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/andrewsapw/avalio/app"
	"github.com/andrewsapw/avalio/monitors"
	"github.com/andrewsapw/avalio/notificators"
	"github.com/andrewsapw/avalio/resources"
)

func StartAvalio() {
	configPath := flag.String("config", "", "config path")
	flag.Parse()

	config, err := app.ParseConfig(*configPath)
	if err != nil {
		log.Default().Fatal(err.Error())
		os.Exit(1)
	}

	var logLevel slog.Level
	switch strings.ToLower(config.LogLevel) {
	case "info":
		logLevel = slog.LevelInfo
	case "debug":
		logLevel = slog.LevelDebug
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})
	logger := slog.New(handler)

	logger.Info("Loading configuration file", "config_path", *configPath)

	resources, err := resources.BuildResources(&config.Resources, logger)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	notificators, err := notificators.BuildNotificators(&config.Notificators, logger)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	monitors, err := monitors.BuildMonitors(&config.Monitors, logger)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	application := app.NewApplication(resources, notificators, monitors, logger)

	err = application.Run()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
