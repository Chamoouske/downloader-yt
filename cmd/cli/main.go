package main

import (
	"flag"
	"fmt"
	"os"
    "os/signal"
    "syscall"

	"downloader/internal/domain"
	"downloader/internal/infra/progress"
	"downloader/internal/infra/termux"
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

    video := domain.Video{URL: *url}

    downloader := youtube.NewKkdaiDownloader(termux.NewTermuxNotifyer())
    progressBar := progress.NewTerminalProgressBar()

    useCase := usecase.DownloadVideoUseCase{Downloader: downloader}
    err := useCase.Execute(video, progressBar)
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

    go func() {
        <-sigChan
        useCase.Downloader.Cancel(video)
    }()

    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Println("Download complete.")
    }
}
