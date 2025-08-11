package main

import (
	"flag"
	"fmt"
	"os"

	termux "downloader/internal/infra/notifyer"
	"downloader/internal/infra/progress"
	"downloader/internal/infra/youtube"
	"downloader/internal/usecase"
)

var url = flag.String("v", "", "Video must not be null")

func main() {
	flag.Parse()
	if *url == "" {
		flag.Usage()
		fmt.Println("Usage: downloader -v <url>")
		os.Exit(1)
	}

	downloader := youtube.NewKkdaiDownloader(termux.NewTermuxNotifyer())
	progressBar := progress.NewTerminalProgressBar()

	useCase := usecase.DownloadVideoUseCase{Downloader: downloader}
	err := useCase.Execute(*url, progressBar)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Download complete.")
	}
}
