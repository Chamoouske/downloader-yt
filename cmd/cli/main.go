package main

import (
	"fmt"
	"os"

	"downloader/internal/domain"
	"downloader/internal/infra/progress"
	"downloader/internal/infra/termux"
	"downloader/internal/infra/youtube"
	"downloader/internal/usecase"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: downloader <url>")
        os.Exit(1)
    }

    url := os.Args[1]
    video := domain.Video{URL: url}

    downloader := youtube.NewKkdaiDownloader(termux.NewTermuxNotifyer())
    progressBar := progress.NewTerminalProgressBar()

    useCase := usecase.DownloadVideoUseCase{Downloader: downloader}
    err := useCase.Execute(video, progressBar)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Println("Download complete.")
    }
}
