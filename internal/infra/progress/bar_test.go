package progress

import (
	"testing"

	progressbar "github.com/schollz/progressbar/v3"
)

// MockProgressBarClient is a mock implementation of the ProgressBarClient interface
type MockProgressBarClient struct {
	NewOptions64Func func(total int64, options ...progressbar.Option) *progressbar.ProgressBar
	Set64Func        func(value int64)
	FinishFunc       func()
}

func (m *MockProgressBarClient) NewOptions64(total int64, options ...progressbar.Option) *progressbar.ProgressBar {
	if m.NewOptions64Func != nil {
		return m.NewOptions64Func(total, options...)
	}
	return nil
}

func (m *MockProgressBarClient) Set64(value int64) {
	if m.Set64Func != nil {
		m.Set64Func(value)
	}
}

func (m *MockProgressBarClient) Finish() {
	if m.FinishFunc != nil {
		m.FinishFunc()
	}
}

func TestTerminalProgressBar_Start(t *testing.T) {
	// Arrange
	mockBarClient := &MockProgressBarClient{
		NewOptions64Func: func(total int64, options ...progressbar.Option) *progressbar.ProgressBar {
			// Assert that NewOptions64 is called with the correct total
			if total != 100 {
				t.Errorf("expected total 100, got %d", total)
			}
			return nil // Return nil as we don't need a real progressbar.ProgressBar for this test
		},
	}
	progressBar := NewTerminalProgressBar(mockBarClient)

	// Act
	progressBar.Start(100)

	// Assertions are done within NewOptions64Func
}

func TestTerminalProgressBar_Update(t *testing.T) {
	// Arrange
	mockBarClient := &MockProgressBarClient{
		Set64Func: func(value int64) {
			// Assert that Set64 is called with the correct value
			if value != 50 {
				t.Errorf("expected value 50, got %d", value)
			}
		},
	}
	progressBar := NewTerminalProgressBar(mockBarClient)

	// Act
	progressBar.Update(50)

	// Assertions are done within Set64Func
}

func TestTerminalProgressBar_Finish(t *testing.T) {
	// Arrange
	finished := false
	mockBarClient := &MockProgressBarClient{
		FinishFunc: func() {
			// Assert that Finish is called
			finished = true
		},
	}
	progressBar := NewTerminalProgressBar(mockBarClient)

	// Act
	progressBar.Finish()

	// Assert
	if !finished {
		t.Error("expected Finish to be called, but it wasn't")
	}
}
