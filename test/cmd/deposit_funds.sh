#!/bin/bash

net=${BITCOIN_NET:=regtest}
if [[ "$net" == "maintest" ]];then
  echo "You shouldn't use this script on ${net}"
  exit 1
fi

if [[ -z $PARTY1_BASE_ADDRESS ]] || [[ -z $PARTY2_BASE_ADDRESS ]];then
	echo "Base addresses not set, source create_addresses.sh script before."
	exit 1
fi

conf="bitcoin.${net}.conf"
bitcoincli="bitcoin-cli -conf=`pwd`/conf/${conf}"

$bitcoincli sendtoaddress $PARTY1_BASE_ADDRESS 0.35
$bitcoincli sendtoaddress $PARTY2_BASE_ADDRESS 0.4

if [[ "${net}" == "regtest" ]];then
  $bitcoincli generate 1
fi
