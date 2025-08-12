package domain

import "os"

type Downloader interface {
	Download(video Video, progress ProgressBar) error
	Finalize(msg string) error
	Cancel(file os.File) error
}
