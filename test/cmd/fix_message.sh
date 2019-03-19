#!/bin/bash

net=${BITCOIN_NET:=regtest}
conf="--conf ./conf/bitcoin.${net}.conf"
oracle_params="--oraclename olivia --rpoints 4"
fix_message="dlccli oracle messages fix"
value=$1

$fix_message $conf $oracle_params \
    --fixingtime "2019-03-30T12:00:00Z" \
    --fixingvalue $value \
    > osig.json && cat osig.json
