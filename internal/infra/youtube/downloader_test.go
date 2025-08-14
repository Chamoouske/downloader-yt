package youtube

import (
	"downloader/internal/domain"
	"downloader/pkg/config"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kkdai/youtube/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockYoutubeClient mocks the YoutubeClient interface
type MockYoutubeClient struct {
	mock.Mock
}

func (m *MockYoutubeClient) GetVideo(videoID string) (*youtube.Video, error) {
	args := m.Called(videoID)
	return args.Get(0).(*youtube.Video), args.Error(1)
}

func (m *MockYoutubeClient) GetStream(video *youtube.Video, format *youtube.Format) (io.ReadCloser, int64, error) {
	args := m.Called(video, format)
	return args.Get(0).(io.ReadCloser), args.Get(1).(int64), args.Error(2)
}

// MockOsFs mocks the OsFs interface
type MockOsFs struct {
	mock.Mock
}

func (m *MockOsFs) OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	args := m.Called(name, flag, perm)
	return args.Get(0).(*os.File), args.Error(1)
}

func (m *MockOsFs) Remove(name string) error {
	args := m.Called(name)
	return args.Error(0)
}


// MockNotifyer mocks the domain.Notifyer interface
type MockNotifyer struct {
	mock.Mock
}

func (m *MockNotifyer) Notify(notification domain.Notification) error {
	args := m.Called(notification)
	return args.Error(0)
}

// MockDatabase mocks the domain.Database interface
type MockDatabase[T any] struct {
	mock.Mock
}

func (m *MockDatabase[T]) Save(key string, value domain.Video) error {
	m.Called(key, value)

	return nil
}

func (m *MockDatabase[T]) Get(key string) (domain.Video, error) {
	args := m.Called(key)
	return args.Get(0).(domain.Video), nil
}

func (m *MockDatabase[T]) Remove(id string) error {
	m.Called(id)

	return nil
}

// MockProgressBar mocks the domain.ProgressBar interface
type MockProgressBar struct {
	mock.Mock
}

func (m *MockProgressBar) Start(total int64) {
	m.Called(total)
}

func (m *MockProgressBar) Update(current int64) {
	m.Called(current)
}

func (m *MockProgressBar) Finish() {
	m.Called()
}

func (m *MockProgressBar) GetCurrent() int64 {
	args := m.Called()
	return args.Get(0).(int64)
}

func (m *MockProgressBar) GetTotal() int64 {
	args := m.Called()
	return args.Get(0).(int64)
}

func TestNewKkdaiDownloader(t *testing.T) {
	mockNotifyer := new(MockNotifyer)
	mockDB := new(MockDatabase[domain.Video])
	mockYtClient := new(MockYoutubeClient)
	mockFs := new(MockOsFs)
	cfg := &config.Config{}

	downloader := NewKkdaiDownloader(mockNotifyer, mockDB, mockYtClient, mockFs, cfg)

	assert.NotNil(t, downloader)
	assert.Equal(t, mockNotifyer, downloader.notifyer)
	assert.Equal(t, mockDB, downloader.db)
	assert.Equal(t, mockYtClient, downloader.ytClient)
	assert.Equal(t, mockFs, downloader.fs)
	assert.Equal(t, cfg, downloader.cfg)
}

