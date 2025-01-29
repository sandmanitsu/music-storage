package app

import (
	"fmt"
	"music_storage/internal/config"
	"music_storage/internal/logger"
	"music_storage/internal/repository"
	"music_storage/internal/service"
	"music_storage/internal/storage/postgresql"
	"music_storage/internal/transport/router"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// Initializes whole application
// @title           Music Storage API
// @version         1.0
// @description     This testing task to create music storage API

// @host      localhost:8080
// @BasePath  /api/v1

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func Run(config *config.Config, logger *logger.Logger) {
	logger.Info("starting app")

	storage, err := postgresql.NewPostgreSQL(config.DB)
	if err != nil {
		logger.Error("failed to init sqlite storage", logger.Err(err))
		os.Exit(1)
	}

	repositories := repository.NewRepository(logger, storage)
	services := service.NewService(logger, repositories)

	handler := router.NewHandler(services)
	router := handler.Init(logger)

	server := http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.Host, config.Port),
		Handler: router,
	}

	logger.Info(fmt.Sprintf("server starting on %s:%d", config.Host, config.Port))
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			logger.Error("error starting server", logger.Err(err))
			os.Exit(1)
		}
	}()

	logger.Info("server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	logger.Info("graceful shutdown")
}
