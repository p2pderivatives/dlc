# wallet

This package provides a secure(?) hierarchical deterministic wallet. The hierarchical deterministic address manager is provided by btcsuite's [waddrmgr](http://godoc.org/github.com/btcsuite/btcwallet/waddrmgr).

## Assumptions

Right now, this library assumes only one cointype will be used (`waddrmgr.KeyScopeBIP0084`) and there will be only one account associated with that cointype. If this is to change, this library will need to be refactored a little. 

## How address management works

When a wallet is created, `Manager` is created as part of the wallet. `Manager` is the root manager; it handles the root HD key (m/). A `ScopedKeyManager` is a sub key manager under the main root key manager; each scoped key managagers handles the cointype key for a particular key scope (m/purpose/cointype). 

Under each `ScopedKeyManager` are `Account` types associated with that `ScopedKeyManager`.

For more information on address management, please consult the original [godoc](https://godoc.org/github.com/btcsuite/btcwallet/waddrmgr).

## How UTXO management works

Right now, to ask about UTXOS, the wallet will query the running `bitcoind` instance by using the rpc command `ListUnspent()`.
When a public address is generated, it is also registered to `bitcoind`, so `bitcoind` knows to keep track of transactions associated with that addresss.