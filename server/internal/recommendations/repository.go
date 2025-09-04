package recommendations

import (
	"context"
	"database/sql"
	"errors"
	"strings"
)

type Repository interface {
	getRecommendations(ctx context.Context, userId string, topK int) ([]Recommendation, error)
}


type repo struct {db *sql.DB}

func NewRepository(db *sql.DB) Repository {
	return &repo{db: db}
}


func (r *repo) getRecommendations(ctx context.Context, userId string, topK int)( []Recommendation, error) {
		userId = strings.TrimSpace(userId)
	if userId == "" || topK < 0 {
		return nil, errors.New("Invalid input") 
	}
	
	selUsers := `
	SELECT id,
       ROW_NUMBER() OVER (ORDER BY id) AS rank
	FROM users
	WHERE id != $1
		AND active = true
	LIMIT $2
	`
	rows, err := r.db.QueryContext(ctx, selUsers, userId, topK)
	if err != nil {
		return nil, err
	}

	var out []Recommendation
	for rows.Next() {
		var rec Recommendation
		if err := rows.Scan(&rec.UserID, &rec.Rank); err != nil {
			return nil, err
		}
		out = append(out, rec)
	}
	
	return out, nil

}
