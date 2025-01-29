package service

import (
	"music_storage/internal/logger"
	"music_storage/internal/repository"
)

type Service struct {
	Track TrackManager
}

func NewService(logger *logger.Logger, repos *repository.Repositories) *Service {
	return &Service{
		Track: NewTrackService(logger, repos.Track),
	}
}
