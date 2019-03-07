package dlccli

import (
	"encoding/json"
	"fmt"

	_oracle "github.com/p2pderivatives/dlc/internal/oracle"
	"github.com/p2pderivatives/dlc/pkg/oracle"
	"github.com/spf13/cobra"
)

var oracleName string
var oracleRpoints int
var fixingValue int

// oracleCmd represents the oracle command
var oracleCmd = &cobra.Command{
	Use:   "oracle",
	Short: "oracle commands",
}

var oracleRpointsCmd = &cobra.Command{
	Use:   "rpoints",
	Short: "Get commited R points from Oracle",
	Run: func(cmd *cobra.Command, args []string) {
		o := initOracle()

		p, err := o.PubkeySet(parseFixingTimeFlag())
		errorHandler(err)

		pjson, err := json.Marshal(p)
		errorHandler(err)
		fmt.Println(string(pjson))
	},
}

var oracleMsgsCmd = &cobra.Command{
	Use: "messages",
}

var oracleFixMsgCmd = &cobra.Command{
	Use:   "fix",
	Short: "Fix message",
	Run: func(cmd *cobra.Command, args []string) {
		o := initOracle()
		o.InitDB()

		msgs := oracle.NumberToByteMsgs(fixingValue, oracleRpoints)

		ftime := parseFixingTimeFlag()
		o.FixMsgs(ftime, msgs)
		s, err := o.SignMsg(ftime)
		errorHandler(err)

		sjson, err := json.Marshal(s)
		errorHandler(err)

		fmt.Println(string(sjson))
	},
}

func initOracle() *_oracle.Oracle {
	netParams := loadChainParams(bitcoinConf)
	o, err := _oracle.New(oracleName, netParams, oracleRpoints)
	errorHandler(err)

	return o
}

func init() {
	// subcomand root
	oracleCmd.PersistentFlags().IntVar(
		&oracleRpoints, "rpoints", 0, "number of commited R points")
	oracleCmd.MarkPersistentFlagRequired("rpoints")
	oracleCmd.PersistentFlags().StringVar(
		&oracleName, "oraclename", "", "oracle name")
	oracleCmd.MarkPersistentFlagRequired("oraclename")
	oracleCmd.PersistentFlags().StringVar(
		&fixingTime, "fixingtime", "", "fixing time")
	oracleCmd.MarkPersistentFlagRequired("fixingtime")
	rootCmd.AddCommand(oracleCmd)

	// Rpoints
	oracleCmd.AddCommand(oracleRpointsCmd)

	// messagees
	oracleCmd.AddCommand(oracleMsgsCmd)

	// fix message
	oracleFixMsgCmd.PersistentFlags().IntVar(
		&fixingValue, "fixingvalue", 0, "fixing value")
	oracleFixMsgCmd.MarkPersistentFlagRequired("fixingvalue")
	oracleMsgsCmd.AddCommand(oracleFixMsgCmd)
}
