package users

import (
	"context"
	"database/sql"
	"errors"
	"strings"
)

type Repository interface  {
	SetActive(ctx context.Context, userId string, active bool) (int64, error)
	Create(ctx context.Context, u User) error
	List(ctx context.Context) ([]User, error)
	ListByID(ctx context.Context, userId string) ([]User, error)

}

type repo struct {db *sql.DB}

func NewRepository(db *sql.DB) Repository {
	return &repo{db: db}
}

func (r *repo) 	SetActive(ctx context.Context, userId string, active bool) (int64, error) {
	userId = strings.TrimSpace(userId)
	if userId == "" {
		return 0, errors.New("user id required")
	}

	const q = `UPDATE users SET active = $2 WHERE id = $1`
	res, err := r.db.ExecContext(ctx, q, userId, active)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()

}

func (r *repo) Create(ctx context.Context, u User) error {
	u.ID = strings.TrimSpace(u.ID)
	if u.ID == ""  {
		return errors.New("id is required")
	}

	const q = `
			INSERT INTO users (id, x, y, active, distance)
			VALUES ($1, $2, $3, FALSE, $4)
		`
	_, err := r.db.ExecContext(ctx, q, u.ID, u.X, u.Y, u.Distance)
	return err
}

func (r * repo)	List(ctx context.Context) ([]User, error){
	const q = `SELECT id, x, y, active, distance FROM users ORDER by id`
	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.X, &u.Y, &u.Active, &u.Distance); err != nil {
			return nil, err
		}
		out = append(out, u)
	}

	return out, rows.Err()
}


func (r *repo)ListByID(ctx context.Context, userId string) ([]User, error){
	userId = strings.TrimSpace(userId)
	if userId == "" {
		return nil, errors.New("user_id required")
	}

	const q = `SELECT id, x, y, active, distance FROM users WHERE id = $1`
	rows, err := r.db.QueryContext(ctx, q, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var out []User
	for rows.Next() {
		var u User
		if err:= rows.Scan(&u.ID, &u.X, &u.Y, &u.Active, &u.Distance); err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, rows.Err()
}
