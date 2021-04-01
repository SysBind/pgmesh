package logical

import (
	"context"
	"errors"

	"github.com/sysbind/pgmesh/postgres"
)

var (
	errWalLevelNotLogical error = errors.New("wal level should be logical on source cluster")
)

func PreFlight(ctx context.Context, src, dest postgres.ConnConfig) error {
	src_db, err := postgres.Connect(ctx, src)

	if err != nil {
		return err
	}
	var wal_level string
	if err = src_db.QueryRow(ctx, "SHOW wal_level").Scan(&wal_level); err != nil {
		return err
	}
	if wal_level != "logical" {
		return errWalLevelNotLogical
	}
	return nil
}
