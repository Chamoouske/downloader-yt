package termux

import (
	"downloader/internal/domain"
	"errors"
	"testing"
)

// MockCmd is a mock implementation of the Cmd interface
type MockCmd struct {
	RunFunc func() error
}

func (m *MockCmd) Run() error {
	if m.RunFunc != nil {
		return m.RunFunc()
	}
	return nil
}

// MockCommander is a mock implementation of the Commander interface
type MockCommander struct {
	CommandFunc func(name string, arg ...string) Cmd
}

func (m *MockCommander) Command(name string, arg ...string) Cmd {
	if m.CommandFunc != nil {
		return m.CommandFunc(name, arg...)
	}
	return &MockCmd{}
}

func TestTermuxNotifyer_Notify(t *testing.T) {
	t.Run("should send notification successfully", func(t *testing.T) {
		// Arrange
		mockCmd := &MockCmd{
			RunFunc: func() error {
				return nil
			},
		}
		mockCommander := &MockCommander{
			CommandFunc: func(name string, arg ...string) Cmd {
				return mockCmd
			},
		}
		notifyer := NewTermuxNotifyer(mockCommander)
		notification := domain.Notification{Title: "Test Title", Message: "Test Message"}

		// Act
		err := notifyer.Notify(notification)

		// Assert
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("should return error when command fails", func(t *testing.T) {
		// Arrange
		expectedErr := errors.New("command failed")
		mockCmd := &MockCmd{
			RunFunc: func() error {
				return expectedErr
			},
		}
		mockCommander := &MockCommander{
			CommandFunc: func(name string, arg ...string) Cmd {
				return mockCmd
			},
		}
		notifyer := NewTermuxNotifyer(mockCommander)
		notification := domain.Notification{Title: "Test Title", Message: "Test Message"}

		// Act
		err := notifyer.Notify(notification)

		// Assert
		if err == nil {
			t.Error("expected an error, but got nil")
		}
		if !errors.Is(err, expectedErr) {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
	})
}