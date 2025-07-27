package domain

type Video struct {
    URL string
}

type Downloader interface {
    Download(video Video, progress ProgressBar) error
}

type ProgressBar interface {
    Start(total int64)
    Update(current int64)
    Finish()
}
