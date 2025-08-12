package webserver

import (
	"context"
	"downloader/internal/domain"
	"downloader/internal/infra/progress"
	"downloader/internal/usecase"
	logger "downloader/pkg/log"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var log = logger.GetLogger("web_server")

type WebServer struct {
	server     *http.Server
	downloadUC usecase.DownloadVideoUseCase
}

func NewWebServer(downloader domain.Downloader) *WebServer {
	return &WebServer{downloadUC: usecase.DownloadVideoUseCase{Downloader: downloader}}
}

func (w *WebServer) Start(port int) {
	mux := http.NewServeMux()
	mux.HandleFunc("/download", w.listItems)

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

func (ws *WebServer) listItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	url := r.URL.Query().Get("url")
	if url == "" {
		http.Error(w, "URL parameter is required", http.StatusBadRequest)
		return
	}

	go ws.downloadUC.Execute(url, progress.NewTerminalProgressBar())
	json.NewEncoder(w).Encode(returnHttp{Message: "Download iniciado"})
}

type returnHttp struct {
	Message string `json:"message"`
}
