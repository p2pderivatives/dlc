#!/bin/bash -euC

mockery=$(command -v mockery)

$mockery -dir=./pkg/wallet -name=Wallet -outpkg=walletmock -output=./internal/mocks/walletmock
$mockery -dir=./internal/rpc -name=Client  -outpkg=rpcmock -output=./internal/mocks/rpcmock
