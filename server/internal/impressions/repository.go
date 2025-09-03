package impressions

import (
	"context"
	"database/sql"
)

type Repository interface {
	List(ctx context.Context, viewerId string) ([]Impression, error)
	Create(ctx context.Context, imp Impression) error
}

type repo struct {db *sql.DB}

func NewRepository(db *sql.DB) Repository { return &repo{db: db}}

func (r *repo) List(ctx context.Context, viewerId string) ([] Impression, error) {
	// Implementation
}



func (r *repo) Create(ctx context.Context, imp Impression) error {
	// Implementation
}
