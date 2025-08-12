package main

import (
	"downloader/internal/infra/notifyer/server"
	webserver "downloader/internal/infra/web_server"
	"downloader/internal/infra/youtube"
	logger "downloader/pkg/log"
	"flag"
	"os"
	"strconv"
)

var log = logger.GetLogger("server")
var port = flag.Int("p", 0, "Port must not be null")

func main() {
	flag.Parse()

	if *port == 0 && os.Getenv("PORT") == "" {
		flag.Usage()
		log.Error("Usage: downloader -p <port>")
		os.Exit(1)
	}

	notifyer := server.NewServerNotifyer("http://localhost:8080/notify")

	svr := webserver.NewWebServer(youtube.NewKkdaiDownloader(notifyer))
	if *port != 0 {
		svr.Start(*port)
	} else {
		portEnv := os.Getenv("PORT")
		if portEnv == "" {
			log.Error("Port is not set in environment variables")
			os.Exit(1)
		}

		portValue, err := strconv.Atoi(portEnv)
		if err != nil {
			log.Error("Invalid PORT environment variable", "error", err)
			os.Exit(1)
		}

		svr.Start(portValue)
	}
}
