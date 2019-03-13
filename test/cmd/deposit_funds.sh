#!/bin/bash

bitcoincli="bitcoin-cli -conf=`pwd`/conf/bitcoin.regtest.conf"
dlc_params="--conf ./conf/bitcoin.regtest.conf --walletdir ./wallets/regtest"
create_address="dlccli wallets addresses create $dlc_params"
balance="dlccli wallets balance $dlc_params"
alice_params="--walletname alice --pubpass pub_alice"
bob_params="--walletname bob --pubpass pub_bob"

addr_a=`$create_address $alice_params`
addr_b=`$create_address $bob_params`
$bitcoincli sendtoaddress $addr_a 10.0
$bitcoincli sendtoaddress $addr_b 10.0
$bitcoincli generate 1

echo "Balance Alice: $($balance $alice_params)"
echo "Balance Bob: $($balance $bob_params)"
