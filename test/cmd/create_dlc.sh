#!/bin/bash

net=${BITCOIN_NET:=regtest}
conf="bitcoin.${net}.conf"
conf_param="--conf ./conf/${conf}"

echo "Getting oracle's pubkey"
oracle_pubkey_file="opub.json"
dlccli oracle rpoints $conf_param \
    --oraclename "olivia" \
    --rpoints $RPOINTS \
    --fixingtime $FIX_TIME \
> $oracle_pubkey_file && cat $oracle_pubkey_file
echo -e ""

if [[ -z $PARTY1_TRANSFER_ADDRESS ]] || [[ -z $PARTY2_TRANSFER_ADDRESS ]];then
	echo "Transfer addresses not set, source create_addresses.sh script before."
	exit 1
fi

cmd="dlccli contracts"

if [[ "$PREMIUM_AMOUNT" -ne 0 ]];then
	cmd="$cmd createwithpremium \
        --premiumamount $PREMIUM_AMOUNT \
        --premiumdestaddress $PREMIUM_DEST_ADDRESS \
        --premiumpayingparty $PREMIUM_PAYING_PARTY"
else
  cmd="$cmd create"
fi

echo "Creating DLC"
cmd="$cmd $conf_param \
        --walletdir	$WALLET_DIR \
        --oracle_pubkey $oracle_pubkey_file \
        --fixingtime $FIX_TIME \
        --fund1 $PARTY1_FUND \
        --fund2 $PARTY2_FUND \
        --address1 $PARTY1_TRANSFER_ADDRESS \
        --address2 $PARTY2_TRANSFER_ADDRESS \
        --fundtx_feerate $FUND_TX_FEERATE \
        --redeemtx_feerate $REDEEM_TX_FEERATE \
        --deals_file $DEALS_FILE \
        --refund_locktime $REFUND_LOCKTIME \
        --wallet1 $PARTY1_WALLET_NAME \
        --wallet2 $PARTY2_WALLET_NAME \
        --pubpass1 $PARTY1_PUB_PASS \
        --pubpass2 $PARTY2_PUB_PASS \
        --privpass1 $PARTY1_PRIV_PASS \
        --privpass2 $PARTY2_PRIV_PASS"

if [[ -n "$PARTY1_CHANGE_ADDRESS" ]];then
  cmd="$cmd --change_address1 $PARTY1_CHANGE_ADDRESS"
fi

if [[ -n "$PARTY2_CHANGE_ADDRESS" ]];then
  cmd="$cmd --change_address2 $PARTY2_CHANGE_ADDRESS"
fi

echo $cmd
$cmd
