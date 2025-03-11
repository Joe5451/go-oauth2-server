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
	// SQL query with left join to fetch user and associated social accounts
	query := `
        SELECT u.id, u.email, u.name, u.avatar, s.id AS social_account_id, s.provider, s.provider_user_id, s.email, s.name, s.avatar FROM users u
        LEFT JOIN social_accounts s ON u.id = s.user_id
        WHERE u.email = @email
	`

	args := pgx.NamedArgs{
		"email": email,
	}

	rows, err := r.conn.Query(context.Background(), query, args)
	if err != nil {
		return domain.User{}, err
	}
	defer rows.Close()

	var user domain.User
	var socialAccounts []domain.SocialAccount

	for rows.Next() {
		var account domain.SocialAccount

		err := rows.Scan(
			&user.ID, &user.Email, &user.Name, &user.Avatar,
			&account.ID, &account.Provider, &account.ProviderUserID, &account.Email, &account.Name, &account.Avatar,
		)
		if err != nil {
			return domain.User{}, fmt.Errorf("Error Fetching user and social account: %w", err)
		}
		socialAccounts = append(socialAccounts, account)
	}

	if user.ID == 0 {
		return domain.User{}, domain.ErrUserNotFound
	}

	user.SocialAccounts = socialAccounts
	return user, nil
}

func (r *PostgresUserRepository) UpdateOrCreateSocialAccount(socialAccount domain.SocialAccount) (domain.SocialAccount, error) {
	query := `
		INSERT INTO social_accounts (user_id, provider, provider_user_id, email, name, avatar)
		VALUES (@user_id, @provider, @provider_user_id, @email, @name, @avatar)
		ON CONFLICT (provider, provider_user_id)
		DO UPDATE SET
			email = EXCLUDED.email,
			name = EXCLUDED.name,
			avatar = EXCLUDED.avatar,
			updated_at = CURRENT_TIMESTAMP
		RETURNING id, user_id, provider, provider_user_id, email, name, avatar, created_at, updated_at
	`

	args := pgx.NamedArgs{
		"user_id":          socialAccount.UserID,
		"provider":         socialAccount.Provider,
		"provider_user_id": socialAccount.ProviderUserID,
		"email":            socialAccount.Email,
		"name":             socialAccount.Name,
		"avatar":           socialAccount.Avatar,
	}

	err := r.conn.QueryRow(context.Background(), query, args).Scan(
		&socialAccount.ID,
		&socialAccount.UserID,
		&socialAccount.Provider,
		&socialAccount.ProviderUserID,
		&socialAccount.Email,
		&socialAccount.Name,
		&socialAccount.Avatar,
		&socialAccount.CreatedAt,
		&socialAccount.UpdatedAt,
	)

	if err != nil {
		return domain.SocialAccount{}, fmt.Errorf("failed to insert or update social account: %w", err)
	}

	return socialAccount, nil
}

func (r *PostgresUserRepository) GetSocialAccountByProviderUserID(providerUserID string) (domain.SocialAccount, error) {
	query := `
		SELECT id, provider, provider_user_id, user_id, created_at, updated_at FROM social_accounts WHERE provider_user_id = @provider_user_id
	`

	args := pgx.NamedArgs{
		"provider_user_id": providerUserID,
	}

	var socialAccount domain.SocialAccount

	err := r.conn.QueryRow(context.Background(), query, args).Scan(
		&socialAccount.ID,
		&socialAccount.Provider,
		&socialAccount.ProviderUserID,
		&socialAccount.UserID,
		&socialAccount.CreatedAt,
		&socialAccount.UpdatedAt,
	)
	if err != nil {
		return domain.SocialAccount{}, fmt.Errorf("Error fetching social account: %w", err)
	}

	return socialAccount, nil
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
