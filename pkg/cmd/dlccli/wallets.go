package dlccli

import (
	"encoding/hex"
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
var walletsCmd = func() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "wallets",
		Short: "Wallet commands",
	}

	return cmd
}

var walletsCreateCmd = func() *cobra.Command {
	cmd := &cobra.Command{
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

	cmd.Flags().StringVar(&walletDir, "walletdir", "", "directory path to store wallets")
	cmd.MarkFlagRequired("walletdir")
	cmd.Flags().StringVar(&walletName, "walletname", "", "wallet name")
	cmd.MarkFlagRequired("walletname")
	cmd.Flags().StringVar(&pubpass, "pubpass", "", "public passphrase")
	cmd.MarkFlagRequired("pubpass")
	cmd.Flags().StringVar(&privpass, "privpass", "", "private passphrase")
	cmd.MarkFlagRequired("privpass")

	return cmd
}

var walletsSeedCmd = func() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "seed",
		Short: "Generate seed",
		Run: func(cmd *cobra.Command, args []string) {
			seed, err := hdkeychain.GenerateSeed(hdkeychain.RecommendedSeedLen)
			errorHandler(err)

			seedHex := hex.EncodeToString(seed)
			fmt.Println(seedHex)
		},
	}
	return cmd
}

var addrsCmd = func() *cobra.Command {
	return &cobra.Command{
		Use:   "addresses",
		Short: "Address command",
	}
}

var addrsCreateCmd = func() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create address",
		Run: func(cmd *cobra.Command, args []string) {
			w, _ := openWallet(pubpass, walletDir, walletName)
			addr, err := w.NewAddress()
			errorHandler(err)

			fmt.Printf("%s\n", addr.EncodeAddress())
		},
	}

	cmd.Flags().StringVar(&walletDir, "walletdir", "", "directory path to store wallets")
	cmd.MarkFlagRequired("walletdir")
	cmd.Flags().StringVar(&walletName, "walletname", "", "wallet name")
	cmd.MarkFlagRequired("walletname")
	cmd.Flags().StringVar(&pubpass, "pubpass", "", "public passphrase")
	cmd.MarkFlagRequired("pubpass")

	return cmd
}

var balanceCmd = func() *cobra.Command {
	cmd := &cobra.Command{
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

	cmd.Flags().StringVar(&walletDir, "walletdir", "", "directory path to store wallets")
	cmd.MarkFlagRequired("walletdir")
	cmd.Flags().StringVar(&walletName, "walletname", "", "wallet name")
	cmd.MarkFlagRequired("walletname")
	cmd.Flags().StringVar(&pubpass, "pubpass", "", "public passphrase")
	cmd.MarkFlagRequired("pubpass")

	return cmd
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
	subRootCmd := walletsCmd()
	rootCmd.AddCommand(subRootCmd)

	// create wallet
	subRootCmd.AddCommand(walletsCreateCmd())

	// seed
	subRootCmd.AddCommand(walletsSeedCmd())

	// balance
	subRootCmd.AddCommand(balanceCmd())

	// addresses sub command root
	addrsRootCmd := addrsCmd()
	subRootCmd.AddCommand(addrsRootCmd)

	// addresse create
	addrsRootCmd.AddCommand(addrsCreateCmd())
}
