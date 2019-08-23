#!/bin/bash

function create_address_for_party1() {
	local party1_wallet_name=${PARTY1_WALLET_NAME:=alice}
	local party1_pub_pass=${PARTY1_PUB_PATH:=pub_alice}
	local params="--walletname $party1_wallet_name --pubpass $party1_pub_pass"
	create_address "${params}"
}

function create_address_for_party2() {
	local party2_wallet_name=${PARTY2_WALLET_NAME:=bob}
	local party2_pub_pass=${PARTY2_PUB_PATH:=pub_bob}
	local params="--walletname $party2_wallet_name --pubpass $party2_pub_pass"
	create_address "${params}"
}

function create_address() {
	local net=${BITCOIN_NET:=regtest}
	local conf="bitcoin.${net}.conf"
	local bitcoincli="bitcoin-cli -conf=`pwd`/conf/${conf}"
	local wallet_dir=${WALLET_DIR:=./wallets/${net}}
	local dlc_params="--conf ./conf/${conf} --walletdir $wallet_dir"
	echo `dlccli wallets addresses create $dlc_params $1`
}

export PARTY1_BASE_ADDRESS=`create_address_for_party1`
export PARTY1_TRANSFER_ADDRESS=`create_address_for_party1`
export PARTY1_CHANGE_ADDRESS=`create_address_for_party1`

export PARTY2_BASE_ADDRESS=`create_address_for_party2`
export PARTY2_TRANSFER_ADDRESS=`create_address_for_party2`
export PARTY2_CHANGE_ADDRESS=`create_address_for_party2`

echo "Base address for party1: ${PARTY1_BASE_ADDRESS}"
echo "Transfer address for party1: ${PARTY1_TRANSFER_ADDRESS}"
echo "Change address for party1: ${PARTY1_CHANGE_ADDRESS}"

echo "Base address for party2: ${PARTY2_BASE_ADDRESS}"
echo "Transfer address for party2: ${PARTY2_TRANSFER_ADDRESS}"
echo "Change address for party2: ${PARTY2_CHANGE_ADDRESS}"

if [[ $PREMIUM_AMOUNT -ne 0 ]];then
	if [[ $PREMIUM_PAYING_PARTY -eq 0 ]];then
		export PREMIUM_DEST_ADDRESS=`create_address_for_party2`
	else
		export PREMIUM_DEST_ADDRESS=`create_address_for_party1`
	fi
	echo "Premium destination address: ${PREMIUM_DEST_ADDRESS}"
fi
