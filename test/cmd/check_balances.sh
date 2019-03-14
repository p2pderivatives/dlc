#!/bin/bash

dlc_params="--conf ./conf/bitcoin.regtest.conf --walletdir ./wallets/regtest"
balance="dlccli wallets balance $dlc_params"
alice_params="--walletname alice --pubpass pub_alice"
bob_params="--walletname bob --pubpass pub_bob"

echo "Balance Alice: $($balance $alice_params)"
echo "Balance Bob: $($balance $bob_params)"
