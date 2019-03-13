package dlccli

import (
	"github.com/spf13/cobra"
)

var contractsCmd = &cobra.Command{
	Use:   "contracts",
	Short: "Contract commands",
}

func init() {
	// subcommand root
	rootCmd.AddCommand(contractsCmd)

	// create contract
	contractsCmd.AddCommand(initCreateContractCmd())

	// view contract
	// contractsCmd.AddCommand(initViewContractCmd())
}
