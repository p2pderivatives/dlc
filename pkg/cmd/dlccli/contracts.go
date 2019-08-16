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

	// create contract with premium command
	contractsCmd.AddCommand(initCreateContractWithPremiumCmd())

	// subcommand deals
	contractsCmd.AddCommand(dealsCmd)

	// fix deal
	dealsCmd.AddCommand(initFixDealCmd())
}
