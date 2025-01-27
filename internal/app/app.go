package app

import (
	"fmt"
	"log/slog"
	"music_storage/internal/config"
	log "music_storage/internal/logger"
	"music_storage/internal/repository"
	"music_storage/internal/service"
	"music_storage/internal/storage/sqlite"
	"music_storage/internal/transport/router"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func Run(config *config.Config, logger *slog.Logger) {
	logger.Info("starting app")

	storage, err := sqlite.NewSQLite(config.DB.StoragePath)
	if err != nil {
		logger.Error("failed to init sqlite storage", log.Err(err))
		os.Exit(1)
	}

	repositories := repository.NewRepository(storage)
	services := service.NewService(logger, repositories)

	handler := router.NewHandler(services)
	router := handler.Init()

	server := http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.Host, config.Port),
		Handler: router,
	}

	logger.Info(fmt.Sprintf("server starting on %s:%d", config.Host, config.Port))
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			logger.Error("error starting server", log.Err(err))
			os.Exit(1)
		}
	}()

	logger.Info("server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	logger.Info("graceful shutdown")
}
