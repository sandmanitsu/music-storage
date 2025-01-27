package service

import (
	"log/slog"
	"music_storage/internal/repository"
	"net/url"
)

type TrackManager interface {
	List(params url.Values) error
}

type TrackService struct {
	repos  repository.Track
	logger *slog.Logger
}

func NewTrackService(logger *slog.Logger, repo repository.Track) *TrackService {
	return &TrackService{repos: repo}
}

func (s *TrackService) List(params url.Values) error {
	filter := map[string]interface{}{}
	filter["id"] = params.Get("id")
	filter["group_name"] = params.Get("group_name")
	filter["song"] = params.Get("song")
	filter["text"] = params.Get("text")
	filter["realise_date"] = params.Get("realise_date")
	filter["link"] = params.Get("link")

	s.repos.Get(repository.ListParamInput{
		Filter: filter,
		Offset: params.Get("offset"),
		Limit:  params.Get("limit"),
	})

	return nil
}
