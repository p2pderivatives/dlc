package dlccli

import (
	"fmt"
	"path/filepath"

	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/btcsuite/btcwallet/walletdb"
	"github.com/spf13/cobra"

	"github.com/p2pderivatives/dlc/internal/dlcmgr"
	_wallet "github.com/p2pderivatives/dlc/internal/wallet"
	"github.com/p2pderivatives/dlc/pkg/wallet"
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
		chainParams := loadChainParams(bitcoinConf)

		// TODO: give seed as command line parameter
		seed, err := hdkeychain.GenerateSeed(hdkeychain.RecommendedSeedLen)
		errorHandler(err)

		w, err := _wallet.CreateWallet(chainParams,
			seed, []byte(pubpass), []byte(privpass),
			walletDir, walletName)
		errorHandler(err)

		err = w.Close()
		errorHandler(err)

		_, wdb := openWallet(pubpass, walletDir, walletName)
		defer wdb.Close()
		_, err = dlcmgr.Create(wdb)
		errorHandler(err)
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
		w, _ := openWallet(pubpass, walletDir, walletName)
		addr, err := w.NewAddress()
		errorHandler(err)

		fmt.Printf("%s\n", addr.EncodeAddress())
	},
}

var balanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "Check total balance",
	Run: func(cmd *cobra.Command, args []string) {
		w, _ := openWallet(pubpass, walletDir, walletName)
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

func openWallet(pubpass string, dir string, name string) (wallet.Wallet, walletdb.DB) {
	chainParams := loadChainParams(bitcoinConf)
	rpcclient := initRPCClient()
	wdb := openWalletDB(dir, name)

	w, err := _wallet.Open(wdb, []byte(pubpass), chainParams, rpcclient)
	errorHandler(err)

	return w, wdb
}

func openWalletDB(dir string, name string) walletdb.DB {
	dbpath := filepath.Join(dir, name+".db")
	wdb, err := walletdb.Open("bdb", dbpath)
	errorHandler(err)
	return wdb
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
