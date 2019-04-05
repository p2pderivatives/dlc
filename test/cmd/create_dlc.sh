#!/bin/bash

net=${BITCOIN_NET:=regtest}
conf="bitcoin.${net}.conf"
conf="--conf ./conf/${conf}"
walletdir="--walletdir ./wallets/${net}"
create_address="dlccli wallets addresses create"
alicep_params="--walletname alicep --pubpass pub_alicep"
bobp_params="--walletname bobp --pubpass pub_bobp"

echo "Getting oracle's pubkey"
oracle_pubkey_file="opub.json"
dlccli oracle rpoints $conf \
    --oraclename "olivia" \
    --rpoints 4 \
    --fixingtime "2019-04-30T12:00:00Z" \
> $oracle_pubkey_file && cat $oracle_pubkey_file
echo -e ""

echo "Creating addresses"
addr1=`$create_address $conf $walletdir $alicep_params`
addr2=`$create_address $conf $walletdir $bobp_params`
echo "address1: $addr1"
echo "address2: $addr2"
chaddr1=`$create_address $conf $walletdir $alicep_params`
chaddr2=`$create_address $conf $walletdir $bobp_params`
echo "change address1: $chaddr1"
echo "change address2: $chaddr2"
echo -e ""

echo "Creating DLC"
cmd="dlccli contracts create $conf $walletdir \
        --oracle_pubkey $oracle_pubkey_file \
        --fixingtime 2019-04-30T12:00:00Z \
        --fund1 2000 \
        --fund2 3333 \
        --address1 $addr1 \
        --address2 $addr2 \
        --change_address1 $chaddr1 \
        --change_address2 $chaddr2 \
        --fundtx_feerate 20 \
        --redeemtx_feerate 20 \
        --deals_file ./test/cmd/deals_qa.csv \
        --refund_locktime 574196 \
        --wallet1 alice \
        --wallet2 bob \
        --pubpass1 pub_alice \
        --pubpass2 pub_bob \
        --privpass1 priv_alice \
        --privpass2 priv_bob"

if [[ "$chaddr1" -ne "" ]];then
  cmd="$cmd --change_address1 $chaddr1"
fi
if [[ "$chaddr2" -ne "" ]];then
  cmd="$cmd --change_address2 $chaddr2"
fi

echo $cmd
$cmd
