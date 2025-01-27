package repository

import (
	"database/sql"
	"fmt"
	"log"
	"music_storage/internal/domain"
	"strings"
)

const (
	table     = "tracks" // table name
	pageLimit = 10       // max amount tracks per page
)

type Track interface {
	Get(params ListParamInput)
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
// todo. realise filter
func (r *TrackRepository) Get(params ListParamInput) {
	where, values := r.whereStatement(params)

	if params.Limit != "" {
		values = append(values, params.Limit)
	} else {
		values = append(values, pageLimit)
	}

	if params.Offset != "" {
		values = append(values, params.Offset)
	}

	quary := fmt.Sprintf("SELECT id, group_name, song, text, realise_date, link FROM %s%s LIMIT = ? OFFSET = ?", table, where)
	fmt.Println(quary)
	// rows, err := r.db.Query(fmt.Sprintf("SELECT id, group_name, song, text, realise_date, link FROM %s", table))
	rows, err := r.db.Query(quary, values...) // ???? here is error!
	if err != nil {
		log.Fatal(err) // todo. return error + create const error in storage to it
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

	fmt.Println(tracks)
}

// Prepare where statement
func (r *TrackRepository) whereStatement(params ListParamInput) (string, []interface{}) {
	var values []interface{}
	var where []string

	for k, v := range params.Filter {
		if v == "" {
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
