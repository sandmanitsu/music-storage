package app

import (
	"log/slog"
	"music_storage/internal/config"
	log "music_storage/internal/logger"
	"music_storage/internal/repository"
	"music_storage/internal/service"
	"music_storage/internal/storage/sqlite"
	"music_storage/internal/transport/router"
	"os"
)

func Run(config *config.Config, logger *slog.Logger) {
	logger.Info("starting app")

	storage, err := sqlite.NewSQLite(config.DB.StoragePath)
	if err != nil {
		logger.Error("failed to init sqlite storage", log.Err(err))
		os.Exit(1)
	}

	repositories := repository.NewRepository(storage)
	services := service.NewService(repositories)

	handlers := router.NewHandler(services)
	router := handlers.Init()
	router.Run()
}
