#!/bin/bash

net=${BITCOIN_NET:=regtest}
dlc_params="--conf ./conf/bitcoin.${net}.conf --walletdir $WALLET_DIR"
party1_params="--wallet $PARTY1_WALLET_NAME --pubpass $PARTY1_PUB_PASS --privpass $PARTY1_PRIV_PASS --contractor_type 0"
fix_deal="dlccli contracts deals fix"
oracle_sig="--oracle_sig ./osig.json"
dlc_id="--dlcid $1"

cmd="$fix_deal $dlc_params $oracle_sig $party1_params $dlc_id"
echo $cmd && $cmd
