package domain

type Server interface {
	Start(port int)
	Stop()
	Download(video Video) error
}
