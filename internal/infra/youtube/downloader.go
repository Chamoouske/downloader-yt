package youtube

import (
	"downloader/internal/domain"
	"downloader/pkg/config"
	logger "downloader/pkg/log"
	"downloader/pkg/utils"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/google/uuid"
	yt "github.com/kkdai/youtube/v2"
)

var log = logger.GetLogger("youtube")

type KkdaiDownloader struct {
	notifyer domain.Notifyer
	db       domain.Database[domain.Video]
}

func NewKkdaiDownloader(notifyer domain.Notifyer, db domain.Database[domain.Video]) *KkdaiDownloader {
	return &KkdaiDownloader{notifyer: notifyer, db: db}
}

func (d *KkdaiDownloader) Download(video domain.Video, progress domain.ProgressBar) error {
	client := yt.Client{}
	cfg := config.GetConfig()

	ytVideo, err := client.GetVideo(video.URL)
	if err != nil {
		return fmt.Errorf("error fetching video info: %w", err)
	}

	format := &ytVideo.Formats[0]

	stream, size, err := client.GetStream(ytVideo, format)
	if err != nil {
		return fmt.Errorf("error getting video stream: %w", err)
	}

	id := uuid.NewString()
	fileName := utils.SanitizeFilename(id + "." + strings.Split(format.MimeType, ";")[0][6:])

	outFile, err := os.OpenFile(filepath.Join(cfg.VideoDir, fileName), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o777)
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
	if err := os.Remove(file.Name()); err != nil {
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
