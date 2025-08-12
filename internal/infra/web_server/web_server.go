package webserver

import (
	"context"
	"downloader/internal/domain"
	"downloader/internal/infra/progress"
	"downloader/internal/usecase"
	"downloader/pkg/config"
	logger "downloader/pkg/log"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
)

var log = logger.GetLogger("web_server")

type WebServer struct {
	server     *http.Server
	downloadUC usecase.DownloadVideoUseCase
}

type returnHttp struct {
	Message string `json:"message"`
}

func NewWebServer(downloader domain.Downloader) *WebServer {
	return &WebServer{downloadUC: usecase.DownloadVideoUseCase{Downloader: downloader}}
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
		log.Info(fmt.Sprintf("server listen on port %d", port))
		if err := w.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error(fmt.Sprintf("listen: %v", err))
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
	log.Info("server stopped")
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

	go ws.downloadUC.Execute(usecase.Solicitation{URL: url, Requester: requester}, progress.NewTerminalProgressBar())
	json.NewEncoder(w).Encode(returnHttp{Message: "Download iniciado"})
}

func (ws *WebServer) download(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	videoRoot := config.GetConfig().VideoDir

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

	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	w.Header().Set("Content-Type", "video/mp4")
	w.Header().Set("Cache-Control", "private, max-age=86400")

	log.Info(fmt.Sprintf("Download de %s iniciado", filename))
	http.ServeContent(w, r, filename, stat.ModTime().UTC(), f)
	tic := time.Tick(30 * time.Minute)

	<-tic
	err = os.Remove(cleanPath)
	if err != nil {
		log.Error(fmt.Sprintf("Nao foi possivel excluir o arquivo: %s", err))
	}
}
