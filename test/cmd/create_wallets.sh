#!/bin/bash

dlc_params="--conf ./conf/bitcoin.regtest.conf --walletdir ./wallets/regtest"
alice_params="--walletname alice --pubpass pub_alice --privpass priv_alice"
bob_params="--walletname bob --pubpass pub_bob --privpass priv_bob"
create_wallet="dlccli wallets create"

echo "Creating Alice's wallet"
$create_wallet $dlc_params $alice_params
echo "Creating Bob's wallet"
$create_wallet $dlc_params $bob_params
