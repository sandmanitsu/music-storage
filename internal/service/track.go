package service

import (
	"encoding/json"
	"fmt"
	"music_storage/internal/domain"
	"music_storage/internal/logger"
	"music_storage/internal/repository"
	"net/url"
	"strings"
	"time"
)

type TrackManager interface {
	List(params url.Values) ([]domain.Track, error)
	Delete(id int) error
	Text(id int) ([]string, error)
	Update(data TrackInput) error
	Add(data TrackAddInput) error
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

type TrackAddInput struct {
	GroupName string `json:"group"`
	Song      string `json:"song"`
}

type SongInfo struct {
	RealiseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

func (s *TrackService) Add(input TrackAddInput) error {
	if input.GroupName == "" {
		return fmt.Errorf("error: fielf group is empty")
	}
	if input.Song == "" {
		return fmt.Errorf("error: song group is empty")
	}

	// todo. Пробросить верный url для запроса.
	// url := fmt.Sprintf("https://localhost/info?group=%s&song=%s", input.GroupName, input.Song)

	// response, err := http.Get(url)
	// if err != nil {
	// 	s.logger.Debug(fmt.Sprintf("error get song info - url: %s", url), s.logger.Err(err))
	// }
	// defer response.Body.Close()

	// body, err := io.ReadAll(response.Body)
	// if err != nil {
	// 	s.logger.Debug("error reading boby song info", s.logger.Err(err))

	// 	return err
	// }

	// ! mock data
	var err error
	body := []byte(`{"releaseDate":"16.07.2006","text":"Ooh baby, don't you know I suffer?\nOoh baby, canyou hear me moan?\nYou caught me under false pretenses\nHow longbefore you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou setmy soul alights","link":"https://www.youtube.com/watch?v=Xsp3_a-PMTw"}`)

	var songInfo SongInfo
	err = json.Unmarshal(body, &songInfo)
	if err != nil {
		s.logger.Debug("error marshal JSON with song info", s.logger.Err(err))

		return err
	}

	date, err := time.Parse("02.01.2006", songInfo.RealiseDate)
	if err != nil {
		s.logger.Debug("error parse date in song info", s.logger.Err(err))

		return err
	}

	err = s.repos.Add(domain.Track{
		GroupName:   input.GroupName,
		Song:        input.Song,
		Text:        songInfo.Text,
		RealiseDate: date.Format("2006-01-02"),
		Link:        songInfo.Link,
	})
	if err != nil {
		return err
	}

	return nil
}
