package dlccli

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/spf13/cobra"

	"github.com/dgarage/dlc/internal/rpc"
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

		w, err := _wallet.CreateWallet(netParams,
			seed, []byte(pubpass), []byte(privpass),
			walletDir, walletName)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = w.Close()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

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
		w := openWallet()
		addr, err := w.NewAddress()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("%s\n", addr.EncodeAddress())
	},
}

var balanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "Check total balance",
	Run: func(cmd *cobra.Command, args []string) {
		w := openWallet()
		utxos, err := w.ListUnspent()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		var total btcutil.Amount
		for _, utxo := range utxos {
			amt, err := btcutil.NewAmount(utxo.Amount)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			total += amt
		}

		fmt.Println(total.ToBTC())
	},
}

func openWallet() wallet.Wallet {
	netParams, err := loadNetParams(bitcoinConf)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	rpcclient, err := rpc.NewClient(bitcoinConf)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	w, err := _wallet.OpenWallet(
		netParams, []byte(pubpass), walletDir, walletName, rpcclient)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return w
}

func init() {
	// subcommand root
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
