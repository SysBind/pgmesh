package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/sysbind/pgmesh/postgres"
	pgutil "github.com/sysbind/pgmesh/postgres/util"
)

func init() {
	rootCmd.AddCommand(pubsubCmd)
}

var pubsubCmd = &cobra.Command{
	Use:   "pubsub",
	Short: "Create logical replication between two databases",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting PubSub")

		stdout, stderr, _ := pgutil.DumpGlobals(postgres.ConnConfig{
			Host:     "localhost",
			Database: "moodle",
			User:     "postgres",
			Pass:     "q1w2e3r4"})

		for line := range stdout {
			fmt.Println(line)
		}

		for line := range stderr {
			log.Fatal(line)
		}
	},
}
