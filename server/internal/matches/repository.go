package matches

import (
	"context"
	"database/sql"
)

type Repository interface {
	ListByUser(ctx context.Context, userID string) ([]Match, error)
	Create(ctx context.Context, m Match) error
}


type repo struct {db *sql.DB}

func NewRepository(db *sql.DB) Repository {
	return &repo{db: db}
}


func (r *repo) ListByUser(ctx context.Context, userID string) ([]Match, error) {
	// Implementation
}

func (r *repo) Create(ctx context.Context, m Match) error {
	// Implementation
}
