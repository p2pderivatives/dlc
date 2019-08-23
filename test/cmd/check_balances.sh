#!/bin/bash

net=${BITCOIN_NET:=regtest}
conf="bitcoin.${net}.conf"
wallet_dir=${WALLET_DIR:="./wallets/${net}"}
dlc_params="--conf ./conf/${conf} --walletdir $wallet_dir"
balance="dlccli wallets balance $dlc_params"
party1_params="--walletname ${PARTY1_WALLET_NAME} --pubpass ${PARTY1_PUB_PASS}"
party2_params="--walletname ${PARTY2_WALLET_NAME} --pubpass ${PARTY2_PUB_PASS}"

echo "Party1 (DLC Wallet): $($balance $party1_params)"
echo "Party2 (DLC Wallet): $($balance $party2_params)"
