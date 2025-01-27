package main

import (
	"log"
	"log/slog"
	"music_storage/internal/app"
	"music_storage/internal/config"
	"music_storage/internal/logger"
)

func main() {
	log.Println("Getting config...")
	config := config.MustLoad()

	log.Println("logger initializing...")
	logger := logger.NewLogger(config.Env)
	logger.Info("logger started!", slog.String("env", config.Env))

	app.Run(config, logger)
}
