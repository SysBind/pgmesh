package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/sysbind/pgmesh/postgres"
	pgcopy "github.com/sysbind/pgmesh/postgres/copy"
)

var slack int // number to add to all copied sequence values

func init() {
	copyseqCmd.Flags().StringVarP(&srcHost, "source-host", "", "", "Source database host (required)")
	copyseqCmd.MarkFlagRequired("source-host")
	copyseqCmd.Flags().StringVarP(&srcDB, "source-database", "", "", "Source database name (required)")
	copyseqCmd.MarkFlagRequired("source-database")
	copyseqCmd.Flags().StringVarP(&destHost, "dest-host", "", "", "Destination database host (required)")
	copyseqCmd.MarkFlagRequired("dest-host")
	copyseqCmd.Flags().StringVarP(&destDB, "dest-database", "", "", "Destination database name (required)")
	copyseqCmd.MarkFlagRequired("dest-database")
	copyseqCmd.Flags().IntVarP(&slack, "slack", "", 0, "a number to add to each sequence value copied (to avoid conflict if there are still potential, non-critical writes to source db, like in maintnence mode of some apps)")
	rootCmd.AddCommand(copyseqCmd)
}

var copyseqCmd = &cobra.Command{
	Use:   "copyseq",
	Short: "Copies all sequence values from source to dest with optional slack",
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

		fmt.Println("Calling CopySequences..")
		if err := pgcopy.CopySequences(ctx, src, dest, slack); err != nil {
			log.Fatal(err)
		}
	},
}
