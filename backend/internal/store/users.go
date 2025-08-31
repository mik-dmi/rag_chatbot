package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UsersStore struct {
	client *sql.DB
}
type PostgreUser struct {
	UserID    string   `json:"user_id"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Password  password `json:"-"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
	ISActive  bool     `json:"is_active"`
}
type password struct {
	text *string
	hash []byte
}

func (p *password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.text = &text
	p.hash = hash

	return nil
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
		&user.Password.hash,
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

func (s *UsersStore) CreateUser(ctx context.Context, tx *sql.Tx, user *PostgreUser) error {
	const query = `
    INSERT INTO users (username, password, email)
    VALUES ($1, $2, $3)
    RETURNING user_id, created_at, updated_at;
    `

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.client.
		QueryRowContext(ctx, query,
			user.Username,
			user.Password.hash,
			user.Email,
		).
		Scan(
			&user.UserID,
			&user.CreatedAt,
			&user.UpdatedAt,
		)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "user_email_key"`:
			return ErrDuplicateEmail
		case err.Error() == `pq: duplicate key value violates unique constraint "user_username_key"`:
			return ErrDuplicateUsername
		default:
			return err
		}
	}
	if err != nil {
		return err
	}
	return nil
}

func (s *UsersStore) CreateAndInvite(ctx context.Context, user *PostgreUser, token string, invitationExp time.Duration) error {
	return withTx(s.client, ctx, func(tx *sql.Tx) error {

		if err := s.CreateUser(ctx, tx, user); err != nil {
			return err
		}

		if err := s.createUserInvitation(ctx, tx, token, invitationExp, user.UserID); err != nil {
			return err
		}
		return nil

	})

}

func (s *UsersStore) createUserInvitation(ctx context.Context, tx *sql.Tx, token string, exp time.Duration, userID string) error {

	query := `INSERT INTO user_invitations (token, user_id, expiry ) VALUES ($1, $2, $3)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, token, userID, time.Now().Add(exp))
	if err != nil {
		return err
	}

	return nil
}

func (s *UsersStore) Activate(ctx context.Context, token string) error {

	return withTx(s.client, ctx, func(tx *sql.Tx) error {
		// 1. finsd the user the token belongs to
		user, err := s.getUserFromInvitation(ctx, tx, token)
		if err != err {
			return err

		}
		// 2. update the user
		user.ISActive = true

		if err := s.update(ctx, tx, user); err != nil {
			return err
		}
		// clean the invitations

		if err := s.deleteUserInvitations(ctx, tx, user.UserID); err != nil {
			return nil

		}

		return nil
	})
}

func (s *UsersStore) deleteUserInvitations(ctx context.Context, tx *sql.Tx, userID string) error {

	query := `DELETE FROM user WHERE user_id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}
	return nil
}

func (s *UsersStore) getUserFromInvitation(ctx context.Context, tx *sql.Tx, plainToken string) (*PostgreUser, error) {
	query := `
		SELECT u.user_id, u.username, u.email, u.created_at, u.updated_at, u.is_active
		FROM users u
		JOIN user_invitations ui ON u.user_id = ui.user_id 
		WHERE ui.token 	=$1 AND ui.expiry > $2  
		
		`
	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	user := &PostgreUser{}
	err := tx.QueryRowContext(ctx, query, hashToken, time.Now()).Scan(
		&user.UserID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.ISActive,
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

func (s *UsersStore) update(ctx context.Context, tx *sql.Tx, user *PostgreUser) error {

	query := `UPDATE users SET username = $1 , email = $2 , is_active = $3 WHERE id = $4 `

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, user.Username, user.Email, user.ISActive, user.UserID)
	if err != nil {
		return err
	}
	return nil

}

func (s *UsersStore) delete(ctx context.Context, tx *sql.Tx, userID string) error {

	query := ` DELETE FROM users WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}
	return nil

}
func (s *UsersStore) Delete(ctx context.Context, userID string) error {

	return withTx(s.client, ctx, func(tx *sql.Tx) error {
		if err := s.delete(ctx, tx, userID); err != nil {
			return err
		}

		if err := s.deleteUserInvitations(ctx, tx, userID); err != nil {
			return err
		}
		return nil
	})

}
