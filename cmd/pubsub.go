package cmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	"github.com/sysbind/pgmesh/postgres"
	pgcopy "github.com/sysbind/pgmesh/postgres/copy"
	pglogical "github.com/sysbind/pgmesh/postgres/logical"
)

func init() {
	rootCmd.AddCommand(pubsubCmd)
}

var pubsubCmd = &cobra.Command{
	Use:   "pubsub",
	Short: "Create logical replication between two databases",
	Run: func(cmd *cobra.Command, args []string) {
		src := postgres.ConnConfig{
			Host:     "localhost",
			Database: "moodle",
			User:     "postgres",
			Pass:     "q1w2e3r4"}

		dest := postgres.ConnConfig{
			Host:     "localhost",
			Database: "moodle2",
			User:     "postgres",
			Pass:     "q1w2e3r4"}

		ctx := context.Background()
		if err := pgcopy.CopySchema(ctx, src, dest); err != nil {
			log.Fatal(err)
		}
		if err := pgcopy.CopyPrimeKeys(ctx, src, dest); err != nil {
			log.Fatal(err)
		}
		if err := pglogical.PubSub(ctx, src, dest); err != nil {
			log.Fatal(err)
		}
	},
}
