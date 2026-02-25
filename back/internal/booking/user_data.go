package booking

import (
	"context"
	"database/sql"
	"search-job/internal/models"
)

// CreateUser — сохраняет нового пользователя с зашифрованным паролем
func (r *Repo) CreateUser(ctx context.Context, u *models.User) error {
	query := `
		INSERT INTO users (email, password, created_at)
		VALUES ($1, $2, NOW())
		RETURNING id, created_at`

	return r.db.QueryRow(ctx, query, u.Email, u.Password).Scan(&u.ID, &u.CreatedAt)
}

// GetUserByEmail — ищет пользователя по email для логина
func (r *Repo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT id, email, password, created_at FROM users WHERE email = $1`

	var u models.User
	err := r.db.QueryRow(ctx, query, email).Scan(&u.ID, &u.Email, &u.Password, &u.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Пользователь не найден
		}
		return nil, err
	}
	return &u, nil
}
