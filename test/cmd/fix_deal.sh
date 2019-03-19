#!/bin/bash

net=${BITCOIN_NET:=regtest}
dlc_params="--conf ./conf/bitcoin.${net}.conf --walletdir ./wallets/${net}"
alice_params="--wallet alice --pubpass pub_alice --privpass priv_alice --contractor_type 0"
fix_deal="dlccli contracts deals fix"
oracle_sig="--oracle_sig ./osig.json"
dlc_id="--dlcid $1"

cmd="$fix_deal $dlc_params $oracle_sig $alice_params $dlc_id"
echo $cmd && $cmd
