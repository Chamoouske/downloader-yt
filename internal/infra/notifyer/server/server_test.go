package server

import (
	"bytes"
	"downloader/internal/domain"
	"errors"
	"io"
	"net/http"
	"testing"
)

// MockHTTPClient is a mock implementation of the HTTPClient interface
type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if m.DoFunc != nil {
		return m.DoFunc(req)
	}
	return nil, nil
}

// MockBuffer is a mock implementation of the Buffer interface
type MockBuffer struct {
	BytesFunc    func() []byte
	StringFunc   func() string
	WriteFunc    func(p []byte) (n int, err error)
	WriteStringFunc func(s string) (n int, err error)
	ReadFunc     func(p []byte) (n int, err error)
}

func (m *MockBuffer) Bytes() []byte {
	if m.BytesFunc != nil {
		return m.BytesFunc()
	}
	return nil
}

func (m *MockBuffer) String() string {
	if m.StringFunc != nil {
		return m.StringFunc()
	}
	return ""
}

func (m *MockBuffer) Write(p []byte) (n int, err error) {
	if m.WriteFunc != nil {
		return m.WriteFunc(p)
	}
	return 0, nil
}

func (m *MockBuffer) WriteString(s string) (n int, err error) {
	if m.WriteStringFunc != nil {
		return m.WriteStringFunc(s)
	}
	return 0, nil
}

func (m *MockBuffer) Read(p []byte) (n int, err error) {
	if m.ReadFunc != nil {
		return m.ReadFunc(p)
	}
	return 0, io.EOF
}

func TestServerNotifyer_Notify(t *testing.T) {
	t.Run("should send notification successfully", func(t *testing.T) {
		// Arrange
		mockHttpClient := &MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewBufferString(""))}, nil
			},
		}
		mockBuffer := &MockBuffer{
			ReadFunc: func(p []byte) (n int, err error) {
				return len(p), io.EOF
			},
		}
		notifyer := NewServerNotifyer("http://test.com", mockHttpClient, func(data []byte) Buffer { return mockBuffer })
		notification := domain.Notification{Title: "Test", Message: "123", To: "user"}

		// Act
		err := notifyer.Notify(notification)

		// Assert
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("should return error when http client returns error", func(t *testing.T) {
		// Arrange
		expectedErr := errors.New("http error")
		mockHttpClient := &MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return nil, expectedErr
			},
		}
		mockBuffer := &MockBuffer{
			ReadFunc: func(p []byte) (n int, err error) {
				return len(p), io.EOF
			},
		}
		notifyer := NewServerNotifyer("http://test.com", mockHttpClient, func(data []byte) Buffer { return mockBuffer })
		notification := domain.Notification{Title: "Test", Message: "123", To: "user"}

		// Act
		err := notifyer.Notify(notification)

		// Assert
		if err == nil {
			t.Error("expected an error, but got nil")
		}
		if err != expectedErr {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
	})

	t.Run("should return error when http status is not OK", func(t *testing.T) {
		// Arrange
		mockHttpClient := &MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{StatusCode: http.StatusBadRequest, Body: io.NopCloser(bytes.NewBufferString(""))}, nil
			},
		}
		mockBuffer := &MockBuffer{
			ReadFunc: func(p []byte) (n int, err error) {
				return len(p), io.EOF
			},
		}
		notifyer := NewServerNotifyer("http://test.com", mockHttpClient, func(data []byte) Buffer { return mockBuffer })
		notification := domain.Notification{Title: "Test", Message: "123", To: "user"}

		// Act
		err := notifyer.Notify(notification)

		// Assert
		if err == nil {
			t.Error("expected an error, but got nil")
		}
		if err.Error() != "received non-200 response: 400" {
			t.Errorf("expected error 'received non-200 response: 400', got %v", err)
		}
	})
}
