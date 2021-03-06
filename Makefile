BITCOIN_NET?=regtest
export BITCOIN_NET

run_bitcoind:
	./scripts/run_bitcoind.sh

reindex_bitciond:
	./scripts/run_bitcoind.sh --reindex

stop_bitcoind:
	./scripts/stop_bitcoind.sh

clean_bitcoind:
	./scripts/stop_bitcoind.sh &> /dev/null || true
	rm -rf ./bitcoind/regtest

mocks:
	./scripts/generate_mocks.sh

cli:
	dep ensure
	go install ./cmd/dlccli
