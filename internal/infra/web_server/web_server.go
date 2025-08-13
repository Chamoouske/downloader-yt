package webserver

import (
	"context"
	"downloader/internal/domain"
	"downloader/internal/infra/progress"
	"downloader/internal/usecase"
	"downloader/pkg/config"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

type WebServer struct {
	server     *http.Server
	downloadUC usecase.DownloadVideoUseCase
	db         domain.Database[domain.Video]
	cfg        *config.Config
	log        *slog.Logger
}

type returnHttp struct {
	Message string `json:"message"`
}

func NewWebServer(downloader domain.Downloader, db domain.Database[domain.Video], cfg *config.Config, appLogger *slog.Logger) *WebServer {
	return &WebServer{downloadUC: usecase.DownloadVideoUseCase{Downloader: downloader}, db: db, cfg: cfg, log: appLogger}
}

func (w *WebServer) Start(port int) {
	mux := mux.NewRouter()
	mux.HandleFunc("/video/download", w.addVideoNaFilaDeDownload).Methods("GET")
	mux.HandleFunc("/video/{id}", w.download).Methods("GET")

	w.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	go func() {
		w.log.Info(fmt.Sprintf("server listen on port %d", port))
		if err := w.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			w.log.Error(fmt.Sprintf("listen: %v", err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	w.Stop()
}

func (w *WebServer) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = w.server.Shutdown(ctx)
	w.log.Info("server stopped")
}

func (ws *WebServer) addVideoNaFilaDeDownload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	url := r.URL.Query().Get("url")
	if url == "" {
		http.Error(w, "URL parameter is required", http.StatusBadRequest)
		return
	}
	requester := r.URL.Query().Get("requester")
	if requester == "" {
		http.Error(w, "requester parameter is required", http.StatusBadRequest)
		return
	}

	progressBarClient := &progress.DefaultProgressBarClient{}
	progressBar := progress.NewTerminalProgressBar(progressBarClient)

	go ws.downloadUC.Execute(usecase.Solicitation{URL: url, Requester: requester}, progressBar)
	json.NewEncoder(w).Encode(returnHttp{Message: "Download iniciado"})
}

func (ws *WebServer) download(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	video, err := ws.db.Get(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	videoRoot := ws.cfg.VideoDir

	filename := id + ".mp4"
	fullPath := filepath.Join(videoRoot, filename)

	cleanRoot, _ := filepath.Abs(videoRoot)
	cleanPath, err := filepath.Abs(fullPath)
	if err != nil || len(cleanPath) < len(cleanRoot) || cleanPath[:len(cleanRoot)] != cleanRoot {
		http.Error(w, "invalid path", http.StatusForbidden)
		return
	}

	f, err := os.Open(cleanPath)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.mp4"`, video.Filename))
	w.Header().Set("Content-Type", "video/mp4")
	w.Header().Set("Cache-Control", "private, max-age=86400")

	ws.log.Info(fmt.Sprintf("Download de %s iniciado", video.Filename))

	cleanupFn := func() {
		f.Close()
		if err := os.Remove(cleanPath); err != nil {
			ws.log.Error(fmt.Sprintf("Erro ao remover arquivo: %s", err))
		} else {
			ws.log.Info(fmt.Sprintf("Arquivo %s removido após download", cleanPath))
		}
	}

	seeker := &cleanupReadSeeker{
		file:    f,
		reader:  f,
		seeker:  f,
		cleanup: cleanupFn,
	}
	defer seeker.triggerCleanup()

	http.ServeContent(w, r, filename, stat.ModTime().UTC(), seeker)
}

type cleanupReadSeeker struct {
	file    *os.File
	reader  io.Reader
	seeker  io.Seeker
	cleanup func()
	once    sync.Once
}

func (r *cleanupReadSeeker) Read(p []byte) (int, error) {
	n, err := r.reader.Read(p)
	if err != nil {
		r.triggerCleanup()
	}
	return n, err
}

func (r *cleanupReadSeeker) Seek(offset int64, whence int) (int64, error) {
	return r.seeker.Seek(offset, whence)
}

func (r *cleanupReadSeeker) triggerCleanup() {
	r.once.Do(r.cleanup)
}
