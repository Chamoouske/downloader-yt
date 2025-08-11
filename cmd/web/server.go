package main

import (
	webserver "downloader/internal/infra/web_server"
	"downloader/internal/infra/youtube"
	logger "downloader/pkg/log"
	"flag"
	"os"
)

var log = logger.GetLogger("server")
var port = flag.Int("p", 0, "Port must not be null")

func main() {
	flag.Parse()

	if *port == 0 {
		flag.Usage()
		log.Error("Usage: downloader -p <port>")
		os.Exit(1)
	}

	svr := webserver.NewWebServer(youtube.NewKkdaiDownloader(nil))

	svr.Start(*port)
}
