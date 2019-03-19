run_bitcoind:
	./scripts/run_bitcoind.sh

stop_bitcoind:
	./scripts/stop_bitcoind.sh

clean_bitcoind:
	./scripts/stop_bitcoind.sh &> /dev/null || true
	rm -rf ./bitcoind/regtest

generate_mocks:
	./scripts/generate_mocks.sh

cli:
	dep ensure
	go install ./cmd/dlccli
