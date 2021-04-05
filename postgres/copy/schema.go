package copy

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/sysbind/pgmesh/postgres"
	pgutil "github.com/sysbind/pgmesh/postgres/util"
)

// CopySchema copies basic schema from src to dest (pg_dump --section=pre-data)
func CopySchema(ctx context.Context, src, dest postgres.ConnConfig) error {
	stdout, stderr, err := pgutil.DumpSchema(src, pgutil.PreData)
	if err != nil {
		return err
	}
	target_db, err := postgres.Connect(ctx, dest)
	if err != nil {
		return err
	}
	fmt.Printf("CopySchema: Target DB contacted at %s/%s\n",
		dest.Host,
		dest.Database)
	var statement string = ""
	for line := range stdout {
		if line == "" || strings.HasPrefix(line, "--") {
			continue
		}
		statement += line
		if statement[len(statement)-1] == ';' {
			_, err := target_db.Exec(ctx, statement)
			if err != nil {
				log.Fatal(err)
			}
			statement = ""
		}
	}

	for line := range stderr {
		if line == "" {
			continue
		}
		return errors.New(line)
	}

	return nil
}

// CopyPrimeKeys copies primary-keys from src to dest
// (pg_dump --section=post-data, than filter only primary keys)
// This is required for initating the logical replication properly.
// Other contraints and indexes are ignored in this stage to speed-up
// the synchronization phase.
func CopyPrimeKeys(ctx context.Context, src, dest postgres.ConnConfig) error {
	stdout, stderr, err := pgutil.DumpSchema(src, pgutil.PostData)
	if err != nil {
		return err
	}
	target_db, err := postgres.Connect(ctx, dest)
	if err != nil {
		return err
	}
	fmt.Printf("CopyPrimeKeys: Target DB contacted at %s/%s\n",
		dest.Host,
		dest.Database)
	var statement string = ""
	for line := range stdout {
		if line == "" || strings.HasPrefix(line, "--") {
			continue
		}
		statement += line
		if statement[len(statement)-1] == ';' {
			if !strings.Contains(statement, "PRIMARY KEY") {
				statement = ""
				continue
			}
			_, err := target_db.Exec(ctx, statement)
			if err != nil {
				log.Fatal(err)
			}
			statement = ""
		}
	}

	for line := range stderr {
		if line == "" {
			continue
		}
		return errors.New(line)
	}

	return nil
}

// CopyContraints copies all other constraints (non primary-keys)
// (pg_dump --section=post-data, than filter out primary keys)
func CopyConstraints(ctx context.Context, src, dest postgres.ConnConfig) error {
	stdout, stderr, err := pgutil.DumpSchema(src, pgutil.PostData)
	if err != nil {
		return err
	}
	target_db, err := postgres.Connect(ctx, dest)
	if err != nil {
		return err
	}
	fmt.Printf("CopyConstraints: Target DB contacted at %s/%s\n",
		dest.Host,
		dest.Database)
	var statement string = ""
	for line := range stdout {
		if line == "" || strings.HasPrefix(line, "--") {
			continue
		}
		statement += line
		if statement[len(statement)-1] == ';' {
			if strings.Contains(statement, "PRIMARY KEY") {
				statement = ""
				continue
			}
			_, err := target_db.Exec(ctx, statement)
			if err != nil {
				log.Fatal(err)
			}
			statement = ""
		}
	}

	for line := range stderr {
		if line == "" {
			continue
		}
		return errors.New(line)
	}

	return nil
}
