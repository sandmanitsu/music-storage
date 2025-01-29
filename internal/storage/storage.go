package storage

import "errors"

var (
	ErrRecordsNotFound = errors.New("error record not found")
	ErrNoRowsAffected  = errors.New("no rows affected")
)
