#!/bin/bash -euC

bitcoincli=$(command -v bitcoin-cli)

function stop() {
  $bitcoincli -datadir=./bitcoind -conf=bitcoin.regtest.conf stop &> /dev/null
  echo $?
}

echo "Stopping bitcoind"
while true;do
  if [[ "$(stop)" == "1" ]];then
    sleep 0.2
    break
  fi
  sleep 0.1
done

echo "Deleting regtest directory"
rm -rf ./bitcoind/regtest
