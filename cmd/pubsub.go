package cmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	"github.com/sysbind/pgmesh/postgres"
	pgcopy "github.com/sysbind/pgmesh/postgres/copy"
	pglogical "github.com/sysbind/pgmesh/postgres/logical"
)

var (
	srcHost  string
	srcDB    string
	destHost string
	destDB   string
)

func init() {
	pubsubCmd.Flags().StringVarP(&srcHost, "source-host", "", "", "Source database host (required)")
	pubsubCmd.MarkFlagRequired("source-host")
	pubsubCmd.Flags().StringVarP(&srcDB, "source-database", "", "", "Source database name (required)")
	pubsubCmd.MarkFlagRequired("source-database")
	pubsubCmd.Flags().StringVarP(&destHost, "dest-host", "", "", "Destination database host (required)")
	pubsubCmd.MarkFlagRequired("dest-host")
	pubsubCmd.Flags().StringVarP(&destDB, "dest-database", "", "", "Destination database name (required)")
	pubsubCmd.MarkFlagRequired("dest-database")
	rootCmd.AddCommand(pubsubCmd)
}

var pubsubCmd = &cobra.Command{
	Use:   "pubsub",
	Short: "Create logical replication between two databases",
	Run: func(cmd *cobra.Command, args []string) {
		src := postgres.ConnConfig{
			Host:     srcHost,
			Database: srcDB,
			User:     "postgres",
			Pass:     "q1w2e3r4"}

		dest := postgres.ConnConfig{
			Host:     destHost,
			Database: destDB,
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
		if err := pgcopy.CopyConstraints(ctx, src, dest); err != nil {
			log.Fatal(err)
		}
	},
}
