package dlccli

import (
	"github.com/spf13/cobra"
)

var contractsCmd = &cobra.Command{
	Use:   "contracts",
	Short: "Contract commands",
}

var dealsCmd = &cobra.Command{
	Use:   "deals",
	Short: "Deals commands",
}

func init() {
	// subcommand root
	rootCmd.AddCommand(contractsCmd)

	// create contract
	contractsCmd.AddCommand(initCreateContractCmd())

	// subcommand deals
	contractsCmd.AddCommand(dealsCmd)

	// fix deal
	dealsCmd.AddCommand(initFixDealCmd())
}
