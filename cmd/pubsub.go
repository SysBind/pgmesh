package cmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	"github.com/sysbind/pgmesh/postgres"
	pgcopy "github.com/sysbind/pgmesh/postgres/copy"
	pglogical "github.com/sysbind/pgmesh/postgres/logical"
)

var detach bool

func init() {
	pubsubCmd.Flags().StringVarP(&srcHost, "source-host", "", "", "Source database host (required)")
	pubsubCmd.MarkFlagRequired("source-host")
	pubsubCmd.Flags().StringVarP(&srcDB, "source-database", "", "", "Source database name (required)")
	pubsubCmd.MarkFlagRequired("source-database")
	pubsubCmd.Flags().StringVarP(&srcUser, "source-user", "", "postgres", "Source database user")
	pubsubCmd.Flags().StringVarP(&srcPass, "source-pass", "", "", "Source database password")
	pubsubCmd.Flags().StringVarP(&destHost, "dest-host", "", "", "Destination database host (required)")
	pubsubCmd.MarkFlagRequired("dest-host")
	pubsubCmd.Flags().StringVarP(&destDB, "dest-database", "", "", "Destination database name (required)")
	pubsubCmd.MarkFlagRequired("dest-database")
	pubsubCmd.Flags().BoolVarP(&detach, "detach", "", false, "Detach previously etstablished logical replication")
	pubsubCmd.Flags().StringVarP(&destUser, "dest-user", "", "postgres", "Destination database user")
	pubsubCmd.Flags().StringVarP(&destPass, "dest-pass", "", "", "Destination database password")

	rootCmd.AddCommand(pubsubCmd)
}

var pubsubCmd = &cobra.Command{
	Use:   "pubsub",
	Short: "Create / Detach logical replication between two databases",
	Run: func(cmd *cobra.Command, args []string) {
		src := postgres.ConnConfig{
			Host:     srcHost,
			Database: srcDB,
			User:     srcUser,
			Pass:     srcPass}

		dest := postgres.ConnConfig{
			Host:     destHost,
			Database: destDB,
			User:     destUser,
			Pass:     destPass}

		ctx := context.Background()
		if detach {
			if err := pglogical.Detach(ctx, dest); err != nil {
				log.Fatal(err)
			}
			return
		}
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
