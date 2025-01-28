package repository

import (
	"database/sql"
	"fmt"
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

	limitPlaceholderNum := len(values) + 1
	if params.Limit != "" {
		values = append(values, params.Limit)
	} else {
		values = append(values, defaultLimit)
	}

	offsetPlaceholderNum := len(values) + 1
	if params.Offset != "" {
		values = append(values, params.Offset)
	} else {
		values = append(values, defaultOffset)
	}

	_ = where
	query := fmt.Sprintf(
		"SELECT id, group_name, song, song_text, TO_CHAR(realise_date::DATE, 'YYYY-MM-DD') AS realise_date, link FROM %s%s LIMIT $%d OFFSET $%d",
		table,
		where,
		limitPlaceholderNum,
		offsetPlaceholderNum,
	)
	fmt.Println(query)
	rows, err := r.db.Query(query, values...)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	var tracks []domain.Track
	for rows.Next() {
		var track domain.Track
		if err := rows.Scan(&track.ID, &track.GroupName, &track.Song, &track.Text, &track.RealiseDate, &track.Link); err != nil {
			return nil, err
		}
		tracks = append(tracks, track)
	}

	return tracks, nil
}

// Prepare where statement
func (r *TrackRepository) whereStatement(params ListParamInput) (string, []interface{}) {
	var values []interface{}
	var where []string
	i := 1
	for k, v := range params.Filter {
		if v == "" {
			continue
		}

		if k == "song_text" {
			values = append(values, fmt.Sprintf("%%%s%%", v))
			where = append(where, fmt.Sprintf("%s LIKE $%d", k, i))
			i++
			continue
		}

		values = append(values, v)
		where = append(where, fmt.Sprintf("%s = $%d", k, i))
		i++
	}

	if len(where) == 0 {
		return "", values
	}

	return fmt.Sprintf(" WHERE %s", strings.Join(where, " AND ")), values
}
