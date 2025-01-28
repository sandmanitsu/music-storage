package service

import (
	"log/slog"
	"music_storage/internal/repository"
)

type Service struct {
	Track TrackManager
}

func NewService(logger *slog.Logger, repos *repository.Repositories) *Service {
	return &Service{
		Track: NewTrackService(logger, repos.Track),
	}
}
