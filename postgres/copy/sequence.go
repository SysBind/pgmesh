package copy

import (
	"context"

	"github.com/sysbind/pgmesh/postgres"
)

// CopySchema copies basic schema from src to dest (pg_dump --section=pre-data)
func CopySequences(ctx context.Context, src, dest postgres.ConnConfig) error {
	src_db, err := postgres.Connect(ctx, src)
	if err != nil {
		return err
	}

	rows, err := src_db.Query(ctx, "SELECT schemaname,sequencename,coalesce(last_value, 1) AS lastval FROM pg_sequences")
	if err != nil {
		return err
	}
	defer rows.Close()

	var schema, sequence string
	var lastval int
	for rows.Next() {
		if err = rows.Scan(&schema,
			&sequence,
			&lastval); err != nil {
			return err
		}
	}
	// Check for errors from iterating over rows.
	err = rows.Err()

	return err
}