func TestKkdaiDownloader_Download_Success(t *testing.T) {
	mockNotifyer := new(MockNotifyer)
	mockDB := new(MockDatabase[domain.Video])
	mockYtClient := new(MockYoutubeClient)
	mockFs := new(MockOsFs)
	mockProgressBar := new(MockProgressBar)
	cfg := &config.Config{VideoDir: t.TempDir()}

	downloader := NewKkdaiDownloader(mockNotifyer, mockDB, mockYtClient, mockFs, cfg)

	testVideo := domain.Video{URL: "test_video_id", Requester: "test_requester"}
	ytVideo := &youtube.Video{Title: "Test Video Title", Formats: []youtube.Format{{MimeType: "video/mp4; codecs=\"avc1.64001F\"", Quality: "hd720"}}}
	mockReadCloser := io.NopCloser(strings.NewReader("video_content"))

	mockYtClient.On("GetVideo", testVideo.URL).Return(ytVideo, nil)
	mockYtClient.On("GetStream", ytVideo, &ytVideo.Formats[0]).Return(mockReadCloser, int64(len("video_content")), nil)
	// Create a real temporary file for the mock to return
	tmpFile, err := os.CreateTemp(cfg.VideoDir, "test_video_*.mp4")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name()) // Clean up the temporary file after the test
	defer tmpFile.Close()           // Close the file handle

	mockFs.On("OpenFile", mock.Anything, mock.Anything, mock.Anything).Return(tmpFile, nil)
	mockDB.On("Save", mock.Anything, mock.Anything).Return().Once()
	mockNotifyer.On("Notify", mock.Anything).Return(nil).Once()
	mockProgressBar.On("Start", mock.Anything).Return().Once()
	mockProgressBar.On("Update", mock.Anything).Return().Times(1) // Called once at the end by the TeeReader
	mockProgressBar.On("Finish").Return().Once()

	err = downloader.Download(testVideo, mockProgressBar)
	assert.NoError(t, err)

	mockYtClient.AssertExpectations(t)
	mockFs.AssertExpectations(t)
	mockDB.AssertExpectations(t)
	mockNotifyer.AssertExpectations(t)
	mockProgressBar.AssertExpectations(t)

	// Verify the content of the temporary file
	content, err := os.ReadFile(tmpFile.Name())
	assert.NoError(t, err)
	assert.Equal(t, "video_content", string(content))
}

func TestKkdaiDownloader_Download_GetVideoError(t *testing.T) {
	mockNotifyer := new(MockNotifyer)
	mockDB := new(MockDatabase[domain.Video])
	mockYtClient := new(MockYoutubeClient)
	mockFs := new(MockOsFs)
	mockProgressBar := new(MockProgressBar)
	cfg := &config.Config{}

	downloader := NewKkdaiDownloader(mockNotifyer, mockDB, mockYtClient, mockFs, cfg)

	testVideo := domain.Video{URL: "test_video_id"}
	expectedErr := errors.New("failed to get video")

	mockYtClient.On("GetVideo", testVideo.URL).Return(&youtube.Video{}, expectedErr)

	err := downloader.Download(testVideo, mockProgressBar)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error fetching video info")

	mockYtClient.AssertExpectations(t)
	mockNotifyer.AssertNotCalled(t, "Notify", mock.Anything)
	mockDB.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
	mockProgressBar.AssertNotCalled(t, "Start", mock.Anything)
}

func TestKkdaiDownloader_Download_GetStreamError(t *testing.T) {
	mockNotifyer := new(MockNotifyer)
	mockDB := new(MockDatabase[domain.Video])
	mockYtClient := new(MockYoutubeClient)
	mockFs := new(MockOsFs)
	mockProgressBar := new(MockProgressBar)
	cfg := &config.Config{}

	downloader := NewKkdaiDownloader(mockNotifyer, mockDB, mockYtClient, mockFs, cfg)

	testVideo := domain.Video{URL: "test_video_id"}
	ytVideo := &youtube.Video{Formats: []youtube.Format{{MimeType: "video/mp4"}}}
	expectedErr := errors.New("failed to get stream")

	mockYtClient.On("GetVideo", testVideo.URL).Return(ytVideo, nil)
	mockYtClient.On("GetStream", ytVideo, &ytVideo.Formats[0]).Return(&DummyReadCloser{}, int64(0), expectedErr)

	err := downloader.Download(testVideo, mockProgressBar)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error getting video stream")

	mockYtClient.AssertExpectations(t)
	mockNotifyer.AssertNotCalled(t, "Notify", mock.Anything)
	mockDB.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
	mockProgressBar.AssertNotCalled(t, "Start", mock.Anything)
}

