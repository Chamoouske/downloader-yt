package domain

import "os"

type Downloader interface {
	Download(video Video, progress ProgressBar) error
	Finalize(file os.File) error
	Cancel(file os.File) error
}
