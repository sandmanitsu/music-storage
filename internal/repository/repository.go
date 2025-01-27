package repository

import "music_storage/internal/storage/sqlite"

type Repositories struct {
	Track Track
}

// Create Repository struct that include other data repositories
func NewRepository(storage *sqlite.Storage) *Repositories {
	return &Repositories{Track: NewTrackRepository(storage.DB)}
}
