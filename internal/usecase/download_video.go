package usecase

import (
	"downloader/internal/domain"
)

type DownloadVideoUseCase struct {
	Downloader domain.Downloader
}

func (uc *DownloadVideoUseCase) Execute(url string, progress domain.ProgressBar) error {
	return uc.Downloader.Download(domain.Video{URL: url}, progress)
}
