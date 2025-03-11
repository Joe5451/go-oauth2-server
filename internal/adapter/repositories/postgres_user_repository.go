package repositories

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/Joe5451/go-oauth2-server/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type PostgresUserRepository struct {
	conn *pgx.Conn
}

func NewPostgresUserRepository(conn *pgx.Conn) *PostgresUserRepository {
	return &PostgresUserRepository{
		conn: conn,
	}
}

func (r *PostgresUserRepository) CreateUser(user domain.User) (domain.User, error) {
	query := `
        INSERT INTO users (email, password, name, avatar) VALUES (@email, @password, @name, @avatar) RETURNING id, email, name, avatar
	`

	args := pgx.NamedArgs{
		"email":    user.Email,
		"password": user.Password,
		"name":     user.Name,
		"avatar":   user.Avatar,
	}

	err := r.conn.QueryRow(context.Background(), query, args).Scan(&user.ID, &user.Email, &user.Name, &user.Avatar)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.User{}, domain.ErrDuplicateEmail
		}
		return domain.User{}, fmt.Errorf("database error: %w", err)
	}

	return user, nil
}

func (r *PostgresUserRepository) GetUser(userID int64) (domain.User, error) {
	query := `
		SELECT id, email, password, name, avatar, created_at, updated_at FROM users WHERE id = @id
	`

	args := pgx.NamedArgs{
		"id": userID,
	}

	var user domain.User
	err := r.conn.QueryRow(context.Background(), query, args).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.Avatar,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, err
	}

	return user, nil
}

func (r *PostgresUserRepository) GetUserByEmail(email string) (domain.User, error) {
	query := `
		SELECT
			id,
			email,
			phone_number,
			username,
			gender,
			avatar,
			created_at,
			updated_at,
		FROM users WHERE email = @email
	`

	args := pgx.NamedArgs{
		"email": email,
	}

	var user domain.User
	err := r.conn.QueryRow(context.Background(), query, args).Scan(
		&user.ID,
		&user.Email,
		&user.PhoneNumber,
		&user.Username,
		&user.Gender,
		&user.Avatar,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (r *PostgresUserRepository) FirstOrCreateSocialAccount(
	provider, providerUserID string,
) (domain.SocialAccount, error) {
	// pending implementation
	return domain.SocialAccount{}, nil
}

func (r *PostgresUserRepository) CreateSocialAccount(account domain.SocialAccount) domain.SocialAccount {
	return domain.SocialAccount{}
}

func (r *PostgresUserRepository) UpdateSocialAccount(account domain.SocialAccount) error {
	query := `
        UPDATE social_accounts SET user_id = @user_id, updated_at = CURRENT_TIMESTAMP WHERE id = @social_account_id
    `

	args := pgx.NamedArgs{
		"user_id":           account.UserID,
		"social_account_id": account.ID,
	}

	cmdTag, err := r.conn.Exec(context.Background(), query, args)
	if err != nil {
		log.Printf("Error updating social account: %v\n", err)
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("no social account updated (account: %v)", account)
	}

	return nil
}

func (r *PostgresUserRepository) UpdateUser(usreID int64, user domain.User) error {
	return nil
}
