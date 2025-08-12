package server

import (
	"downloader/internal/domain"
	logger "downloader/pkg/log"
	"fmt"
	"net/http"
)

var log = logger.GetLogger("server")

type ServerNotifyer struct {
	URL string
}

func NewServerNotifyer(url string) *ServerNotifyer {
	return &ServerNotifyer{URL: url}
}

func (s *ServerNotifyer) Notify(notification domain.Notification) error {
	log.Info(fmt.Sprintf("Sending notification: %s", notification.Message))
	httpClient := &http.Client{}
	req, err := http.NewRequest("GET", s.URL, nil)
	if err != nil {
		log.Error(fmt.Sprintf("Error creating request: %v", err))
	}
	req.Header.Set("Content-Type", "application/json")
	q := req.URL.Query()
	q.Add("message", notification.Message)
	req.URL.RawQuery = q.Encode()

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
