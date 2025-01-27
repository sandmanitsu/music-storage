package service

import (
	"log/slog"
	"music_storage/internal/domain"
	sl "music_storage/internal/logger"
	"music_storage/internal/repository"
	"net/url"
)

type TrackManager interface {
	List(params url.Values) ([]domain.Track, error)
}

type TrackService struct {
	repos  repository.Track
	logger *slog.Logger
}

func NewTrackService(logger *slog.Logger, repo repository.Track) *TrackService {
	return &TrackService{repos: repo}
}

func (s *TrackService) List(params url.Values) ([]domain.Track, error) {
	filter := map[string]interface{}{}
	filter["id"] = params.Get("id")
	filter["group_name"] = params.Get("group_name")
	filter["song"] = params.Get("song")
	filter["text"] = params.Get("text")
	filter["realise_date"] = params.Get("realise_date")
	filter["link"] = params.Get("link")

	tracks, err := s.repos.Get(repository.ListParamInput{
		Filter: filter,
		Offset: params.Get("offset"),
		Limit:  params.Get("limit"),
	})
	if err != nil {
		s.logger.Debug("error executing sql", sl.Err(err))
		return nil, err
	}

	return tracks, nil
}
