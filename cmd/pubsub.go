package cmd

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

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

		dump := pgutil.DumpGlobals(postgres.ConnConfig{
			Host:     "localhost",
			Database: "moodle",
			User:     "postgres",
			Pass:     "q1w2e3r4"})

		stdout, _ := dump.StdoutPipe()
		stderr, _ := dump.StderrPipe()
		err := dump.Start()
		if err != nil {
			log.Fatal(err)
		}

		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			fmt.Println(scanner.Text()) // Println will add back the final '\n'
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
		}

		stderr_str, err := io.ReadAll(stderr)
		if err != nil {
			log.Fatal(err)
		}

		if err := dump.Wait(); err != nil {
			fmt.Printf("%s", stderr_str)
			log.Fatal(err)
		}
	},
}
