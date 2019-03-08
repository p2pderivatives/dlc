#!/bin/bash -euC

bitcoind=$(command -v bitcoind)
bitcoincli=$(command -v bitcoin-cli)
opts=( -datadir=./bitcoind -conf=./bitcoin.regtest.conf )

# start deamon if not running
function getnetworkinfo() {
  $bitcoincli "${opts[@]}" getnetworkinfo &> /dev/null
  echo $?
}
if [[ "$(getnetworkinfo)" -ne "0" ]];then
  $bitcoind "${opts[@]}"

  # wait until accepting rpc requests
  while true;do
    if [[ "$(getnetworkinfo)" -eq "0" ]];then
      break
    fi
    sleep 0.1
  done
else
  echo "Bitcoind is already running"
fi


function getblockcount() {
  echo $($bitcoincli "${opts[@]}" getblockcount)
}
height=$(getblockcount)
echo "Block Height: ${height}"

blocks=$((101 - $height))
if [[ "$blocks" -gt "0" ]];then
  echo "Generating initial regtest blocks"
  $bitcoincli "${opts[@]}" generate $blocks &> /dev/null
  height=$(getblockcount)
  echo "Block Height: ${height}"
fi
