package repository

import (
	"database/sql"
	"fmt"
	"log"
	"music_storage/internal/domain"
)

type Track interface {
	Get()
}

type TrackRepository struct {
	db *sql.DB
}

// Create Track Repository
func NewTrackRepository(db *sql.DB) *TrackRepository {
	return &TrackRepository{db: db}
}

// Return tracks by filter params
// todo. realise filter
func (r *TrackRepository) Get() {
	rows, err := r.db.Query("SELECT id, group_name, song, text, realise_date, link FROM tracks")
	if err != nil {
		log.Fatal(err)
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
