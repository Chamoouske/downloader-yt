package youtube

import (
	"downloader/internal/domain"
	"fmt"
	"io"
	"os"
	"strings"

	yt "github.com/kkdai/youtube/v2"
)

type KkdaiDownloader struct {
  file	   os.File
  notifyer domain.Notifyer
}

func NewKkdaiDownloader(notifyer domain.Notifyer) *KkdaiDownloader {
    return &KkdaiDownloader{notifyer: notifyer}
}

func (d *KkdaiDownloader) Download(video domain.Video, progress domain.ProgressBar) error {
    client := yt.Client{}

    ytVideo, err := client.GetVideo(video.URL)
    if err != nil {
        return fmt.Errorf("error fetching video info: %w", err)
    }

    format := &ytVideo.Formats[0]

    stream, size, err := client.GetStream(ytVideo, format)
    if err != nil {
        return fmt.Errorf("error getting video stream: %w", err)
    }

    fileName := ytVideo.Title + "." + strings.Split(format.MimeType, ";")[0][6:]

    outFile, err := os.Create(fileName)
    if err != nil {
        return fmt.Errorf("error creating file: %w", err)
    }
    d.file = *outFile
    defer outFile.Close()

    progress.Start(size)

    proxyReader := io.TeeReader(stream, &progressWriter{progress: progress})

    _, err = io.Copy(outFile, proxyReader)
    if err != nil {
        return fmt.Errorf("error saving video: %w", err)
    }

    progress.Finish()
    d.Finalize()
    return nil
}

func (d *KkdaiDownloader) Finalize() error {
    msg := d.file.Name() + " was downloaded"

    notification := &domain.Notification {
        Title:   "Download Finalized",
        Message: msg,
    }

    if err := d.notifyer.Notify(*notification); err != nil {
	return fmt.Errorf("erro to notify user: %w", err)
    }

    return nil
}


func (d *KkdaiDownloader) Cancel(video domain.Video) error {
  if err := d.file.Close(); err != nil {
  	return fmt.Errorf("error closing file: %w", err)
  }
  if err := os.Remove(d.file.Name()); err != nil {
  	return fmt.Errorf("error deleting file: %w", err)
  }

  return nil
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
