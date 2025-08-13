package youtube

import (
	"downloader/internal/domain"
	"downloader/pkg/config"
	logger "downloader/pkg/log"
	"downloader/pkg/utils"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/google/uuid"
	yt "github.com/kkdai/youtube/v2"
)

var log *slog.Logger

func init() {
	cfg, _ := config.NewConfig()
	log = logger.NewLogger(cfg).With("component", "youtube")
}

// YoutubeClient interface to abstract youtube client operations
type YoutubeClient interface {
	GetVideo(videoID string) (*yt.Video, error)
	GetStream(video *yt.Video, format *yt.Format) (io.ReadCloser, int64, error)
}

// OsFs interface to abstract os file system operations
type OsFs interface {
	OpenFile(name string, flag int, perm os.FileMode) (*os.File, error)
	Remove(name string) error
}

// KkdaiDownloader struct with injected dependencies
type KkdaiDownloader struct {
	notifyer domain.Notifyer
	db       domain.Database[domain.Video]
	ytClient YoutubeClient
	fs       OsFs
	cfg      *config.Config
}

func NewKkdaiDownloader(notifyer domain.Notifyer, db domain.Database[domain.Video], ytClient YoutubeClient, fs OsFs, cfg *config.Config) *KkdaiDownloader {
	return &KkdaiDownloader{notifyer: notifyer, db: db, ytClient: ytClient, fs: fs, cfg: cfg}
}

func (d *KkdaiDownloader) Download(video domain.Video, progress domain.ProgressBar) error {
	ytVideo, err := d.ytClient.GetVideo(video.URL)
	if err != nil {
		return fmt.Errorf("error fetching video info: %w", err)
	}

	format := &ytVideo.Formats[0]

	stream, size, err := d.ytClient.GetStream(ytVideo, format)
	if err != nil {
		return fmt.Errorf("error getting video stream: %w", err)
	}

	id := uuid.NewString()
	fileName := utils.SanitizeFilename(id + "." + strings.Split(format.MimeType, ";")[0][6:])

	outFile, err := d.fs.OpenFile(filepath.Join(d.cfg.VideoDir, fileName), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o777)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	d.configCancelSignal(*outFile)
	defer outFile.Close()

	log.Info(fmt.Sprintf("Download do vídeo %s iniciado!", ytVideo.Title))
	progress.Start(size)

	proxyReader := io.TeeReader(stream, &progressWriter{progress: progress})

	_, err = io.Copy(outFile, proxyReader)
	if err != nil {
		return fmt.Errorf("error saving video: %w", err)
	}

	progress.Finish()
	video.Filename = utils.SanitizeFilename(ytVideo.Title)
	if d.db != nil {
		d.db.Save(id, video)
	}
	if d.notifyer != nil {
		d.Finalize(domain.Notification{
			Title:   ytVideo.Title,
			Message: id,
			To:      video.Requester,
		})
	}
	return nil
}

func (d *KkdaiDownloader) Finalize(notification domain.Notification) error {
	if err := d.notifyer.Notify(notification); err != nil {
		return fmt.Errorf("erro to notify user: %w", err)
	}

	return nil
}

func (d *KkdaiDownloader) Cancel(file os.File) error {
	if err := file.Close(); err != nil {
		return fmt.Errorf("error closing file: %w", err)
	}
	if err := d.fs.Remove(file.Name()); err != nil {
		return fmt.Errorf("error deleting file: %w", err)
	}

	return nil
}

func (d *KkdaiDownloader) configCancelSignal(file os.File) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		d.Cancel(file)
	}()
}

// DefaultYoutubeClient implements YoutubeClient using kkdai/youtube/v2
type DefaultYoutubeClient struct{}

func (c *DefaultYoutubeClient) GetVideo(videoID string) (*yt.Video, error) {
	client := yt.Client{}
	return client.GetVideo(videoID)
}

func (c *DefaultYoutubeClient) GetStream(video *yt.Video, format *yt.Format) (io.ReadCloser, int64, error) {
	client := yt.Client{}
	return client.GetStream(video, format)
}

// DefaultOsFs implements OsFs using os package
type DefaultOsFs struct{}

func (fs *DefaultOsFs) OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(name, flag, perm)
}

func (fs *DefaultOsFs) Remove(name string) error {
	return os.Remove(name)
}

type progressWriter struct {
	total    int64
	current  int64
	progress domain.ProgressBar
}

func (pw *progressWriter) Write(p []byte) (int, error) {
	n := len(p)
	pw.current += int64(n)
	pw.progress.Update(pw.current)
	return n, nil
}
