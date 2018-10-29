# wallet

This package provides a secure(?) hierarchical deterministic wallet. The hierarchical deterministic address manager is provided by btcsuite's [waddrmgr](http://godoc.org/github.com/btcsuite/btcwallet/waddrmgr). 

## How address management works

When a wallet is created, `Manager` is created as part of the wallet. A `ScopedKeyManager` is a sub key manager under the main root key manager; each scoped key managagers handles the cointype key for a particular key scope (m/purpose/cointype)

