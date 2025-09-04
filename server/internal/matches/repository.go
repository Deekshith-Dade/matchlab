package matches

import (
	"context"
	"database/sql"
	"errors"
	"strings"
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
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil, errors.New("userID required")
	}

	const q = `
		SELECT u1, u2, at
		FROM matches
		WHERE u1 = $1 OR u2 = $1
		ORDER BY at DESC;
	`

	rows, err := r.db.QueryContext(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Match
	for rows.Next() {
		var m Match
		if err := rows.Scan(&m.U1, &m.U2, &m.At); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

func (r *repo) Create(ctx context.Context, m Match) error {
	m.U1 = strings.TrimSpace(m.U1)
	m.U2 = strings.TrimSpace(m.U2)
	if m.U1 == "" || m.U2 == "" {
		return errors.New("both user IDs required")
	}
	if m.U1 == m.U2 {
		return errors.New("cannot match a user with itself")
	}

	// Enforce consistent ordering so (u1,u2) and (u2,u1) are treated the same
	if m.U1 > m.U2 {
		m.U1, m.U2 = m.U2, m.U1
	}

	var atArg any
	if m.At.IsZero() {
		atArg = nil // let DB default to NOW()
	} else {
		atArg = m.At.UTC()
	}

	const ins = `
		INSERT INTO matches (u1, u2, at)
		VALUES ($1, $2, COALESCE($3, NOW()))
		ON CONFLICT (u1, u2) DO NOTHING
	`
	_, err := r.db.ExecContext(ctx, ins, m.U1, m.U2, atArg)
	return err
}
