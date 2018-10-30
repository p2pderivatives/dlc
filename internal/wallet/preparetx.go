package wallet

import (
	"sort"

	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcwallet/wallet/txauthor"
	_ "github.com/btcsuite/btcwallet/walletdb/bdb" // blank import for bolt db driver
	"github.com/btcsuite/btcwallet/wtxmgr"
)

// byAmount defines the methods needed to satisify sort.Interface to
// sort credits by their output amount.
type byAmount []wtxmgr.Credit

func (s byAmount) Len() int           { return len(s) }
func (s byAmount) Less(i, j int) bool { return s[i].Amount < s[j].Amount }
func (s byAmount) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// FundTx adds inputs to a transaction until amount.
func (w *Wallet) FundTx(eligible []wtxmgr.Credit) txauthor.InputSource {

	// outs := []*wire.OutPoint{}
	// total := int64(0)
	// addfee := int64(0)
	// for _, utxo := range list {
	// 	txid, _ := chainhash.NewHashFromStr(utxo.TxID)
	// 	outs = append(outs, wire.NewOutPoint(txid, utxo.Vout))
	// 	a, _ := btcutil.NewAmount(utxo.Amount)
	// 	total += int64(a)
	// 	addfee = int64(len(outs)) * dlc.DlcTxInSize * efee
	// 	if amount+addfee <= total {
	// 		if amount+addfee == total {
	// 			break
	// 		}
	// 		addfee += dlc.DlcTxOutSize * efee
	// 		if amount+addfee <= total {
	// 			break
	// 		}
	// 	}
	// }
	// if amount+addfee > total {
	// 	return fmt.Errorf("short of bitcoin")
	// }
	// for _, out := range outs {
	// 	tx.AddTxIn(wire.NewTxIn(out, nil, nil))
	// }
	// if amount+addfee == total {
	// 	return nil
	// }
	// change := total - (amount + addfee)
	// pkScript := w.P2WPKHpkScript(w.GetFakePublicKey())
	// tx.AddTxOut(wire.NewTxOut(change, pkScript))
	// return nil

	// Pick largest outputs first.  This is only done for compatibility with
	// previous tx creation code, not because it's a good idea.
	sort.Sort(sort.Reverse(byAmount(eligible)))

	// Current inputs and their total value.  These are closed over by the
	// returned input source and reused across multiple calls.
	currentTotal := btcutil.Amount(0)
	currentInputs := make([]*wire.TxIn, 0, len(eligible))
	currentScripts := make([][]byte, 0, len(eligible))
	currentInputValues := make([]btcutil.Amount, 0, len(eligible))

	return func(target btcutil.Amount) (btcutil.Amount, []*wire.TxIn,
		[]btcutil.Amount, [][]byte, error) {

		for currentTotal < target && len(eligible) != 0 {
			nextCredit := &eligible[0]
			eligible = eligible[1:]
			nextInput := wire.NewTxIn(&nextCredit.OutPoint, nil, nil)
			currentTotal += nextCredit.Amount
			currentInputs = append(currentInputs, nextInput)
			currentScripts = append(currentScripts, nextCredit.PkScript)
			currentInputValues = append(currentInputValues, nextCredit.Amount)
		}
		return currentTotal, currentInputs, currentInputValues, currentScripts, nil
	}

}
