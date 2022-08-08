package handler

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"
	"github.com/linhmtran168/511transit/internal/data"
	"github.com/rs/zerolog"
)

func NewAppHandler(logger zerolog.Logger, repository data.DataRepository) http.Handler {
	r := chi.NewRouter()

	r.Use(httplog.RequestLogger(logger))
	r.Use(middleware.RealIP)
	r.Use(middleware.Compress(5))
	r.Use(middleware.Heartbeat("/ping"))

	// If other environment besides local, serve static files from dist folder
	if os.Getenv("ENV") != "" {
		workDir, _ := os.Getwd()
		filesDir := filepath.Join(workDir, "web/dist")
		StaticServe(r, "/", filesDir)
	}

	// Websocket endpoint
	websocketHandler := NewWebSocketHandler(repository)
	r.Get("/ws", websocketHandler.ConnectHandler)

	return r
}
