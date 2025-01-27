package service

import (
	"errors"
	"log/slog"
	"music_storage/internal/repository"
)

var (
	ErrValidateInputParams = errors.New("error validating input params")
)

type Service struct {
	Track TrackManager
}

func NewService(logger *slog.Logger, repos *repository.Repositories) *Service {
	return &Service{
		Track: NewTrackService(logger, repos.Track),
	}
}
