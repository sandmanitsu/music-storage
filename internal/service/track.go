package service

import (
	"music_storage/internal/domain"
	"music_storage/internal/logger"
	"music_storage/internal/repository"
	"net/url"
	"strings"
)

type TrackManager interface {
	List(params url.Values) ([]domain.Track, error)
	Delete(id int) error
	Text(id int) ([]string, error)
	Update(data TrackInput) error
}

type TrackService struct {
	repos  repository.Track
	logger *logger.Logger
}

func NewTrackService(logger *logger.Logger, repo repository.Track) *TrackService {
	return &TrackService{repos: repo}
}

func (s *TrackService) List(params url.Values) ([]domain.Track, error) {
	filter := map[string]interface{}{}
	filter["id"] = params.Get("id")
	filter["group_name"] = params.Get("group_name")
	filter["song"] = params.Get("song")
	filter["song_text"] = params.Get("text")
	filter["realise_date"] = params.Get("realise_date")
	filter["link"] = params.Get("link")

	tracks, err := s.repos.Get(repository.ListParamInput{
		Filter: filter,
		Offset: params.Get("offset"),
		Limit:  params.Get("limit"),
	})
	if err != nil {
		// s.logger.Debug("error executing sql") //, sl.Err(err))
		return nil, err
	}

	return tracks, nil
}

func (s *TrackService) Text(id int) ([]string, error) {
	text, err := s.repos.Text(id)
	if err != nil {
		return nil, err
	}

	chorus := strings.Split(text, "%chorus%")
	return chorus, nil
}

func (s *TrackService) Delete(id int) error {
	return s.repos.Delete(id)
}

type TrackInput struct {
	ID          int     `json:"id"`
	GroupName   *string `json:"group_name"`
	Song        *string `json:"song"`
	Text        *string `json:"text"`
	RealiseDate *string `json:"realise_date"`
	Link        *string `json:"link"`
}

func (s *TrackService) Update(input TrackInput) error {
	data := map[string]interface{}{}
	data["group_name"] = input.GroupName
	data["song"] = input.Song
	data["song_text"] = input.Text
	data["realise_date"] = input.RealiseDate
	data["link"] = input.Link

	return s.repos.Update(input.ID, data)
}
