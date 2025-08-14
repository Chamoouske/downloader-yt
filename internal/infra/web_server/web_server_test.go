package webserver

import (
	"bytes"
	"downloader/internal/domain"
	"downloader/pkg/config"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
)

// MockDownloader is a mock implementation of the domain.Downloader interface
type MockDownloader struct {
	DownloadFunc func(video domain.Video, progress domain.ProgressBar) error
	FinalizeFunc func(notification domain.Notification) error
	CancelFunc   func(file os.File) error
}

func (m *MockDownloader) Download(video domain.Video, progress domain.ProgressBar) error {
	if m.DownloadFunc != nil {
		return m.DownloadFunc(video, progress)
	}
	return nil
}

func (m *MockDownloader) Finalize(notification domain.Notification) error {
	if m.FinalizeFunc != nil {
		return m.FinalizeFunc(notification)
	}
	return nil
}

func (m *MockDownloader) Cancel(file os.File) error {
	if m.CancelFunc != nil {
		return m.CancelFunc(file)
	}
	return nil
}

// MockDatabase is a mock implementation of the domain.Database interface
type MockDatabase[T any] struct {
	SaveFunc func(id string, v T) error
	GetFunc  func(id string) (T, error)
	RemoveFunc func(id string) error
}

func (m *MockDatabase[T]) Save(id string, v T) error {
	if m.SaveFunc != nil {
		return m.SaveFunc(id, v)
	}
	return nil
}

func (m *MockDatabase[T]) Get(id string) (T, error) {
	if m.GetFunc != nil {
		return m.GetFunc(id)
	}
	var zero T
	return zero, nil
}

func (m *MockDatabase[T]) Remove(id string) error {
	if m.RemoveFunc != nil {
		return m.RemoveFunc(id)
	}
	return nil
}

// MockProgressBar is a mock implementation of the domain.ProgressBar interface
type MockProgressBar struct {
	StartFunc  func(total int64)
	UpdateFunc func(current int64)
	FinishFunc func()
}

func (m *MockProgressBar) Start(total int64) {
	if m.StartFunc != nil {
		m.StartFunc(total)
	}
}

func (m *MockProgressBar) Update(current int64) {
	if m.UpdateFunc != nil {
		m.UpdateFunc(current)
	}
}

func (m *MockProgressBar) Finish() {
	if m.FinishFunc != nil {
		m.FinishFunc()
	}
}

func TestWebServer_addVideoNaFilaDeDownload(t *testing.T) {
	t.Run("should return bad request if URL is missing", func(t *testing.T) {
		// Arrange
		req := httptest.NewRequest("GET", "/video/download?requester=test", nil)
		rr := httptest.NewRecorder()
		ws := NewWebServer(nil, nil, nil, slog.Default())

		// Act
		ws.addVideoNaFilaDeDownload(rr, req)

		// Assert
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
		}
		if !bytes.Contains(rr.Body.Bytes(), []byte("URL parameter is required")) {
			t.Errorf("expected error message, got %s", rr.Body.String())
		}
	})

	t.Run("should return bad request if requester is missing", func(t *testing.T) {
		// Arrange
		req := httptest.NewRequest("GET", "/video/download?url=test_url", nil)
		rr := httptest.NewRecorder()
		ws := NewWebServer(nil, nil, nil, slog.Default())

		// Act
		ws.addVideoNaFilaDeDownload(rr, req)

		// Assert
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
		}
		if !bytes.Contains(rr.Body.Bytes(), []byte("requester parameter is required")) {
			t.Errorf("expected error message, got %s", rr.Body.String())
		}
	})

	t.Run("should start download and return success message", func(t *testing.T) {
		// Arrange
		mockDownloader := &MockDownloader{
			DownloadFunc: func(video domain.Video, progress domain.ProgressBar) error {
				return nil
			},
		}
		req := httptest.NewRequest("GET", "/video/download?url=test_url&requester=test_requester", nil)
		rr := httptest.NewRecorder()
		ws := NewWebServer(mockDownloader, nil, nil, slog.Default())

		// Act
		ws.addVideoNaFilaDeDownload(rr, req)

		// Assert
		if rr.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
		}
		if !bytes.Contains(rr.Body.Bytes(), []byte("Download iniciado")) {
			t.Errorf("expected success message, got %s", rr.Body.String())
		}
	})
}

func TestWebServer_download(t *testing.T) {
	t.Run("should return not found if video is not in db", func(t *testing.T) {
		// Arrange
		mockDb := &MockDatabase[domain.Video]{
			GetFunc: func(id string) (domain.Video, error) {
				return domain.Video{}, errors.New("not found")
			},
		}
		req := httptest.NewRequest("GET", "/video/test_id", nil)
		rr := httptest.NewRecorder()
		ws := NewWebServer(nil, mockDb, nil, slog.Default())

		router := mux.NewRouter()
		router.HandleFunc("/video/{id}", ws.download).Methods("GET")
		router.ServeHTTP(rr, req)

		// Assert
		if rr.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, rr.Code)
		}
	})

	t.Run("should return forbidden for invalid path", func(t *testing.T) {
		// Arrange
		mockDb := &MockDatabase[domain.Video]{
			GetFunc: func(id string) (domain.Video, error) {
				return domain.Video{Filename: "../evil.mp4"}, nil
			},
		}
		cfg := &config.Config{VideoDir: "./videos"}
		req := httptest.NewRequest("GET", "/video/test_id", nil)
		rr := httptest.NewRecorder()
		ws := NewWebServer(nil, mockDb, cfg, slog.Default())

		router := mux.NewRouter()
		router.HandleFunc("/video/{id}", ws.download).Methods("GET")
		router.ServeHTTP(rr, req)

		// Assert
		if rr.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, rr.Code)
		}
	})

	t.Run("should return internal server error if file cannot be opened", func(t *testing.T) {
		// Arrange
		mockDb := &MockDatabase[domain.Video]{
			GetFunc: func(id string) (domain.Video, error) {
				return domain.Video{Filename: "test_video.mp4"}, nil
			},
		}
		cfg := &config.Config{VideoDir: "./videos"}
		req := httptest.NewRequest("GET", "/video/test_id", nil)
		rr := httptest.NewRecorder()
		ws := NewWebServer(nil, mockDb, cfg, slog.Default())

		router := mux.NewRouter()
		router.HandleFunc("/video/{id}", ws.download).Methods("GET")
		router.ServeHTTP(rr, req)

		// Assert
		if rr.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, rr.Code)
		}
	})
}
