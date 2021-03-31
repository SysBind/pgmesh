package cmd

import (
	"context"
	"fmt"
	"log"
	"strings"

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

		ctx := context.Background()

		stdout, stderr, err := pgutil.DumpSchema(postgres.ConnConfig{
			Host:     "localhost",
			Database: "moodle",
			User:     "postgres",
			Pass:     "q1w2e3r4"}, pgutil.PreData)
		if err != nil {
			log.Fatal(err)
		}
		target_db, err := postgres.Connect(ctx, postgres.ConnConfig{
			Host:     "localhost",
			Database: "moodle2",
			User:     "postgres",
			Pass:     "q1w2e3r4"})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Target DB contacted at %s/%s\n", "localhost", "moodle2")
		var statement string = ""
		for line := range stdout {
			fmt.Printf("Scanning %s\n", line)
			if line == "" || strings.HasPrefix(line, "--") {
				continue
			}
			statement += line
			if statement[len(statement)-1] == ';' {
				fmt.Printf("Running %s\n", statement)
				tag, err := target_db.Exec(ctx, statement)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println(string(tag))
				statement = ""
			}
		}

		for line := range stderr {
			log.Fatal(line)
		}
	},
}