func TestKkdaiDownloader_Download_OpenFileError(t *testing.T) {
	mockNotifyer := new(MockNotifyer)
	mockDB := new(MockDatabase[domain.Video])
	mockYtClient := new(MockYoutubeClient)
	mockFs := new(MockOsFs)
	mockProgressBar := new(MockProgressBar)
	cfg := &config.Config{VideoDir: t.TempDir()}

	downloader := NewKkdaiDownloader(mockNotifyer, mockDB, mockYtClient, mockFs, cfg)

	testVideo := domain.Video{URL: "test_video_id"}
	ytVideo := &youtube.Video{Formats: []youtube.Format{{MimeType: "video/mp4"}}}
	mockReadCloser := io.NopCloser(strings.NewReader("video_content"))
	expectedErr := errors.New("failed to open file")

	mockYtClient.On("GetVideo", testVideo.URL).Return(ytVideo, nil)
	mockYtClient.On("GetStream", ytVideo, &ytVideo.Formats[0]).Return(mockReadCloser, int64(len("video_content")), nil)
	mockFs.On("OpenFile", mock.Anything, mock.Anything, mock.Anything).Return(&os.File{}, expectedErr) // Return a dummy file and an error

	err := downloader.Download(testVideo, mockProgressBar)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error creating file")

	mockYtClient.AssertExpectations(t)
	mockFs.AssertExpectations(t)
	mockNotifyer.AssertNotCalled(t, "Notify", mock.Anything)
	mockDB.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
	mockProgressBar.AssertNotCalled(t, "Start", mock.Anything)
}

func TestKkdaiDownloader_Finalize_Success(t *testing.T) {
	mockNotifyer := new(MockNotifyer)
	mockDB := new(MockDatabase[domain.Video])
	mockYtClient := new(MockYoutubeClient)
	mockFs := new(MockOsFs)
	cfg := &config.Config{}

	downloader := NewKkdaiDownloader(mockNotifyer, mockDB, mockYtClient, mockFs, cfg)
	notification := domain.Notification{Title: "Test Title"}

	mockNotifyer.On("Notify", notification).Return(nil).Once()

	err := downloader.Finalize(notification)
	assert.NoError(t, err)

	mockNotifyer.AssertExpectations(t)
}

func TestKkdaiDownloader_Finalize_NotifyError(t *testing.T) {
	mockNotifyer := new(MockNotifyer)
	mockDB := new(MockDatabase[domain.Video])
	mockYtClient := new(MockYoutubeClient)
	mockFs := new(MockOsFs)
	cfg := &config.Config{}

	downloader := NewKkdaiDownloader(mockNotifyer, mockDB, mockYtClient, mockFs, cfg)
	notification := domain.Notification{Title: "Test Title"}
	expectedErr := errors.New("notify error")

	mockNotifyer.On("Notify", notification).Return(expectedErr).Once()

	err := downloader.Finalize(notification)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "erro to notify user")

	mockNotifyer.AssertExpectations(t)
}

func TestKkdaiDownloader_Cancel_Success(t *testing.T) {
	mockFs := new(MockOsFs)
	// For testing purposes, we can create a temporary file and mock its Close and Remove operations
	tmpFile, err := os.CreateTemp("", "test_cancel_*.tmp")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name()) // Ensure cleanup

	mockFs.On("Remove", tmpFile.Name()).Return(nil).Once()

	downloader := &KkdaiDownloader{fs: mockFs}
	err = downloader.Cancel(*tmpFile)
	assert.NoError(t, err)
	mockFs.AssertExpectations(t)
}

func TestKkdaiDownloader_Cancel_CloseFileError(t *testing.T) {
	mockFs := new(MockOsFs)
	// Create a dummy file that we can control the Close error
	mockFile := &mockFileWithError{
		File:     &os.File{}, // A dummy file, Close() will be mocked
		closeErr: errors.New("close error"),
	}

	downloader := &KkdaiDownloader{fs: mockFs}
	err := downloader.Cancel(*mockFile.File) // Pass the embedded os.File
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error closing file")
	mockFs.AssertNotCalled(t, "Remove", mock.Anything) // Remove should not be called if Close fails
}

