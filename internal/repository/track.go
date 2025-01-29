package repository

import (
	"database/sql"
	"fmt"
	"log/slog"
	"music_storage/internal/domain"
	"music_storage/internal/logger"
	"music_storage/internal/storage"
	"reflect"
	"strings"
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
	Update(id int, data map[string]interface{}) error
}

type TrackRepository struct {
	logger *logger.Logger
	db     *sql.DB
}

// Create Track Repository
func NewTrackRepository(logger *logger.Logger, db *sql.DB) *TrackRepository {
	return &TrackRepository{db: db, logger: logger}
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
		r.logger.Debug("error executing query", slog.String("query", query), r.logger.Err(err))

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
	values := make([]interface{}, 0, 5)
	where := make([]string, 0, 5)
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
		r.logger.Debug("error deleting track", r.logger.Err(err))

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
		r.logger.Debug(fmt.Sprintf("db error: getting song_text. query - %s", query), r.logger.Err(err))

		return "", err
	}

	return text, nil
}

func (r *TrackRepository) Update(id int, data map[string]interface{}) error {
	fields, values := parseData(data)

	values = append(values, id)
	idPlaceholderNum := len(values)

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = $%d", table, strings.Join(fields, ","), idPlaceholderNum)
	result, err := r.db.Exec(query, values...)
	if err != nil {
		r.logger.Debug(fmt.Sprintf("error insert into %s query: %s", table, query), r.logger.Err(err))

		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		r.logger.Debug(fmt.Sprintf("error getting affected rows %s", table), r.logger.Err(err))
	}

	if rows == 0 {
		r.logger.Debug(fmt.Sprintf("error no rows affected %s", table))

		return storage.ErrNoRowsAffected
	}

	return nil
}

// Prepare field and values
// Filtering values on nil
func parseData(data map[string]interface{}) ([]string, []interface{}) {
	fields := make([]string, 0, 5)
	values := make([]interface{}, 0, 5)
	i := 1
	for f, v := range data {
		if reflect.ValueOf(v).IsNil() {
			continue
		}

		values = append(values, v)
		fields = append(fields, fmt.Sprintf("%s = $%d ", f, i))
		i++
	}

	return fields, values
}
