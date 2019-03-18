#!/bin/bash

net=${BITCOIN_NET:=regtest}
conf="bitcoin.${net}.conf"
dlc_params="--conf ./conf/${conf} --walletdir ./wallets/${net}"
balance="dlccli wallets balance $dlc_params"
alice_params="--walletname alice --pubpass pub_alice"
alicep_params="--walletname alicep --pubpass pub_alicep"
bob_params="--walletname bob --pubpass pub_bob"
bobp_params="--walletname bobp --pubpass pub_bobp"

echo "Alice (DLC Wallet): $($balance $alice_params)"
echo "Alice (Personal Wallet): $($balance $alicep_params)"
echo "Bob (DLC Wallet): $($balance $bob_params)"
echo "Bob (Personal Wallet): $($balance $bobp_params)"
