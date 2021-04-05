package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	srcHost  string
	srcDB    string
	destHost string
	destDB   string
	slack    int // number to add to all copied sequence values
)

var rootCmd = &cobra.Command{
	Use:   "pgmesh",
	Short: "pgmesh automates common postgres logical replication use-cases",
	Long: `pgmesh automates common postgres logical replication use-cases -
                Long version
                Complete documentation is available at https://gitub.com/sysbind/pgmesh`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("pgmesh: no command given (try pgmesh --help)")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
