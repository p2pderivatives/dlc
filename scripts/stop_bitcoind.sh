#!/bin/bash -euC

bitcoincli=$(command -v bitcoin-cli)
net=${BITCOIN_NET:=regtest}
conf="bitcoin.${net}.conf"
opts=( -datadir=./bitcoind -conf=$conf )

$bitcoincli "${opts[@]}" stop