func TestKkdaiDownloader_Cancel_RemoveFileError(t *testing.T) {
	mockFs := new(MockOsFs)
	tmpFile, err := os.CreateTemp("", "test_cancel_*.tmp")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name()) // Ensure cleanup

	expectedErr := errors.New("remove error")
	mockFs.On("Remove", tmpFile.Name()).Return(expectedErr).Once()

	downloader := &KkdaiDownloader{fs: mockFs}
	err = downloader.Cancel(*tmpFile)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error deleting file")
	mockFs.AssertExpectations(t)
}

// Helper struct to mock os.File.Close()
type mockFileWithError struct {
	*os.File
	closeErr error
}

func (m *mockFileWithError) Close() error {
	return m.closeErr
}

// DummyReadCloser returns an error on Read
type DummyReadCloser struct{}

func (d *DummyReadCloser) Read(p []byte) (n int, err error) {
	return 0, errors.New("dummy read error")
}

func (d *DummyReadCloser) Close() error {
	return nil
}

func TestDefaultYoutubeClient_GetVideo(t *testing.T) {
	client := &DefaultYoutubeClient{}
	// This test will hit the actual YouTube API, so it might be slow and fragile.
	// For a true unit test, one would mock the underlying yt.Client.
	// However, for integration-like testing, this provides some value.
	videoID := "dQw4w9WgXcQ" // A famous Rick Astley video ID
	video, err := client.GetVideo(videoID)

	assert.NoError(t, err)
	assert.NotNil(t, video)
	assert.NotEmpty(t, video.Title)
	assert.True(t, len(video.Formats) > 0)
}

func TestDefaultYoutubeClient_GetStream(t *testing.T) {
	client := &DefaultYoutubeClient{}
	// This also hits the actual YouTube API.
	videoID := "dQw4w9WgXcQ"
	ytClient := youtube.Client{}
	video, err := ytClient.GetVideo(videoID)
	assert.NoError(t, err)
	assert.NotNil(t, video)

	format := &video.Formats[0] // Get the first format

	stream, size, err := client.GetStream(video, format)
	assert.NoError(t, err)
	assert.NotNil(t, stream)
	assert.True(t, size > 0)

	defer stream.Close()
	// Read some data to ensure stream is working
	buffer := make([]byte, 1024)
	n, err := stream.Read(buffer)
	assert.NoError(t, err)
	assert.True(t, n > 0)
}

func TestDefaultOsFs_OpenFile(t *testing.T) {
	fs := &DefaultOsFs{}
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "testfile.txt")

	file, err := fs.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	assert.NoError(t, err)
	assert.NotNil(t, file)
	defer file.Close()
	defer os.Remove(filePath)

	_, err = file.WriteString("test content")
	assert.NoError(t, err)

	content, err := os.ReadFile(filePath)
	assert.NoError(t, err)
	assert.Equal(t, "test content", string(content))
}

func TestDefaultOsFs_Remove(t *testing.T) {
	fs := &DefaultOsFs{}
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "file_to_remove.txt")

	// Create a dummy file to remove
	err := os.WriteFile(filePath, []byte("dummy"), 0644)
	assert.NoError(t, err)

	err = fs.Remove(filePath)
	assert.NoError(t, err)

	_, err = os.Stat(filePath)
	assert.True(t, os.IsNotExist(err)) // Ensure file is removed
}

func TestProgressWriter_Write(t *testing.T) {
	mockProgressBar := new(MockProgressBar)
	pw := &progressWriter{progress: mockProgressBar}

	mockProgressBar.On("Update", int64(5)).Return().Once()
	mockProgressBar.On("Update", int64(10)).Return().Once()

	n, err := pw.Write([]byte("hello"))
	assert.NoError(t, err)
	assert.Equal(t, 5, n)
	assert.Equal(t, int64(5), pw.current)

	n, err = pw.Write([]byte("world"))
	assert.NoError(t, err)
	assert.Equal(t, 5, n)
	assert.Equal(t, int64(10), pw.current)

	mockProgressBar.AssertExpectations(t)
}
