#!/bin/bash -euC

bitcoincli=$(command -v bitcoin-cli)
opts="-datadir=./bitcoind -conf=./bitcoin.regtest.conf"

$bitcoincli $opts stop
