package postgres

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

type ConnConfig struct {
	Host     string
	Database string
	User     string
	Pass     string
}

func Connect(ctx context.Context, cfg ConnConfig) (*pgxpool.Pool, error) {
	connstr := fmt.Sprintf("postgres://%s:%s@%s:5432/%s",
		cfg.User, cfg.Pass, cfg.Host, cfg.Database)
	poolconf, err := pgxpool.ParseConfig(connstr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed configuring connection pool: %v\n", err)
		return nil, err
	}
	conn, err := pgxpool.ConnectConfig(ctx, poolconf)

	return conn, err
}
