#!/bin/bash

conf="--conf=`pwd`/conf/bitcoin.regtest.conf"
oracle_params="--oraclename olivia --rpoints 5"
fix_message="dlccli oracle messages fix"

$fix_message $conf $oracle_params \
    --fixingtime "2019-03-30T12:00:00Z" \
    --fixingvalue 30001 \
    > osig.json && cat osig.json
