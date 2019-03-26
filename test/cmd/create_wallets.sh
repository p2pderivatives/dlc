#!/bin/bash

net=${BITCOIN_NET:=regtest}
conf="--conf ./conf/bitcoin.${net}.conf"
dlc_params="${conf} --walletdir ./wallets/${net}"
alice_params=( --walletname alice --pubpass pub_alice --privpass priv_alice )
alice_personal_params=( --walletname alicep --pubpass pub_alicep --privpass priv_alicep )
bob_params=( --walletname bob --pubpass pub_bob --privpass priv_bob )
bob_personal_params=( --walletname bobp --pubpass pub_bobp --privpass priv_bobp )
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

echo "Creating Alice's DLC wallet"
create_wallet "${alice_params[@]}"

echo "Creating Alice's personal wallet"
create_wallet "${alice_personal_params[@]}"

echo "Creating Bob's DLC wallet"
create_wallet "${bob_params[@]}"

echo "Creating Bob's personal wallet"
create_wallet "${bob_personal_params[@]}"
