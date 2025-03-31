package store

import (
	"context"
	"database/sql"
)

type UsersStore struct {
	db *sql.DB
}
type PostgreUser struct {
	UserID    string `json:"userid"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (s *UsersStore) GetUserById(ctx context.Context, userId string) (*PostgreUser, error) {

	return nil, nil
}

func (s *UsersStore) CreateUser(ctx context.Context, user *PostgreUser) error {
	query := `
	INSERT INTO users (name, email) VALUES($1,$2,$3) RETURNING id, created_at
	`

	err := s.db.QueryRowContext(
		ctx,
		query,
		user.Name,
		user.Email,
	).Scan(&user.UserID, &user.CreatedAt)

	if err != nil {
		return err
	}
	return nil

}
