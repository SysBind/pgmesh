package logical

import (
	"context"
	"fmt"

	"github.com/sysbind/pgmesh/postgres"
)

// PubSub creates publication on source database
//        and subcribes to it on dest database
//        to initiate logical replication on all
//        tables
func PubSub(ctx context.Context, src, dest postgres.ConnConfig) error {
	src_db, err := postgres.Connect(ctx, src)
	if err != nil {
		return err
	}

	tag, err := src_db.Exec(ctx, "CREATE PUBLICATION p_upgrade FOR ALL TABLES")
	if err != nil {
		return err
	}
	fmt.Println(tag)

	target_db, err := postgres.Connect(ctx, dest)
	if err != nil {
		return err
	}

	subscribe_cmd := fmt.Sprintf(
		"CREATE SUBSCRIPTION s_upgrade CONNECTION 'host=%s port=5432 dbname=%s user=%s password=%s' PUBLICATION p_upgrade;",
		src.Host, src.Database, src.User, src.Pass)
	tag, err = target_db.Exec(ctx, subscribe_cmd)
	if err != nil {
		return err
	}
	fmt.Println(tag)

	return nil
}

// Detach tears down replication setup done by the PubSub func
func Detach(ctx context.Context, dest postgres.ConnConfig) error {
	target_db, err := postgres.Connect(ctx, dest)
	if err != nil {
		return err
	}
	tag, err := target_db.Exec(ctx, "DROP SUBSCRIPTION s_upgrade;")
	if err != nil {
		return err
	}
	fmt.Println(tag)

	return nil
}
