# Discreet Log Contracts (DLC)
[![CircleCI](https://circleci.com/gh/dgarage/dlc.svg?style=svg)](https://circleci.com/gh/dgarage/dlc)

Discreet Log Contracts (DLC) are smart contracts proposed by Thaddeus Dryja in [this paper](https://adiabat.github.io/dlc.pdf), which allow you to facilitate conditional payment on Bitcoin.
This library is an experimental implementation of DLC, aimed to be used in the Bitcoin mainnet.

## Setup project

### Install dependencies

```
dep ensure
```

### Run test

```
go test ./...
```

## Project Layout

```
.
├── README.md
├── Gopkg.toml
├── Gopkg.lock
├── pkg // Library code that's ok to use by external applications
├── internal // Privade code that you don't want external applications importing
└── vendor
```

## License
[MIT License](https://github.com/dgarage/dlc/blob/master/LICENSE)
