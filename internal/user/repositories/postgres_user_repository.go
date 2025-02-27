package repositories

import (
	"context"
	"log"

	"github.com/Joe5451/go-oauth2-server/internal/domains"
	"github.com/jackc/pgx/v5"
)

type PostgresUserRepository struct {
	conn *pgx.Conn
}

func (r *PostgresUserRepository) Create(user domains.User) (domains.User, error) {
	query := `
        INSERT INTO users (
			email,
			password,
			phone_number,
			username,
			gender,
			avatar,
		) VALUES (
			@email,
			@password,
			@phone_number,
			@username,
			@gender,
			@avatar,
		)
	`

	args := pgx.NamedArgs{
		"email":        user.Email,
		"password":     user.Password,
		"phone_number": user.PhoneNumber,
		"username":     user.Username,
		"gender":       user.Gender,
		"avatar":       user.Avatar,
	}

	_, err := r.conn.Exec(context.Background(), query, args)
	if err != nil {
		log.Println("Error Inserting User")
		return domains.User{}, err
	}

	// WIP: Will return populated user data after implementation
	return domains.User{}, nil
}
