package storage

import "errors"

var (
	// ErrNotFound - object not found in DB.
	ErrNotFound = errors.New("not found")

	// ErrConflict - conflicting request (e.g. unique constraint violated)
	ErrConflict = errors.New("conflict")
)