#!/bin/bash

net=${BITCOIN_NET:=regtest}
conf="--conf ./conf/bitcoin.${net}.conf"
oracle_params="--oraclename olivia --rpoints $RPOINTS"
fix_message="dlccli oracle messages fix"
value=$1

$fix_message $conf $oracle_params \
    --fixingtime $FIX_TIME \
    --fixingvalue $value \
    > osig.json && cat osig.json
