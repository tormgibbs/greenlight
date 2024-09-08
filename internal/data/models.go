package data

import (
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrRecordNotFound A custom error which is returned when a resource could not be found
var (
	ErrRecordNotFound = errors.New("record not found")
)

// Models struct which wraps the MovieModel
type Models struct {
	Movies MovieModel
}

// NewModels creates and returns a Model instance containing initialized MovieModel
func NewModels(db *pgxpool.Pool) Models {
	return Models{
		Movies: MovieModel{db},
	}
}
