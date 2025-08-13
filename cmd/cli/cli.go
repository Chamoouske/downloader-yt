package main

import (
	"flag"
	"fmt"
	"os"

	dependencyinjections "downloader/internal/infra/dependency_injections"
	"downloader/internal/infra/notifyer/server"
	"downloader/internal/infra/progress"
	"downloader/internal/infra/youtube"
	"downloader/internal/usecase"
	"downloader/pkg/config"
	"net/http"
)

var url = flag.String("v", "", "Video must not be null")

func main() {
	flag.Parse()
	if *url == "" {
		flag.Usage()
		fmt.Println("Usage: downloader -v \"<url>\"")
		os.Exit(1)
	}

	cfg, err := config.NewConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	ytClient := &youtube.DefaultYoutubeClient{}
	fs := &youtube.DefaultOsFs{}
	progressBarClient := &progress.DefaultProgressBarClient{}

	// For notifyer, we'll use the server notifyer for now
	// In a real application, this might be chosen based on config
	serverNotifyer := server.NewServerNotifyer(cfg.URLWebhook, &http.Client{}, func(data []byte) server.Buffer { return server.NewDefaultBuffer(data) })
	db := dependencyinjections.NewVideoDatabase()

	downloader := youtube.NewKkdaiDownloader(serverNotifyer, db, ytClient, fs, cfg)
	progressBar := progress.NewTerminalProgressBar(progressBarClient)

	useCase := usecase.DownloadVideoUseCase{Downloader: downloader}
	err = useCase.Execute(usecase.Solicitation{URL: *url}, progressBar)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Download complete.")
	}
}
