#!/bin/bash

dlc_params="--conf ./conf/bitcoin.regtest.conf --walletdir ./wallets/regtest"
alice_params="--walletname alice --pubpass pub_alice --privpass priv_alice"
alice_personal_params="--walletname alicep --pubpass pub_alicep --privpass priv_alicep"
alice_params="--walletname alice --pubpass pub_alice --privpass priv_alice"
bob_params="--walletname bob --pubpass pub_bob --privpass priv_bob"
bob_personal_params="--walletname bobp --pubpass pub_bobp --privpass priv_bobp"
create_wallet="dlccli wallets create"

echo "Creating Alice's DLC wallet"
$create_wallet $dlc_params $alice_params
echo "Creating Alice's personal wallet"
$create_wallet $dlc_params $alice_personal_params
echo "Creating Bob's DLC wallet"
$create_wallet $dlc_params $bob_params
echo "Creating Bob's personal wallet"
$create_wallet $dlc_params $bob_personal_params
