package copy

import (
	"context"
	"fmt"

	"github.com/sysbind/pgmesh/postgres"
)

// CopySequences copies all sequence values from source to target (as before migrating to new)
func CopySequences(ctx context.Context, src, dest postgres.ConnConfig, slack int) error {
	src_db, err := postgres.Connect(ctx, src)
	if err != nil {
		return err
	}

	dest_db, err := postgres.Connect(ctx, dest)
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

		_, err := dest_db.Exec(ctx, fmt.Sprintf("SELECT setval ('%s.%s', %d)",
			schema, sequence, lastval+slack))
		if err != nil {
			return err
		}
	}
	// Check for errors from iterating over rows.
	err = rows.Err()

	return err
}
