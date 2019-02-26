run_bitcoind:
	./scripts/run_bitcoind.sh

stop_bitcoind:
	./scripts/stop_bitcoind.sh

generate_mocks:
	./scripts/generate_mocks.sh

demo_cli:
	go build -o bin/dlcdemo-cli ./cmd/dlcdemo
