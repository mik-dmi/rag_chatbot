package store

import (
	"context"
	"database/sql"
)

type UsersStore struct {
	client *sql.DB
}
type PostgreUser struct {
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
type password struct {
	text *string
	hash []byte
}

func (s *UsersStore) GetUserById(ctx context.Context, userId string) (*PostgreUser, error) {
	query := `
	SELECT   id,
        username,
        email,
        password,
        created_at,
        updated_at * FROM users WHERE id= $1 
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	user := &PostgreUser{}
	err := s.client.QueryRowContext(
		ctx,
		query,
		userId,
	).Scan(
		&user.UserID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
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
	const query = `
    INSERT INTO users (username, password, email)
    VALUES ($1, $2, $3)
    RETURNING user_id, created_at, updated_at;
    `

	err := s.client.
		QueryRowContext(ctx, query,
			user.Username,
			user.Password,
			user.Email,
		).
		Scan(
			&user.UserID,
			&user.CreatedAt,
			&user.UpdatedAt,
		)

	if err != nil {
		return err
	}
	return nil

}
