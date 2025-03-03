package repository

import (
	"music_storage/internal/logger"
	"music_storage/internal/storage/postgresql"
)

type Repositories struct {
	Track Track
}

// Create Repository struct that include other data repositories
func NewRepository(logger *logger.Logger, storage *postgresql.Storage) *Repositories {
	return &Repositories{Track: NewTrackRepository(logger, storage.DB)}
}
