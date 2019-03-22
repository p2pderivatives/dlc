#!/bin/bash -euC

bitcoind=$(command -v bitcoind)
bitcoincli=$(command -v bitcoin-cli)
net=${BITCOIN_NET:=regtest}
conf="bitcoin.${net}.conf"
opts=( -datadir=./bitcoind -conf=$conf )

# start deamon if not running
function getnetworkinfo() {
  $bitcoincli "${opts[@]}" getnetworkinfo &> /dev/null
  echo $?
}

function run_bitcoind() {
  $bitcoind "${opts[@]}"
  echo $?
}

if [[ "$(getnetworkinfo)" -ne "0" ]];then
  if [[ "$(run_bitcoind)" -ne "0" ]];then
    exit 1
  fi

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

if [[ "$net" == "regtest" ]];then
  blocks=$((101 - $height))
  if [[ "$blocks" -gt "0" ]];then
    echo "Generating initial regtest blocks"
    $bitcoincli "${opts[@]}" generate $blocks &> /dev/null
    height=$(getblockcount)
    echo "Block Height: ${height}"
  fi
fi
