#!/bin/bash -euC

bitcoind=$(command -v bitcoind)
bitcoincli=$(command -v bitcoin-cli)

# start deamon
echo "Starting bitcoind"
$bitcoind -datadir=./bitcoind -conf=./bitcoin.regtest.conf &

# wait until accepting rpc requests
echo "Waiting for initilization"
function getnetworkinfo() {
  $bitcoincli -datadir=./bitcoind -conf=bitcoin.regtest.conf getnetworkinfo &> /dev/null
  echo $?
}
while true;do
  if [[ "$(getnetworkinfo)" == "0" ]];then
    break
  fi
  sleep 0.1
done

echo "Generating initial regtest blocks"
$bitcoincli -datadir=./bitcoind -conf=bitcoin.regtest.conf generate 101 &> /dev/null
