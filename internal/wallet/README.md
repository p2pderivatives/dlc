# wallet

This package provides wallet features necessary for implementing DLCs.

This package heavily relies on some of the packages of [`btcsuite/btcwallet`](🇲http://godoc.org/github.com/btcsuite/btcwallet/wallet) for its core functionality.

## Connecting to `bitcoind`

Until bitcoin rpc parameters can be read automatically from a `bitcoin.conf` file, your bitcoin rpc parameters need to be entered manually in the `wallet.go` file, in lines 59-61.
```
	rpcport     = "localhost: REPLACEME"
	rpcusername = "RENAME!"
	rpcpassword = "RENAME!"
```
For your convenience, the defualt `mainnet` port number is `8333`, for `testnet3` is `18332` , and `regnet` is `18443`.
