package actions

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/deekshith-dade/matchlab/internal/matches"
)

type Repository interface {
	List(ctx context.Context, viewerId string) ([]Action, error)
	Create(ctx context.Context, a Action) (*matches.Match, error) 
}

type repo struct {db *sql.DB}

func NewRepository(db *sql.DB) Repository { return &repo{db: db}}

func(r *repo) List(ctx context.Context, viewerId string) ([]Action, error) {
	// implementation
	if viewerId == "" {
		return nil, errors.New("viewer_id required")
	}	

	const q = `
		SELECT viewer_id, viewed_id, kind, at
		FROM actions
		WHERE viewer_id = $1
		ORDER BY at DESC;
		`
	rows, err := r.db.QueryContext(ctx, q, viewerId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var out []Action
	for rows.Next() {
		var a Action
		if err := rows.Scan(&a.ViewerId, &a.ViewedId, &a.Kind, &a.At); err != nil {
			return nil, err
		}

		out = append(out, a)
	}


	return out, nil

}

func (r *repo) Create(ctx context.Context, a Action) (*matches.Match, error) {
	// implementation

	// Implement Checks Later
	a.ViewerId = strings.TrimSpace(a.ViewerId)
	a.ViewedId = strings.TrimSpace(a.ViewedId)
	if a.ViewerId == "" || a.ViewedId == "" {
		return nil, errors.New("ViewerID and ViewedId required")
	}

	kind := strings.ToLower(strings.TrimSpace(a.Kind))
	if kind == "" {
		return nil, errors.New("kind required")
	}

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	var atArg any
	if a.At.IsZero() {
		atArg = nil
	}else {
		atArg = a.At.UTC()
	}

	const ins = `
		INSERT INTO actions (viewer_id, viewed_id, kind, at)
		VALUES ($1, $2, $3, COALESCE($4, NOW()))
		`

	if _, err := tx.ExecContext(ctx, ins, a.ViewerId, a.ViewedId, a.Kind, atArg); err != nil {
		return nil, err
	}

	if kind != "like" {
		if err := tx.Commit(); err != nil {
			return nil, err
		}
		return nil, nil
	}
	
	const checkReciprocal = `
		SELECT 1
		FROM actions
		WHERE viewer_id = $1 AND viewed_id = $2 AND kind = 'like'
		LIMIT 1
		`
	var dummy int
	err = tx.QueryRowContext(ctx, checkReciprocal, a.ViewedId, a.ViewerId).Scan(&dummy)
	if err == sql.ErrNoRows {
		if err := tx.Commit(); err != nil {
			return nil, err
		}
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	u1, u2:= orderedPair(a.ViewerId, a.ViewedId)

	const insMatch = `
		INSERT INTO MATCHES (u1, u2, at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (u1, u2) DO NOTHING
		RETURNING u1, u2, at
		`

		var m matches.Match
	if err := tx.QueryRowContext(ctx, insMatch, u1, u2).Scan(&m.U1, &m.U2, &m.At); err != nil {
			if err == sql.ErrNoRows {
				const selMatch = `SELECT u1, u2, at FROM MATCHES WHERE u1 = $1 AND u2 = $2`
				if err := tx.QueryRowContext(ctx, selMatch, u1, u2).Scan(&m.U1, &m.U2, &m.At); err != nil {
					return nil, err
				}
			}else{
			return nil, err
		}

	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &m, nil
}


func orderedPair(a, b string) (string, string) {
	if a <= b {
		return a, b
	}
	return b, a
}
