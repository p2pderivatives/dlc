package dlccli

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/p2pderivatives/dlc/internal/dlcmgr"
	"github.com/p2pderivatives/dlc/pkg/dlc"
	"github.com/p2pderivatives/dlc/pkg/oracle"
	"github.com/p2pderivatives/dlc/pkg/utils"
	"github.com/spf13/cobra"
)

var dlcid string
var osigfile string
var contractorType int

func runFixDeal(cmd *cobra.Command, args []string) {
	c := initCotractor()

	osig := parseOracleSignedMsg()

	idxs := []int{}
	n := len(osig.Sigs)
	for i := 0; i < n; i++ {
		idxs = append(idxs, i)
	}
	err := c.builder.FixDeal(osig, idxs)
	errorHandler(err)

	cetx, err := c.builder.SignedContractExecutionTx()
	errorHandler(err)

	cetxHex, err := utils.TxToHex(cetx)
	errorHandler(err)
	fmt.Printf("\nCETx hex:\n%s\n", cetxHex)

	cltx, err := c.builder.SignedClosingTx(cetx)
	errorHandler(err)
	cltxHex, err := utils.TxToHex(cltx)
	errorHandler(err)
	fmt.Printf("\nClosingTx hex:\n%s\n", cltxHex)
}

func initCotractor() *Contractor {
	w, wdb := openWallet(pubpass, walletDir, walletName)
	err := w.Unlock([]byte(privpass))
	errorHandler(err)
	mgr, err := dlcmgr.Open(wdb)
	errorHandler(err)

	h, err := chainhash.NewHashFromStr(dlcid)
	errorHandler(err)
	key := h.CloneBytes()
	d, err := mgr.RetrieveContract(key)
	errorHandler(err)

	b := dlc.NewBuilderFromDLC(
		d, dlc.Contractor(contractorType), w)

	return &Contractor{
		wallet:   w,
		builder:  b,
		manager:  mgr,
		pubpass:  pubpass,
		privpass: privpass,
	}
}

func parseOracleSignedMsg() *oracle.SignedMsg {
	data, err := ioutil.ReadFile(osigfile)
	errorHandler(err)

	signedMsg := &oracle.SignedMsg{}
	err = json.Unmarshal(data, signedMsg)
	errorHandler(err)

	return signedMsg
}

func initFixDealCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "fix",
		Short: "Fix deal",
		Run:   runFixDeal,
	}

	cmd.Flags().StringVar(&dlcid, "dlcid", "", "Contract ID")
	cmd.MarkFlagRequired("dlcid")
	cmd.Flags().StringVar(&osigfile, "oracle_sig", "", "Oracle's signed message json file")
	cmd.MarkFlagRequired("oracle_sig")
	cmd.Flags().StringVar(&walletDir, "walletdir", "", "Wallet directory")
	cmd.MarkFlagRequired("walletdir")
	cmd.Flags().StringVar(&walletName, "wallet", "", "Wallet name")
	cmd.MarkFlagRequired("wallet")
	cmd.Flags().IntVar(&contractorType, "contractor_type", 0, "0: first party, 1:second party")
	cmd.MarkFlagRequired("contractor_type")
	cmd.Flags().StringVar(&pubpass, "pubpass", "", "public passphrase")
	cmd.MarkFlagRequired("pubpass")
	cmd.Flags().StringVar(&privpass, "privpass", "", "private passphrase")
	cmd.MarkFlagRequired("privpass")

	return cmd
}
