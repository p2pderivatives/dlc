// Copyright © 2018 Junji Watanabe <junji-watanabe@garage.co.jp>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package dlccli

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/p2pderivatives/dlc/internal/rpc"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var debug bool
var bitcoinConf string
var walletDir string
var fixingTime string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dlccli",
	Short: "DLC command line interface",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	errorHandler(err)
}

func init() {
	initLogger(rootCmd)

	rootCmd.PersistentFlags().StringVar(
		&bitcoinConf, "conf", "", "bitcoin config file")
	rootCmd.MarkPersistentFlagRequired("conf")
}

func initLogger(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVar(
		&debug, "debug", false, "enable debug logs")

	cobra.OnInitialize(func() {
		cfg := zap.NewDevelopmentConfig()
		if debug {
			cfg.Level.SetLevel(zap.DebugLevel)
		} else {
			cfg.Level.SetLevel(zap.InfoLevel)
		}
		logger, err := cfg.Build()
		errorHandler(err)
		zap.ReplaceGlobals(logger)
	})
}

func logger() *zap.Logger {
	return zap.L
}

func loadChainParams(cfgPath string) *chaincfg.Params {
	cfgFile, err := os.Open(cfgPath)
	errorHandler(err)
	defer cfgFile.Close()

	content, err := ioutil.ReadAll(cfgFile)
	errorHandler(err)

	if extractValue(content, "regtest") == "1" {
		return &chaincfg.RegressionNetParams
	} else if extractValue(content, "testnet") == "1" {
		return &chaincfg.TestNet3Params
	} else {
		return &chaincfg.MainNetParams
	}
}

func extractValue(content []byte, key string) string {
	pattern := fmt.Sprintf(`(?m)^\s*%s=([^\s]+)`, key)
	reg, _ := regexp.Compile(pattern)
	matches := reg.FindSubmatch(content)

	if matches == nil {
		return ""
	}

	return strings.Split(string(matches[0]), "=")[1]
}

func parseFixingTimeFlag() time.Time {
	t, err := time.Parse(time.RFC3339, fixingTime)
	errorHandler(err)
	return t
}

func initRPCClient() rpc.Client {
	rpcclient, err := rpc.NewClient(bitcoinConf)
	errorHandler(err)
	return rpcclient
}

func errorHandler(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
