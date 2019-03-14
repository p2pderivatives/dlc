#!/bin/bash

conf="--conf=`pwd`/conf/bitcoin.regtest.conf"
oracle_params="--oraclename olivia --rpoints 5"
fix_message="dlccli oracle rpoints"

$fix_message $conf $oracle_params \
    --fixingtime "2019-03-30T12:00:00Z" \
    > opub.json && cat opub.json

