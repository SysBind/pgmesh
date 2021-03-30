package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(pubsubCmd)
}

var pubsubCmd = &cobra.Command{
	Use:   "pubsub",
	Short: "Create logical replication between two databases",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting PubSub")
	},
}
