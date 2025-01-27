package service

import (
	"fmt"
	"music_storage/internal/repository"
)

type TrackManager interface {
	List()
}

type TrackService struct {
	repos repository.Track
}

func NewTrackService(repo repository.Track) *TrackService {
	return &TrackService{repos: repo}
}

func (s *TrackService) List() {
	fmt.Println("track service")

	s.repos.Get()
}
