package storage

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

func Connect(ctx context.Context) (*pgx.Conn, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))

	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

type StorePostgres struct {
	db *pgx.Conn
}

func NewStore(db *pgx.Conn) *StorePostgres {
	return &StorePostgres{db: db}
}
