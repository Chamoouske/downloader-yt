package main

import (
	dependencyinjections "downloader/internal/infra/dependency_injections"
	serverNotifyer "downloader/internal/infra/notifyer/server"
	webserver "downloader/internal/infra/web_server"
	"downloader/internal/infra/youtube"
	"downloader/pkg/config"
	logger "downloader/pkg/log"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
)

var port = flag.Int("p", 0, "Port must not be null")

func main() {
	flag.Parse()

	cfg, err := config.NewConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	appLogger := logger.NewLogger(cfg)

	if *port == 0 && os.Getenv("PORT") == "" {
		flag.Usage()
		appLogger.Error("Usage: downloader -p <port>")
		os.Exit(1)
	}

	// Dependencies for KkdaiDownloader
	ytClient := &youtube.DefaultYoutubeClient{}
	fs := &youtube.DefaultOsFs{}

	// Dependencies for Notifyer
	serverHttpClient := &http.Client{}
	serverBuffer := func(data []byte) serverNotifyer.Buffer { return serverNotifyer.NewDefaultBuffer(data) }
	notifyer := serverNotifyer.NewServerNotifyer(cfg.URLWebhook, serverHttpClient, serverBuffer)

	// Database
	db := dependencyinjections.NewVideoDatabase()

	// Downloader
	downloader := youtube.NewKkdaiDownloader(notifyer, db, ytClient, fs, cfg)

	// WebServer
	svr := webserver.NewWebServer(downloader, db, cfg, appLogger)

	svr.Start(getPort(appLogger))
}

func getPort(appLogger *slog.Logger) int {
	if *port != 0 {
		return *port
	}

	portEnv := os.Getenv("PORT")
	if portEnv == "" {
		appLogger.Error("Port is not set in environment variables")
		os.Exit(1)
	}

	portValue, err := strconv.Atoi(portEnv)
	if err != nil {
		appLogger.Error("Invalid PORT environment variable", "error", err)
		os.Exit(1)
	}

	return portValue
}
