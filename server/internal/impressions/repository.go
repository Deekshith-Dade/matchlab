package impressions

import (
	"context"
	"database/sql"
	"errors"
	"strings"
)

type Repository interface {
	List(ctx context.Context, viewerId string) ([]Impression, error)
	Create(ctx context.Context, imp Impression) error
}

type repo struct {db *sql.DB}

func NewRepository(db *sql.DB) Repository { return &repo{db: db}}

func (r *repo) List(ctx context.Context, viewerID string) ([]Impression, error) {
	viewerID = strings.TrimSpace(viewerID)
	if viewerID == "" {
		return nil, errors.New("viewer_id required")
	}

	const q = `
		SELECT viewer_id, viewed_id, rank, at
		FROM impressions
		WHERE viewer_id = $1
		ORDER BY at DESC;
	`

	rows, err := r.db.QueryContext(ctx, q, viewerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Impression
	for rows.Next() {
		var imp Impression
		if err := rows.Scan(&imp.ViewerID, &imp.ViewedID, &imp.Rank, &imp.At); err != nil {
			return nil, err
		}
		out = append(out, imp)
	}
	return out, rows.Err()
}

func (r *repo) Create(ctx context.Context, imp Impression) error {
	imp.ViewerID = strings.TrimSpace(imp.ViewerID)
	imp.ViewedID = strings.TrimSpace(imp.ViewedID)
	if imp.ViewerID == "" || imp.ViewedID == "" {
		return errors.New("viewer_id and viewed_id required")
	}
	// If imp.At is zero, let DB default to NOW()
	// Pass nil and COALESCE in SQL to avoid client clock issues.
	var atArg any
	if imp.At.IsZero() {
		atArg = nil
	} else {
		atArg = imp.At.UTC()
	}

	const q = `
		INSERT INTO impressions (viewer_id, viewed_id, rank, at)
		VALUES ($1, $2, $3, COALESCE($4, NOW()))
	`
	_, err := r.db.ExecContext(ctx, q, imp.ViewerID, imp.ViewedID, imp.Rank, atArg)
	return err
}
