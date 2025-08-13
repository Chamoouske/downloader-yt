package server

import (
	"bytes"
	"downloader/internal/domain"
	logger "downloader/pkg/log"
	"encoding/json"
	"fmt"
	"net/http"
)

var log = logger.GetLogger("server")

type ServerNotifyer struct {
	URL string
}

type VideoResponse struct {
	URL string `json:"url"`
	To  string `json:"to"`
}

func NewServerNotifyer(url string) *ServerNotifyer {
	return &ServerNotifyer{URL: url}
}

func (s *ServerNotifyer) Notify(notification domain.Notification) error {
	log.Info(fmt.Sprintf("Sending notification: %s", notification.Title))
	httpClient := &http.Client{}

	obj, err := json.Marshal(VideoResponse{
		URL: fmt.Sprintf("https://downloader.ajaxlima.dev.br/video/%s", notification.Message),
		To:  notification.To,
	})
	if err != nil {
		msgError := fmt.Sprintf("Error creating json obj: %v", err)
		log.Error(msgError)
		return fmt.Errorf("%s", msgError)
	}

	req, err := http.NewRequest("GET", s.URL, bytes.NewBuffer(obj))
	if err != nil {
		msgError := fmt.Sprintf("Error creating request: %v", err)
		log.Error(msgError)
		return fmt.Errorf("%s", msgError)
	}
	req.Header.Set("Content-Type", "application/json")

	log.Info(fmt.Sprintf("Request from: %s", s.URL))
	resp, err := httpClient.Do(req)
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
