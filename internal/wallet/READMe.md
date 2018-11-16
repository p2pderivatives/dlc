# wallet

This package provides wallet features necessary for implementing DLCs. 

## Connecting to `bitcoind`

Until bitcoin rpc parameters can be read automatically from a `bitcoin.conf` file, your bitcoin rpc parameters need to be entered manually in the `wallet.go` file, in lines 59-61.
```
	rpcport     = "localhost: REPLACEME"
	rpcusername = "RENAME!"
	rpcpassword = "RENAME!"
```
For your convenience, the defualt `mainnet` port number is `8333`, for `testnet3` is `18332` , and `regnet` is `18443`.
