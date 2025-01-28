package repository

import (
	"database/sql"
	"fmt"
	"log/slog"
	"music_storage/internal/domain"
	"music_storage/internal/storage"
	"strings"

	sl "music_storage/internal/logger"
)

const (
	table         = "tracks" // table name
	defaultLimit  = 10       // max amount tracks per page
	defaultOffset = 0
)

type Track interface {
	Get(params ListParamInput) ([]domain.Track, error)
	Delete(id int) error
	Text(id int) (string, error)
}

type TrackRepository struct {
	logger *slog.Logger
	db     *sql.DB
}

// Create Track Repository
func NewTrackRepository(logger *slog.Logger, db *sql.DB) *TrackRepository {
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
	rows, err := r.db.Query(query, values...)
	if err != nil {
		r.logger.Debug("error executing query", slog.String("query", query), sl.Err(err))
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

func (r *TrackRepository) Delete(id int) error {
	_, err := r.db.Exec(fmt.Sprintf("DELETE FROM %s WHERE id = $1", table), id)
	if err != nil {
		r.logger.Debug("error deleting track", sl.Err(err))
		return err
	}

	return nil
}

func (r *TrackRepository) Text(id int) (string, error) {
	var text string
	query := fmt.Sprintf("SELECT song_text FROM %s WHERE id = $1", table)
	err := r.db.QueryRow(query, id).Scan(&text)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", storage.ErrRecordsNotFound
		}
		r.logger.Debug(fmt.Sprintf("db error: getting song_text. query - %s", query), sl.Err(err))
		return "", err
	}
	return text, nil
}
