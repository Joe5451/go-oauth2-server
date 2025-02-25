package database

import (
	"context"
	"fmt"

	"github.com/Joe5451/go-oauth2-server/internal/config"
	"github.com/pkg/errors"

	"github.com/jackc/pgx/v5"
)

func NewPostgresDB() (*pgx.Conn, error) {
	databaseUrl := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		config.AppConfig.DBUser,
		config.AppConfig.DBPassword,
		config.AppConfig.DBHost,
		config.AppConfig.DBPort,
		config.AppConfig.DBName,
	)

	conn, err := pgx.Connect(context.Background(), databaseUrl)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to connect to database")
	}

	return conn, nil
}
