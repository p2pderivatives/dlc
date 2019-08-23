#!/bin/bash

net=${BITCOIN_NET:=regtest}

party1_wallet_name=${PARTY1_WALLET_NAME:=alice}
party1_priv_pass=${PARTY1_PRIV_PATH:=priv_alice}
party1_pub_pass=${PARTY1_PUB_PATH:=pub_alice}

party2_wallet_name=${PARTY2_WALLET_NAME:=bob}
party2_priv_pass=${PARTY2_PRIV_PATH:=priv_bob}
party2_pub_pass=${PARTY2_PUB_PATH:=pub_bob}

conf="--conf ./conf/bitcoin.${net}.conf"
dlc_params="${conf} --walletdir ./wallets/${net}"
party1_params=( --walletname ${party1_wallet_name} --pubpass ${party1_pub_pass} --privpass ${party1_priv_pass} )
party2_params=( --walletname ${party2_wallet_name} --pubpass ${party2_pub_pass} --privpass ${party2_priv_pass} )
create_wallet="dlccli wallets create"

function create_wallet() {
  seed=$(dlccli wallets seed $conf)
  wallet_params="--seed $seed"
  for i in "$@";
  do
    wallet_params="$wallet_params $i"
  done
  cmd="$create_wallet $dlc_params $wallet_params"
  echo $cmd
  $cmd
}

echo "Creating Party 1's DLC wallet"
create_wallet "${party1_params[@]}"

echo "Creating Party 2's DLC wallet"
create_wallet "${party2_params[@]}"
