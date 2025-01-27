package repository

import (
	"database/sql"
	"fmt"
	"log"
	"music_storage/internal/domain"
	"strings"
)

const (
	table         = "tracks" // table name
	defaultLimit  = 10       // max amount tracks per page
	defaultOffset = 0
)

type Track interface {
	Get(params ListParamInput) ([]domain.Track, error)
}

type TrackRepository struct {
	db *sql.DB
}

// Create Track Repository
func NewTrackRepository(db *sql.DB) *TrackRepository {
	return &TrackRepository{db: db}
}

type ListParamInput struct {
	Filter        map[string]interface{}
	Limit, Offset string
}

// Return tracks by filter params
func (r *TrackRepository) Get(params ListParamInput) ([]domain.Track, error) {
	where, values := r.whereStatement(params)

	if params.Limit != "" {
		values = append(values, params.Limit)
	} else {
		values = append(values, defaultLimit)
	}

	if params.Offset != "" {
		values = append(values, params.Offset)
	} else {
		values = append(values, defaultOffset)
	}

	quary := fmt.Sprintf(
		"SELECT id, group_name, song, text, realise_date, link FROM %s%s LIMIT ? OFFSET ?",
		table,
		where,
	)
	fmt.Println(quary)
	rows, err := r.db.Query(quary, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tracks []domain.Track
	for rows.Next() {
		var track domain.Track
		if err := rows.Scan(&track.ID, &track.GroupName, &track.Song, &track.Text, &track.RealiseDate, &track.Link); err != nil {
			log.Fatal(err)
		}
		tracks = append(tracks, track)
	}

	return tracks, nil
}

// Prepare where statement
func (r *TrackRepository) whereStatement(params ListParamInput) (string, []interface{}) {
	var values []interface{}
	var where []string

	for k, v := range params.Filter {
		if v == "" {
			continue
		}

		if k == "text" {
			values = append(values, fmt.Sprintf("%%%s%%", v))
			where = append(where, fmt.Sprintf("%s LIKE ?", k))
			continue
		}

		values = append(values, v)
		where = append(where, fmt.Sprintf("%s = ?", k))
	}

	if len(where) == 0 {
		return "", values
	}

	return fmt.Sprintf(" WHERE %s", strings.Join(where, " AND ")), values
}
