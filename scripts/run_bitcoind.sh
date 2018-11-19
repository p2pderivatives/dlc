#!/bin/bash -euC

bitcoind=$(command -v bitcoind)

$bitcoind -datadir=./bitcoind -conf=./bitcoin.regtest.conf
