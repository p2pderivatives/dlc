# Common settings
## The network on which to execute the DLC (regtest, testnet or mainnet)
export BITCOIN_NET="regtest"
## The amount of premium to be included in the fund transaction
export PREMIUM_AMOUNT=200000
## The party that will pay the premium to the other (0 if party 1 is paying or 1 if party 2 is paying)
export PREMIUM_PAYING_PARTY=1
## The time at which the contract matures
export FIX_TIME=$(date -u -v+10M +%FT%TZ)
## The fee rate to use for the fund transaction
export FUND_TX_FEERATE=50
## The fee rate to use for the redeem transaction
export REDEEM_TX_FEERATE=40
## The CSV file containing the deal
export DEALS_FILE="./test/cmd/deals.csv"
## The time after which the refund transaction can be used
export REFUND_LOCKTIME="574196"
## The path where to store the wallet
export WALLET_DIR="./wallets/$BITCOIN_NET"

# Settings for party 1
## The name of the wallet
export PARTY1_WALLET_NAME="alice"
## The private key pass for the wallet
export PARTY1_PRIV_PASS="priv_alice"
## The public key pass for the wallet
export PARTY1_PUB_PASS="pub_alice"
## The fund input by party 1 in the contract
export PARTY1_FUND=20000000

# Settings for party 2
## The name of the wallet
export PARTY2_WALLET_NAME="bob"
## The private key pass for the wallet
export PARTY2_PRIV_PASS="priv_bob"
## The public key pass for the wallet
export PARTY2_PUB_PASS="pub_bob"
## The fund input by party 2 in the contract
export PARTY2_FUND=33333333

#Settings for oracle
## The number of R points to use (number of digit to be signed)
export RPOINTS=4
## The name of the oracle
export ORACLE_NAME="olivia"
