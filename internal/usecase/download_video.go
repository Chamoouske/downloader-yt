package usecase

import "downloader/internal/domain"

type DownloadVideoUseCase struct {
    Downloader domain.Downloader
}

func (uc *DownloadVideoUseCase) Execute(video domain.Video, progress domain.ProgressBar) error {
    return uc.Downloader.Download(video, progress)
}
