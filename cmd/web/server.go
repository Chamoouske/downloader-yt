package main

import (
	dependencyinjections "downloader/internal/infra/dependency_injections"
	"downloader/internal/infra/notifyer/server"
	webserver "downloader/internal/infra/web_server"
	"downloader/internal/infra/youtube"
	"downloader/pkg/config"
	logger "downloader/pkg/log"
	"flag"
	"os"
	"strconv"
)

var log = logger.GetLogger("server")
var port = flag.Int("p", 0, "Port must not be null")

func main() {
	flag.Parse()
	cfg := config.GetConfig()

	if *port == 0 && os.Getenv("PORT") == "" {
		flag.Usage()
		log.Error("Usage: downloader -p <port>")
		os.Exit(1)
	}

	notifyer := server.NewServerNotifyer(cfg.URLWebhook)
	db := *dependencyinjections.GetVideoDatabase()

	svr := webserver.NewWebServer(youtube.NewKkdaiDownloader(notifyer, db), db)

	svr.Start(getPort())
}

func getPort() int {
	if *port != 0 {
		return *port
	}

	if os.Getenv("PORT") == "" {
		log.Error("Port is not set in environment variables")
		os.Exit(1)
	}

	portValue, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Error("Invalid PORT environment variable", "error", err)
		os.Exit(1)
	}

	return portValue
}
