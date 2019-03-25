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

while [[ "$#" -ne "0" ]];do
  opts+=( ${1} )
  shift
done

if [[ "$(getnetworkinfo)" -ne "0" ]];then
  $bitcoind "${opts[@]}"

  # wait until accepting rpc requests
  cnt=0
  while true;do
    if [[ "$(getnetworkinfo)" -eq "0" ]];then
      break
    fi
    if [[ "$cnt" -gt "100" ]];then
      echo "Failed to start bitcoind. see debug.log for more details."
      exit 1
    fi
    sleep 0.1
    cnt=$(($cnt+1))
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
