package dlccli

import (
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/spf13/cobra"

	_wallet "github.com/dgarage/dlc/internal/wallet"
	"github.com/dgarage/dlc/pkg/wallet"
)

// var seed []byte
var pubpass string
var privpass string
var walletName string

// walletCmd represents the wallet command
var walletsCmd = &cobra.Command{
	Use:   "wallets",
	Short: "Wallet commands",
}

var walletsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new wallet",
	Run: func(cmd *cobra.Command, args []string) {
		netParams := loadNetParams(bitcoinConf)

		// TODO: give seed as command line parameter
		seed, err := hdkeychain.GenerateSeed(hdkeychain.RecommendedSeedLen)
		errorHandler(err)
		fmt.Printf("Seed: %s\n", hex.EncodeToString(seed))

		w, err := _wallet.CreateWallet(netParams,
			seed, []byte(pubpass), []byte(privpass),
			walletDir, walletName)
		errorHandler(err)

		err = w.Close()
		errorHandler(err)

		fmt.Println("Wallet created")
	},
}

var addrsCmd = &cobra.Command{
	Use:   "addresses",
	Short: "Address command",
}

var addrsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create address",
	Run: func(cmd *cobra.Command, args []string) {
		w := openWallet(pubpass, walletDir, walletName)
		addr, err := w.NewAddress()
		errorHandler(err)

		fmt.Printf("%s\n", addr.EncodeAddress())
	},
}

var balanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "Check total balance",
	Run: func(cmd *cobra.Command, args []string) {
		w := openWallet(pubpass, walletDir, walletName)
		utxos, err := w.ListUnspent()
		errorHandler(err)

		var total btcutil.Amount
		for _, utxo := range utxos {
			amt, err := btcutil.NewAmount(utxo.Amount)
			errorHandler(err)
			total += amt
		}

		fmt.Println(total.ToBTC())
	},
}

func openWallet(p string, dir string, name string) wallet.Wallet {
	netParams := loadNetParams(bitcoinConf)
	rpcclient := initRPCClient()

	w, err := _wallet.OpenWallet(
		netParams, []byte(p), dir, name, rpcclient)
	errorHandler(err)

	return w
}

func init() {
	// subcommand root
	walletsCmd.PersistentFlags().StringVar(
		&walletDir, "walletdir", "", "directory path to store wallets")
	walletsCmd.MarkPersistentFlagRequired("walletdir")
	rootCmd.AddCommand(walletsCmd)

	// create
	walletsCreateCmd.Flags().StringVar(&walletName, "walletname", "", "wallet name")
	walletsCreateCmd.MarkFlagRequired("walletname")
	walletsCreateCmd.Flags().StringVar(&pubpass, "pubpass", "", "public passphrase")
	walletsCreateCmd.MarkFlagRequired("pubpass")
	walletsCreateCmd.Flags().StringVar(&privpass, "privpass", "", "private passphrase")
	walletsCreateCmd.MarkFlagRequired("privpass")
	walletsCmd.AddCommand(walletsCreateCmd)

	// addresses
	walletsCmd.AddCommand(addrsCmd)

	// addresse create
	addrsCreateCmd.Flags().StringVar(&walletName, "walletname", "", "wallet name")
	addrsCreateCmd.MarkFlagRequired("walletname")
	addrsCreateCmd.Flags().StringVar(&pubpass, "pubpass", "", "public passphrase")
	addrsCreateCmd.MarkFlagRequired("pubpass")
	addrsCmd.AddCommand(addrsCreateCmd)

	// balance
	balanceCmd.Flags().StringVar(&walletName, "walletname", "", "wallet name")
	balanceCmd.MarkFlagRequired("walletname")
	balanceCmd.Flags().StringVar(&pubpass, "pubpass", "", "public passphrase")
	balanceCmd.MarkFlagRequired("pubpass")
	walletsCmd.AddCommand(balanceCmd)

}
