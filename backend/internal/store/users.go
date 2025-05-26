package store

import (
	"context"
	"database/sql"
)

type UsersStore struct {
	client *sql.DB
}
type PostgreUser struct {
	Username  string `json:"username"`
	Password  string `json:"-"`
	Email     string `json:"email"`
	UserID    string `json:"user_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
type password struct {
	text *string
	hash []byte
}

func (s *UsersStore) GetUserById(ctx context.Context, userId string) (*PostgreUser, error) {
	query := `
	SELECT id, username, email , password , created_at * FROM users WHERE id= $1 
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	user := &PostgreUser{}
	err := s.client.QueryRowContext(
		ctx,
		query,
		userId,
	).Scan(
		&user.Name,
		&user.Password,
		&user.Email,
		&user.UserID,
		&user.CreatedAt,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}

func (s *UsersStore) CreateUser(ctx context.Context, user *PostgreUser) error {
	query := `
	INSERT INTO users (name, password , email) VALUES($1,$2,$3) RETURNING id, created_at
	`

	err := s.client.QueryRowContext(
		ctx,
		query,
		user.Name,
		user.Password.hash,
		user.Email,
	).Scan(
		&user.UserID,
		&user.CreatedAt,
	)

	if err != nil {
		return err
	}
	return nil

}
