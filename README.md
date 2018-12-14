# Discreet Log Contracts (DLC)
[![CircleCI](https://circleci.com/gh/p2pderivatives/dlc.svg?style=svg)](https://circleci.com/gh/p2pderivatives/dlc)

Discreet Log Contracts (DLC) are smart contracts proposed by Thaddeus Dryja in [this paper](https://adiabat.github.io/dlc.pdf), which allow you to facilitate conditional payment on Bitcoin.
This library is an experimental implementation of DLC, aimed to be used in the Bitcoin mainnet, and **not ready for production**.

## Overview
This contract schema involves three parties, 2 contractors and an oracle. We call the 2-contractors Alice and Bob, and the oracle Olivia.

* Alice and Bob
  * They want to make a future contract
  * They trust a message fixed and signed by Olivia
* Olivia
  * Olivia publishes messages and signs of them

You can read more details in [this article](https://medium.com/@gertjaap/discreet-log-contracts-invisible-smart-contracts-on-the-bitcoin-blockchain-cc8afbdbf0db) or [the paper](this paper](https://adiabat.github.io/dlc.pdf) for more details.

## Examples
### Communication between 3 parties
[This sample code](https://github.com/p2pderivatives/dlc/blob/master/test/integration/dlc_test.go) demonstrates how 3 parties communicate. An oracle Olivia publishes a n-digit number lottery result everyday, and Alice and Bob bet on the lottery.

In the tese, the following scenarios are tested.

* Alice and Bob make contracts and execute a fixed one.
* Oracle does't publish a valid message and sign, and contractors refund their funding transactions.

### Communication between a contractor and an oracle
[This sample code](https://github.com/p2pderivatives/dlc/blob/master/test/integration/oracle_test.go) demonstrates how a contractor Alice communicates with an oracle Olivia. 
Olivia publishes a various weather information, but Alice uses only some of the info. Contractors can choose which messages to use, and oracle doesn't know which messages are used in contracts (conditions of contracts).


## Development

### Install dependencies

```
dep ensure
```

### Run test

```
go test ./...
```

## License
[MIT License](https://github.com/dgarage/dlc/blob/master/LICENSE)
