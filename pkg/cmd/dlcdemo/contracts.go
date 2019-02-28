package dlccli

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/btcsuite/btcutil"
	"github.com/dgarage/dlc/pkg/dlc"
	"github.com/dgarage/dlc/pkg/oracle"
	"github.com/dgarage/dlc/pkg/wallet"
	"github.com/spf13/cobra"
)

var contractsCmd = &cobra.Command{
	Use:   "contracts",
	Short: "Contract commands",
}

var fund1 int
var fund2 int
var address1 string
var address2 string
var fundtxFeerate int
var redeemtxFeerate int

// var refundLocktime string
var dealsFile string
var wallet1 string
var wallet2 string
var pubpass1 string
var pubpass2 string
var privpass1 string
var privpass2 string

func initCreateContractCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "create",
		Short: "Create contract",
		Run: func(cmd *cobra.Command, args []string) {
			party1 := initFirstParty()
			party2 := initSecondParty()
			fmt.Println(party1)
			fmt.Println(party2)
			fmt.Println("Contract created")
		},
	}

	cmd.Flags().StringVar(&fixingTime, "fixingtime", "", "fixing time")
	cmd.MarkFlagRequired("fixingtime")
	cmd.Flags().IntVar(&fund1, "fund1", 0, "Fund amount of First party (satoshi)")
	cmd.MarkFlagRequired("fund1")
	cmd.Flags().IntVar(&fund2, "fund2", 0, "Fund amount of Second party (satoshi)")
	cmd.MarkFlagRequired("fund2")
	cmd.Flags().StringVar(&address1, "address1", "", "Transfer address of First party")
	cmd.MarkFlagRequired("address1")
	cmd.Flags().StringVar(&address2, "address2", "", "Transfer address of Second party")
	cmd.MarkFlagRequired("address2")
	cmd.Flags().IntVar(&fundtxFeerate, "fundtx_feerate", 0, "Fee rate for fund tx (satoshi/byte)")
	cmd.MarkFlagRequired("fundtx_feerate")
	cmd.Flags().IntVar(&redeemtxFeerate, "redeemtx_feerate", 0, "Fee rate for refund tx, cetx, closing tx (satoshi/byte)")
	cmd.MarkFlagRequired("redeemtx_feerate")
	// cmd.MarkFlagRequired("refund_locktime")
	// cmd.Flags().StringVar(&refund_locktime, "refund_locktime", "", "Locktime of refune tx")
	cmd.Flags().StringVar(&dealsFile, "deals_file", "", "Path to a csv file that contains deals")
	cmd.MarkFlagRequired("deals_file")
	cmd.Flags().StringVar(&walletDir, "walletdir", "", "directory path to store wallets")
	cmd.MarkFlagRequired("walletdir")
	cmd.Flags().StringVar(&wallet1, "wallet1", "", "wallet name of First Party")
	cmd.MarkFlagRequired("wallet1")
	cmd.Flags().StringVar(&wallet2, "wallet2", "", "wallet name of Second Party")
	cmd.MarkFlagRequired("wallet_2")
	cmd.Flags().StringVar(&pubpass1, "pubpass1", "", "Pubpass phrase of First party's wallet")
	cmd.MarkFlagRequired("pubpass1")
	cmd.Flags().StringVar(&pubpass2, "pubpass2", "", "Pubpass phrase of Second party's wallet")
	cmd.MarkFlagRequired("pubpass2")
	cmd.Flags().StringVar(&privpass1, "privpass1", "", "Privpass phrase of First party's wallet")
	cmd.MarkFlagRequired("privpass1")
	cmd.Flags().StringVar(&privpass2, "privpass2", "", "Privpass phrase of Second party's wallet")
	cmd.MarkFlagRequired("privpass2")

	return cmd
}

// Contractor is contractor
type Contractor struct {
	wallet   wallet.Wallet
	builder  *dlc.Builder
	pubpass  string
	privpass string
}

func loadDLCConditions() *dlc.Conditions {
	ftime := parseFixingTimeFlag()

	// cast int to btcutil.Amount
	famt1 := btcutil.Amount(fund1)
	famt2 := btcutil.Amount(fund2)
	ffrate := btcutil.Amount(fundtxFeerate)
	rfrate := btcutil.Amount(redeemtxFeerate)

	// TODO: confirm how to convert timestamp to locktime
	lc := uint32(1)

	deals := loadDeals()

	conds, err := dlc.NewConditions(
		ftime, famt1, famt2, ffrate, rfrate, lc, deals)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return conds
}

func loadDeals() []*dlc.Deal {
	f, err := os.Open(dealsFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	nDigits := 5

	deals := []*dlc.Deal{}
	r := csv.NewReader(bufio.NewReader(f))
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		deal := convertRowToDeal(row, nDigits)
		deals = append(deals, deal)
	}

	return deals
}

func convertRowToDeal(rec []string, nDigits int) *dlc.Deal {
	v, err := strconv.Atoi(rec[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	msgs := oracle.NumberToByteMsgs(v, nDigits)

	amt1, err := strconv.Atoi(rec[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	amt2, err := strconv.Atoi(rec[2])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	deal := dlc.NewDeal(
		btcutil.Amount(amt1),
		btcutil.Amount(amt2),
		msgs)

	return deal
}

func initFirstParty() *Contractor {
	w := openWallet(pubpass1, walletDir, wallet1)
	conds := loadDLCConditions()
	b := dlc.NewBuilder(dlc.FirstParty, w, conds)

	return &Contractor{
		wallet:   w,
		builder:  b,
		pubpass:  pubpass1,
		privpass: privpass1,
	}
}

func initSecondParty() *Contractor {
	w := openWallet(pubpass2, walletDir, wallet2)
	conds := loadDLCConditions()
	b := dlc.NewBuilder(dlc.SecondParty, w, conds)

	return &Contractor{
		wallet:   w,
		builder:  b,
		pubpass:  pubpass2,
		privpass: privpass2,
	}
}

func init() {
	// subcommand root
	rootCmd.AddCommand(contractsCmd)

	// create contract
	contractsCmd.AddCommand(initCreateContractCmd())
}
