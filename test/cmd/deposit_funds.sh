#!/bin/bash

bitcoincli="bitcoin-cli -conf=`pwd`/conf/bitcoin.regtest.conf"
dlc_params="--conf ./conf/bitcoin.regtest.conf --walletdir ./wallets/regtest"
create_address="dlccli wallets addresses create $dlc_params"
alice_params="--walletname alice --pubpass pub_alice"
# alice_personal_params="--walletname alicep --pubpass pub_alicep"
bob_params="--walletname bob --pubpass pub_bob"
# bob_personal_params="--walletname bobp --pubpass pub_bobp"

addr_a=`$create_address $alice_params`
addr_b=`$create_address $bob_params`
$bitcoincli sendtoaddress $addr_a 0.20022035
$bitcoincli sendtoaddress $addr_b 0.33355368
$bitcoincli generate 1
