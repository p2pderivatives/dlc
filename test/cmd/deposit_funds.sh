#!/bin/bash

net=${BITCOIN_NET:=regtest}
if [[ "$net" == "maintest" ]];then
  echo "You shouldn't use this script on ${net}"
  exit 1
fi

conf="bitcoin.${net}.conf"
bitcoincli="bitcoin-cli -conf=`pwd`/conf/${conf}"
dlc_params="--conf ./conf/${conf} --walletdir ./wallets/${net}"
create_address="dlccli wallets addresses create $dlc_params"
alice_params="--walletname alice --pubpass pub_alice"
bob_params="--walletname bob --pubpass pub_bob"

addr_a=`$create_address $alice_params`
addr_b=`$create_address $bob_params`
$bitcoincli sendtoaddress $addr_a 0.00011360
$bitcoincli sendtoaddress $addr_b 0.00012693

if [[ "${net}" == "regtest" ]];then
  $bitcoincli generate 1
fi
