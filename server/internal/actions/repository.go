package actions

import (
	"context"
	"database/sql"
	"net/http"
)

type Repository interface {
	List(ctx context.Context, viewerId string) ([]Action, error)
	Create(ctx context.Context, a Action) error
}

type repo struct {db *sql.DB}

func NewRepository(db *sql.DB) Repository { return &repo{db: db}}

func(r *repo) List(ctx context.Context, viewerId string) ([]Action, error) {
	// implementation
	return http.ErrAbortHandler
}

func (r *repo) Create(ctx context.Context, a Action) error {
	// implementation
	return http.ErrAbortHandler
}
