package domain

type Downloader interface {
    Download(video Video, progress ProgressBar) error
}

