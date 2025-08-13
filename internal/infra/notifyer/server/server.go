package server

import (
	"bytes"
	"downloader/internal/domain"
	"downloader/pkg/config"
	logger "downloader/pkg/log"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

var log *slog.Logger

func init() {
	cfg, _ := config.NewConfig()
	log = logger.NewLogger(cfg).With("component", "server")
}

// HTTPClient interface to abstract http client operations
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Buffer interface to abstract bytes.Buffer operations
type Buffer interface {
	io.Reader
	Bytes() []byte
	String() string
	Write(p []byte) (n int, err error)
	WriteString(s string) (n int, err error)
}

// ServerNotifyer struct with injected dependencies
type ServerNotifyer struct {
	URL        string
	httpClient HTTPClient
	buffer     func([]byte) Buffer
}

type VideoResponse struct {
	URL string `json:"url"`
	To  string `json:"to"`
}

func NewServerNotifyer(url string, httpClient HTTPClient, buffer func([]byte) Buffer) *ServerNotifyer {
	return &ServerNotifyer{URL: url, httpClient: httpClient, buffer: buffer}
}

func (s *ServerNotifyer) Notify(notification domain.Notification) error {
	log.Info(fmt.Sprintf("Sending notification: %s", notification.Title))

	obj, err := json.Marshal(VideoResponse{
		URL: fmt.Sprintf("https://downloader.ajaxlima.dev.br/video/%s", notification.Message),
		To:  notification.To,
	})
	if err != nil {
		msgError := fmt.Sprintf("Error creating json obj: %v", err)
		log.Error(msgError)
		return fmt.Errorf("%s", msgError)
	}

	req, err := http.NewRequest("GET", s.URL, s.buffer(obj))
	if err != nil {
		msgError := fmt.Sprintf("Error creating request: %v", err)
		log.Error(msgError)
		return fmt.Errorf("%s", msgError)
	}
	req.Header.Set("Content-Type", "application/json")

	log.Info(fmt.Sprintf("Request from: %s", s.URL))
	resp, err := s.httpClient.Do(req)
	if err != nil {
		log.Error(fmt.Sprintf("Error sending request: %v", err))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Error(fmt.Sprintf("Received non-200 response: %d", resp.StatusCode))
		return fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	log.Info("Notification sent successfully")
	return nil
}

// DefaultHTTPClient implements HTTPClient using net/http
type DefaultHTTPClient struct{}

func (c *DefaultHTTPClient) Do(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	return client.Do(req)
}

// DefaultBuffer implements Buffer using bytes.Buffer
type DefaultBuffer struct {
	*bytes.Buffer
}

func NewDefaultBuffer(data []byte) Buffer {
	return &DefaultBuffer{bytes.NewBuffer(data)}
}
