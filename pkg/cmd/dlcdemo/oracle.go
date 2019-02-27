package dlccli

import (
	"fmt"
	"os"

	_oracle "github.com/dgarage/dlc/internal/oracle"
	"github.com/dgarage/dlc/pkg/oracle"
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
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		pjson, err := p.ToJSON()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
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
		s, err := o.SignSet(ftime)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		sjson, err := s.ToJSON()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(string(sjson))
	},
}

func initOracle() *_oracle.Oracle {
	netParams := loadNetParams(bitcoinConf)
	o, err := _oracle.New(oracleName, netParams, oracleRpoints)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

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
