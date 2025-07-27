package usecase

import (
  "downloader/internal/domain"
  "os"
  "os/signal"
  "syscall"
  "fmt"
)

type DownloadVideoUseCase struct {
    Downloader domain.Downloader
}

func (uc *DownloadVideoUseCase) Execute(video domain.Video, progress domain.ProgressBar) error {
    uc.Cancel(video)
    return uc.Downloader.Download(video, progress)
}

func (uc *DownloadVideoUseCase) Cancel(video domain.Video) {
  sigChan := make(chan os.Signal, 1)
  signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

  go func() {
      <-sigChan
      fmt.Println("\nDownload cancelado. Removendo arquivo parcial...")
      uc.Downloader.Cancel(video)
    }()
}
