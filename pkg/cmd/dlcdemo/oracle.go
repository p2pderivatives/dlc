package dlccli

import (
	"fmt"
	"os"

	_oracle "github.com/dgarage/dlc/internal/oracle"
	"github.com/spf13/cobra"
)

var oracleName string
var oracleRpoints int

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
		ftime := parseFixingTimeFlag()
		p, err := o.PubkeySet(ftime)
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
}
