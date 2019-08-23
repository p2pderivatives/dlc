# How to Create DLC using CLI

## Setup
### Prerequisites
* bitcoind / bitcoin-cli have been installed
	* e.g. `brew install bitcoin` 
* Contract conditions and deals have been decided

### Install CLI
```bash
go get -u github.com/p2pderivatives/dlc/cmd/dlccli
dlccli --help
```

and go to source directory

```bash
cd $GOPATH/src/github.com/p2pderivatives/dlc
```

### Run bitcoind
```bash
# regtest mode
make run_bitcoind 

# testnet mode
make run_bitcoind BITCOIN_NET=testnet

# mainnet mode
make run_bitcoind BITCOIN_NET=mainnet
```

To clean up regtest, run the following tasks and run bitcoind again

```bash
make stop_bitcoind && make clean_bitcoind
```

if you want to stop bitconind of mainnet/testnet

```bash
make stop_bitcoind BITCOIN_NET=testnet # or mainnet
```


## Steps
* [Set Up Parameters (optional)](#set-up-parameters-optional)
* [Create Wallet](#create-wallet)
* [Create Addresses](#create-addresses)
* [Deposite funds and fees](#deposite-funds-and-fees)
* [Create DLC](#create-dlc)
* [Confirm Created Transactions](#confirm-created-transactions)
* [Send Fund Tx](#send-fund-tx)
* [Fix Message](#fix-message)
* [Execute Contract](#execute-contract)

### Set Up Parameters (optional)

If you wish to use the scripts to execute the different steps, you first need to set up the parameters.
Edit the file `test/cmd/set_parameters.sh` and then source it using:

```bash
source ./test/cmd/set_parameters.sh
```
The parameters are described in the file.

(Note: it is important to use `source` and not execute it directly otherwise the environment variables won't be set.)

### Create Wallet

#### Using a script

Run:
```bash
./test/cmd/create_wallets.sh
```

#### Using commands

##### Create Alice's wallet

```bash
dlccli wallets seed --conf ./conf/bitcoin.regtest.conf

3100c6bd2a1fe93baaa49e7666bd5c3875a86fb7e417947281d757e7f8e6593f
```
Hereafter, replace "seed\_alice" with the hexadecimal number generated above.

```bash
dlccli wallets create \
    --conf ./conf/bitcoin.regtest.conf \
    --walletdir ./wallets/regtest \
    --walletname "alice" \
    --privpass "priv_alice" \
    --pubpass "pub_alice" \
    --seed "seed_alice" 
```

#### Create Bob's wallet

```bash
dlccli wallets seed --conf ./conf/bitcoin.regtest.conf

3100c6bd2a1fe93baaa49e7666bd5c3875a86fb7e417947281d757e7f8e6593f

```
Hereafter, replace "seed\_bob" with the hexadecimal number generated above.

```bash
dlccli wallets create \
    --conf ./conf/bitcoin.regtest.conf \
    --walletdir ./wallets/regtest \
    --walletname "bob" \
    --privpass "priv_bob" \
    --pubpass "pub_bob" \
    --seed "seed_bob" 
```

### Create addresses

#### Using a script

Run:
```bash
source test/cmd/create_addresses.sh
```
(Note: again be careful to use `source` so that environment variables are properly set.)

#### Using commands

The DLC requires three addresses:
* One base address to be used in the fund transaction,
* One change address to be use to send the change from the fund transaction,
* One transfer address to receive the output of the DLC.

For Alice use the following command to generate an address:

```bash
dlccli wallets addresses create \
    --conf ./conf/bitcoin.regtest.conf \
    --walletdir ./wallets/regtest \
    --walletname "alice" \
    --pubpass "pub_alice"

bcrt1q9679haanl0tax3wmylsdr62ft3xfc2yu9g74a4
```

For Bob:

```bash
dlccli wallets addresses create \
    --conf ./conf/bitcoin.regtest.conf \
    --walletdir ./wallets/regtest \
    --walletname "bob" \
    --pubpass "pub_bob"

bcrt1q08q2ks538wqzu2z766rpfh0cttvckz6jsluszm
```

### Deposite Funds and Fees

Note that the below steps are only valid for `regtest`. For `testnet` and `mainnet`, you will have to send the fund using your wallet to the base addresses generated in the previous step.

#### Using a script

Run:
```bash
./test/cmd/deposit_funds.sh
./test/cmd/check_balances.sh
```

#### Using commands

Alice

```bash
% bitcoin-cli -conf=`pwd`/conf/bitcoin.regtest.conf \
	sendtoaddress bcrt1q9679haanl0tax3wmylsdr62ft3xfc2yu9g74a4 0.20022035
	
b4586b8acc967292f760c6217af91d7500cae217a67ec2e890a6e06913cb7986
```

If using `regtest`:
```bash
% bitcoin-cli -conf=`pwd`/conf/bitcoin.regtest.conf generate 1

[
  "5d9aef6cae467233c6d263d85875e1db10b5dfc561e03b8c2d4924b2a5d7bd32"
]
```
```bash
dlccli wallets balance \
	--conf ./conf/bitcoin.regtest.conf \
	--walletdir ./wallets/regtest \
	--walletname "alice" \
	--pubpass "pub_alice"
	
0.20022035
```

Bob

```bash
bitcoin-cli -conf=`pwd`/conf/bitcoin.regtest.conf \
	sendtoaddress bcrt1q08q2ks538wqzu2z766rpfh0cttvckz6jsluszm 0.33355368
	
4181f0d52966edf374b7319d3975c41b61bdd53861b319a0e836f9c8df6f5e01
```

If using `regtest`:
```bash
bitcoin-cli -conf=`pwd`/conf/bitcoin.regtest.conf generate 1

[
  "6f176d42b59de844d0f334ac06aafb348c20fec68ddbc0d6034de94f09fe9d71"
]
```
```bash
dlccli wallets balance \
	--conf ./conf/bitcoin.regtest.conf \
	--walletdir ./wallets/regtest \
	--walletname "bob" \
	--pubpass "pub_bob"
	
0.33355368
```

### Prepare Deals

Prepare deals in csv format (or use the file included in the repository).

Example:
```bash
# value, distribution1(satoshi), distribution2(satoshi)
$ head ./test/cmd/deals.csv
3000,53333333,0
3001,53288903,44429
3002,53244503,88829
3003,53200133,133200
3004,53155792,177541
3005,53111480,221852
3006,53067198,266134
3007,53022946,310386
3008,52978723,354609
3009,52934529,398803
```

### Create DLC

#### Using a script

Run

```bash
./test/cmd/create_dlc.sh
```

#### Using commands

##### Generate R points for oracle

(Note: fix time needs to be greater than current time.)

```bash
$ dlccli oracle rpoints \
    --conf ./conf/bitcoin.regtest.conf \
    --oraclename "olivia" \
    --rpoints 4 \
    --fixingtime "2019-08-30T12:00:00Z" \
> opub.json && cat opub.json
{
  "pubkey": "03a7844731daf02e1fa81249bafe7efe456375ebbcce927dd38e4e4232853dff15",
  "rpoints": [
    "021fc465868ef187ec26ce45407bfadd99839d4dc101f4c4526e2e9a17ad3e811b",
    "02dc200c6cadcf320abbaeabb710a11d1da2e92cb42a1d7d7b19729a270ed51190",
    "03602cf96e40bc6fd08df540cbcc841cce95625bd9d4588d7e8bc113d864c6a5a0",
    "03e2ee1312646d0b22d526299ca293445ff4c242b0cae0aadd76fa1ccb0bd06880"]
}
```

##### Alice and Bob create DLC

(Note: hereafter `address1` and `address2` correspond to the transfer addresses generated in [create addresses](#create-addresses))

```bash
dlccli contracts create \
	--conf ./bitcoind/bitcoin.regtest.conf \
	--oracle_pubkey ./opub.json \
	--fixingtime "2019-08-30T12:00:00Z" \
	--fund1 2000000 \
	--fund2 3333333 \
	--address1 "bcrt1qjndkjszkzqpahdzz8kkc4hgxljlrlswp25cusr" \
	--address2 "bcrt1qkcuc77z0ktnyfv8stz6xrnnc8uawxt7gg055gt" \
	--change_address1 "bcrt1qd9aadr8jf4v2y4pe0l239h2r25tmff9lrz35v9" \
	--change_address2 "bcrt1qwk5j2evm0h3kf7cakd4dlh0e9v38ll0ex64ssj" \
	--fundtx_feerate 50 \
	--redeemtx_feerate 40 \
	--deals_file ./test/cmd/deals.csv \
	--refund_locktime 5000 \
	--walletdir ./wallets/regtest \
	--wallet1 "alice" \
	--wallet2 "bob" \
	--pubpass1 "pub_alice" \
	--pubpass2 "pub_bob" \
	--privpass1 "priv_alice" \
	--privpass2 "priv_bob"
	
Contract created

ContractID:
68a0c4026c76800c33bd5614fec7b3402bf55067dc2670576f146ac26a98b692

FundTx hex:
020000000001028679cb1369e0a690e8c27ea617e2ca00751df97a21c660f7927296cc8a6b58b40000000000ffffffff015e6fdfc8f936e8a019b36138d5bd611bc475399d31b774f3ed6629d5f081410000000000ffffffff03972bd805000000002200208fc69883a973c042c3d636ba6bf16c6d6fcb8e0645b17c49593d4a81b032c0e351d8c338000000001600149f96e83381d6cdd2547c4ba356799f3dd98e8aa5701a993800000000160014b2148ad2f69392741d6aea79c3e9a5b66a95527f02483045022100960ec31876d5250b5c06ec674149ba5229df35553b9dff901195bde65eb572ea02202d2f2f1ef937be7c654d4dfb9942adcf658111adf6b84d9b84d5465ba5d2a2bc012103cb7a859433bbda5ab9d5a956e05e3c0bc32c93ef89333354dd26ab0ff5a4771602473044022043d3f1f92a68ba7cd79903e150502452a5ea7100c1e70f527b649479f070c0e902202588b70052348e5180bcc0491e4b4dc80f8001a7a176a3d7af09e2d98dfad638012103eb9a227fbae493c2d310cccdd9151d0b372be81ccc9665842279c17fd8cf50a800000000

RefundTx hex:
0200000000010174a82a1e758db53eb4098431430c077cec437866bd9c9135af8d2122c6b1ac500000000000feffffff02e79bd60200000000160014498ef5355b87b188ff074b41ab36c9840171a0c3c859010300000000160014dfea60b1b9eeea759f0c988a89c854ad6959f7930400483045022100d81c1ea006d5f1c132339e0f789756ee354ed906a877d5169a81c44d0d36cb5c02207987d8c5433af70d830009b1213ebd20f382760207e0f702cf1a6281c4d811810147304402205b1fa0062f6c1c2993a6953bed4356e788a2fd69d2a3e05d34212496f2d9257f022034bf8abe5994b5ce8d61f91ea661aca58dcf679272c8a3aaa408638083105f010147522102923081bc3ff2e7c72969013265e99b1c4af1c59b6572f6f9734f632aeb7d319821037b8e3201837b670bafaf3fd3e5cb20e952e45cc406040eda27c32f192dad968c52ae88130000
```

If one of the party pays a premium to the other use:

```bash
dlccli contracts createwithpremium ... --premiumamount 200000 --premiumdestaddress "bcrt1qvqpyv5nd49y9y7ndjdum5xy0qduttc7ehhw70e" --premiumpayingparty 0
```
(Note: the other parameters should be set similarly.)

### Confirm Created Transactions

Fund Tx

```bash
$ bitcoin-cli -conf=`pwd`/conf/bitcoin.regtest.conf decoderawtransaction \
020000000001028679cb1369e0a690e8c27ea617e2ca00751df97a21c660f7927296cc8a6b58b40000000000ffffffff015e6fdfc8f936e8a019b36138d5bd611bc475399d31b774f3ed6629d5f081410000000000ffffffff03972bd805000000002200208fc69883a973c042c3d636ba6bf16c6d6fcb8e0645b17c49593d4a81b032c0e351d8c338000000001600149f96e83381d6cdd2547c4ba356799f3dd98e8aa5701a993800000000160014b2148ad2f69392741d6aea79c3e9a5b66a95527f02483045022100960ec31876d5250b5c06ec674149ba5229df35553b9dff901195bde65eb572ea02202d2f2f1ef937be7c654d4dfb9942adcf658111adf6b84d9b84d5465ba5d2a2bc012103cb7a859433bbda5ab9d5a956e05e3c0bc32c93ef89333354dd26ab0ff5a4771602473044022043d3f1f92a68ba7cd79903e150502452a5ea7100c1e70f527b649479f070c0e902202588b70052348e5180bcc0491e4b4dc80f8001a7a176a3d7af09e2d98dfad638012103eb9a227fbae493c2d310cccdd9151d0b372be81ccc9665842279c17fd8cf50a800000000

{
  "txid": "29e7b2e17ebc192fc72be08dbab1d618a72dd4e1547d499c68ac0699b427ff6d",
  "hash": "9f9d293e32816cf3f0d8410dc65c490bb03488b26106237ece39d54d89f8f5b6",
  "version": 2,
  "size": 413,
  "vsize": 251,
  ...
  "vout": [
    {
      "value": 0.98053015,
      "n": 0,
      "scriptPubKey": {
        "asm": "0 c51b2215f583b9544186a9352b2e800860a4a20ba75dabd4ac5554ab7736db64",
        "hex": "0020c51b2215f583b9544186a9352b2e800860a4a20ba75dabd4ac5554ab7736db64",
        "reqSigs": 1,
        "type": "witness_v0_scripthash",
        "addresses": [
          "bcrt1qc5djy904swu4gsvx4y6jkt5qpps2fgst5aw6h49v2422kaekmdjqnrk7rs"
        ]
      }
    },
    ...
  ]
}
```

Refund Tx

```bash
$ bitcoin-cli -conf=`pwd`/conf/bitcoin.regtest.conf decoderawtransaction \
0200000000010174a82a1e758db53eb4098431430c077cec437866bd9c9135af8d2122c6b1ac500000000000feffffff02e79bd60200000000160014498ef5355b87b188ff074b41ab36c9840171a0c3c859010300000000160014dfea60b1b9eeea759f0c988a89c854ad6959f7930400483045022100d81c1ea006d5f1c132339e0f789756ee354ed906a877d5169a81c44d0d36cb5c02207987d8c5433af70d830009b1213ebd20f382760207e0f702cf1a6281c4d811810147304402205b1fa0062f6c1c2993a6953bed4356e788a2fd69d2a3e05d34212496f2d9257f022034bf8abe5994b5ce8d61f91ea661aca58dcf679272c8a3aaa408638083105f010147522102923081bc3ff2e7c72969013265e99b1c4af1c59b6572f6f9734f632aeb7d319821037b8e3201837b670bafaf3fd3e5cb20e952e45cc406040eda27c32f192dad968c52ae88130000

{
  ...
  "vin": [
    {
      "txid": "b7af665256972e1ca53b2ebe13aa471c73aa979fb9be4f42575523ea8444e295",
      "vout": 0,
      "scriptSig": {
        "asm": "",
        "hex": ""
      },
      "txinwitness": [
        "",
        "3044022031b478fcf34e650b652f8f8cccc81be9593e4f5a0a275a3c69c44cf6a8f76f9502206fe697675c50daffd45fda34c9c928b98d9f6f84137ac751e89ce667a3c9642c01",
        "304402203d669277cd5d4ad4be11b5dd10e89af4245e0dc84bc0df71993c99f0b08e62d6022004e18ff9187d0e1c207dafc254f9aef2aef0dd8b771ef6386513f7950475684801",
        "522103b1afa38b000ecf30b589d401fe052225f275719470b058eb8885d6509c9140b82103bf1f253c953480cf4719d30049d053a1bf61fb4e368bf198f1dc50b7f176da1e52ae"
      ],
      "sequence": 4294967294
    }
  ],
  "vout": [
    {
      "value": 0.47619047,
      "n": 0,
      "scriptPubKey": {
        "asm": "0 6797b2e46b6fc41276fe117c8b1ad863f7208975",
        "hex": "00146797b2e46b6fc41276fe117c8b1ad863f7208975",
        "reqSigs": 1,
        "type": "witness_v0_keyhash",
        "addresses": [
          "bcrt1qv7tm9ertdlzpyah7z97gkxkcv0mjpzt4xp066v"
        ]
      }
    },
    {
      "value": 0.50420168,
      "n": 1,
      "scriptPubKey": {
        "asm": "0 b2d65da4e5d306588a6d9d5e55e990a9fb1d9559",
        "hex": "0014b2d65da4e5d306588a6d9d5e55e990a9fb1d9559",
        "reqSigs": 1,
        "type": "witness_v0_keyhash",
        "addresses": [
          "bcrt1qktt9mf896vr93zndn409t6vs48a3m92ev8gtdh"
        ]
      }
    }
  ]
}
```


### Send Fund Tx

Run:
```bash
$ bitcoin-cli -conf=`pwd`/conf/bitcoin.regtest.conf sendrawtransaction \
020000000001028679cb1369e0a690e8c27ea617e2ca00751df97a21c660f7927296cc8a6b58b40000000000ffffffff015e6fdfc8f936e8a019b36138d5bd611bc475399d31b774f3ed6629d5f081410000000000ffffffff03972bd805000000002200208fc69883a973c042c3d636ba6bf16c6d6fcb8e0645b17c49593d4a81b032c0e351d8c338000000001600149f96e83381d6cdd2547c4ba356799f3dd98e8aa5701a993800000000160014b2148ad2f69392741d6aea79c3e9a5b66a95527f02483045022100960ec31876d5250b5c06ec674149ba5229df35553b9dff901195bde65eb572ea02202d2f2f1ef937be7c654d4dfb9942adcf658111adf6b84d9b84d5465ba5d2a2bc012103cb7a859433bbda5ab9d5a956e05e3c0bc32c93ef89333354dd26ab0ff5a4771602473044022043d3f1f92a68ba7cd79903e150502452a5ea7100c1e70f527b649479f070c0e902202588b70052348e5180bcc0491e4b4dc80f8001a7a176a3d7af09e2d98dfad638012103eb9a227fbae493c2d310cccdd9151d0b372be81ccc9665842279c17fd8cf50a800000000

50acb1c622218daf35919cbd667843ec7c070c43318409b43eb58d751e2aa874
```
If in `regtest`, generate a block so that the transaction gets into the blockchain.
```bash
$ bitcoin-cli --conf=`pwd`/conf/bitcoin.regtest.conf generate 1
[
  "3394576eaf0d1fcc8158d03b7e6f72065acc235ecc67790250e895ca27f07d2f"
]
```

### Fix Message

#### Using a script

Run:

```bash
./test/cmd/fix_message.sh 3500
```

#### Using commands

```bash
$ dlccli oracle messages fix \
	--conf ./conf/bitcoin.regtest.conf \
	--oraclename "olivia" \
	--rpoints 4 \
	--fixingtime "2019-08-30T12:00:00Z" \
	--fixingvalue 3500 \
> osig.json && cat osig.json

{
  "sigs": [
    "202654f589c18a0672fc2501af8cc70be7bb5d2afdf40c8e3a7f41981df4e850",
    "409b0f0a358f83fef5ebf546e549b4fb6634e2e4ced96e00336bf99249cb8d86",
    "fdc0ab07f575b4d6e4c7dbf3ad6d2cf853f26cad69a4c1374895c63a19907112",
    "03fc4728ab8268c97a4aaf320bad3f78d6ee67ad42e37f861b759d5d6173a0ab",
    "fb9de05bf985a95ed99199004cb54bcc5b1504c41f5e77dd477bd662c30464e9"
  ],
  "value": 3500
}
```

### Execute Contract

#### Using a script

Run:

```bash
./test/cmd/fix_deal.sh 68a0c4026c76800c33bd5614fec7b3402bf55067dc2670576f146ac26a98b692
```

where the parameter is the id of the contract generated in [create dlc](create-dlc).
 
#### Using commands

```bash
$ dlccli contracts deals fix \
	--conf ./conf/bitcoin.regtest.conf \
	--dlcid 68a0c4026c76800c33bd5614fec7b3402bf55067dc2670576f146ac26a98b692
	--oracle_sig ./osig.json \
	--walletdir ./wallets/regtest \
	--wallet alice \
	--pubpass pub_alice \
	--privpass priv_alice \
	--contractor_type 0

CETx hex
020000000001016eb5564574637cae3470e27c68990293605dc9447638ebd3a40faab0180958dc0000000000ffffffff02662b000000000000220020d741420f02b68e17b0a0d6ff2d8be40f10921a13d887319d9c215b82970b4a8744cad70500000000160014548ac7261ac7b949c9bf51b9fadf3cb495c50149040047304402203f881fedff88e71e2d9857260a58cdd8ff8a9b7c440a5b89c429bca1240a8e3502204c991f90cdcefce1cb8ba96aedaf290b16e564e350ebe0427f8ab0af00b8848c01483045022100cc4c43babc78f7893fa2648b8534e6b9dd0b683d69d8b975186913158f850f6f022076c8abd6c691258c54fc12ac83f968d1889ce70a7fa059e0850333879ffc935d0147522103fea2d25e14159fb1859834dac34cf98c01df324263cb80e304b8d4e8d0e0428421031179e5d1c0f3faad53e9d5f3760705e3d006655e3ad0a5a32afe73639dc50dab52ae00000000


ClosingTx hex
020000000001011ede1f49a28f214a69c4ba7da872c14f55aca6b281f703b2fcc45eee3a55abae0000000000ffffffff01a60900000000000016001480ea81624d813b845b4d61c562c7a3f555be51bc03483045022100c70fd68cb4fba78865d27c83976e1f31b8639808d9e17d0465c3032065c66485022040c795dea6c3b6802c4112436971d5a672640646bdbee3a1928bc2b1e9b7bf070101014d632103d554fe2c8c425ee13e91a06ea4d47daae597c7eebab643f178147738c2c2792767029000b27521031179e5d1c0f3faad53e9d5f3760705e3d006655e3ad0a5a32afe73639dc50dab68ac00000000
```

Finally send the created CETx and ClosingTx to the network using bitcoin-cli as it was done in [send fund transaction](#send-fund-tx).
