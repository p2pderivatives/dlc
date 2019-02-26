package dlccli

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/spf13/cobra"

	"github.com/dgarage/dlc/internal/wallet"
)

// walletCmd represents the wallet command
var walletsCmd = &cobra.Command{
	Use:   "wallets",
	Short: "Operate wallets",
}

var walletsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new wallet",
	Run: func(cmd *cobra.Command, args []string) {
		netParams, err := loadNetParams(bitcoinConf)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// TODO: give seed as command line parameter
		seed, err := hdkeychain.GenerateSeed(hdkeychain.RecommendedSeedLen)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("Seed: %s\n", hex.EncodeToString(seed))

		wallet, err := wallet.CreateWallet(netParams,
			seed, []byte(pubpass), []byte(privpass),
			walletDir, walletName)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = wallet.Close()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Wallet created")
	},
}

func init() {
	rootCmd.AddCommand(walletsCmd)
	walletsCmd.AddCommand(walletsCreateCmd)
}
