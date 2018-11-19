#!/bin/bash -euC

bitcoincli=$(command -v bitcoin-cli)

$bitcoincli -datadir=./bitcoind -conf=./bitcoin.regtest.conf stop
