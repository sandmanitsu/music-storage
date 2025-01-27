package service

import "music_storage/internal/repository"

type Service struct {
	Track TrackManager
}

func NewService(repos *repository.Repositories) *Service {
	return &Service{
		Track: NewTrackService(repos.Track),
	}
}
