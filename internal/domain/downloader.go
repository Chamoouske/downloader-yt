package domain

import "os"

type Downloader interface {
	Download(video Video, progress ProgressBar) error
	Finalize(notification Notification) error
	Cancel(file os.File) error
}
