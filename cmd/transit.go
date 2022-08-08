package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/httplog"
	"github.com/joho/godotenv"
	"github.com/linhmtran168/511transit/internal/data/memory"
	"github.com/linhmtran168/511transit/internal/handler"
	transitapi "github.com/linhmtran168/511transit/internal/transit-api"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Default configuration for http logger
	httplogOptions := httplog.Options{
		JSON:            true,
		TimeFieldFormat: zerolog.TimeFormatUnix,
	}
	initConfiguration(&httplogOptions)

	addr := os.Getenv("SERVER_ADDR")
	httpLogger := httplog.NewLogger(os.Getenv("APP_NAME"), httplogOptions)
	apiClient := transitapi.NewTransitAPI()
	repository := memory.NewMemoryRepository(apiClient)
	appHandler := handler.NewAppHandler(httpLogger, repository)
	server := &http.Server{
		Addr:    addr,
		Handler: appHandler,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("Could not start server")
		}
	}()

	defer func() {
		gracefulShutdown(server)
	}()

	log.Info().Msgf("Start server at %s", addr)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
	log.Info().Msg("os.Interrupt - Shutting Down")

}

func initConfiguration(httplogOptions *httplog.Options) {
	// Local env
	if os.Getenv("ENV") == "" {
		// Setting httplog options also update global zerelog options, so we only need to override httplog options
		*httplogOptions = httplog.Options{JSON: false}
		err := godotenv.Load(".env.local")
		if err != nil {
			log.Fatal().Err(err).Msg("Error loading local .env file")
		}
	}
	// Load base configuration
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading base .env file")
	}

}

func gracefulShutdown(server *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Could not gracefully shutdown server")
	} else {
		log.Info().Msg("Server gracefully stopped")
	}
}
