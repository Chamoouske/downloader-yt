package usecase

import (
	"downloader/internal/domain"
)

type DownloadVideoUseCase struct {
	Downloader domain.Downloader
}

type Solicitation struct {
	URL       string
	Requester string
}

func (uc *DownloadVideoUseCase) Execute(sol Solicitation, progress domain.ProgressBar) error {
	return uc.Downloader.Download(domain.Video{URL: sol.URL, Requester: sol.Requester}, progress)
}
