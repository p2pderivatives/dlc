# How to Create DLC using CLI

## Setup
### Prerequisites
* bitciond / bitcoin-cli have been installed
	* e.g. `brew install bitcoin` 
* Contract conditions and deals have been decided

### Install CLI
```
go get -u github.com/p2pderivatives/dlc/cmd/dlccli
dlccli --help
```

and go to source directory

```
cd $GOPATH/src/github.com/p2pderivatives/dlc
```

### Run bitcoind
```
# regtest mode
make run_bitcoind 

# testnet mode
make run_bitcoind BITCOIN_NET=testnet

# mainnet mode
make run_bitcoind BITCOIN_NET=mainnet
```

To clean up regtest, run the following tasks and run bitciond again

```
make stop_bitcoind && make clean_bitciond
```

if you want to stop bitconind of mainnet/testnet

```
make stop_bitcoind BITCOIN_NET=testnet # or mainnet
```


## Steps
* [Create Wallet](#create-wallet)
* [Deposite funds and fees](#deposite-funds-and-fees)
* [Create DLC](#create-dlc)
* [Confirm Created Transactions](#confirm-created-transactions)
* [Send Fund Tx](#send-fund-tx)
* [Fix Message](#fix-message)
* [Execute Contract](#execute-contract)

### Create Wallet
Alice

```
dlccli wallets seed --conf ./conf/bitcoin.regtest.conf

3100c6bd2a1fe93baaa49e7666bd5c3875a86fb7e417947281d757e7f8e6593f
```

```
dlccli wallets create \
    --conf ./conf/bitcoin.regtest.conf \
    --walletdir ./wallets/regtest \
    --walletname "alice" \
    --privpass "priv_alice" \
    --pubpass "pub_alice" \
    --seed "seed_alice" 

Wallet created
```

Bob

```
dlccli wallets create \
    --conf ./conf/bitcoin.regtest.conf \
    --walletdir ./wallets/regtest \
    --walletname "bob" \
    --privpass "priv_bob" \
    --pubpass "pub_bob" \
    --seed "seed_bob" 

Wallet created
```

or run 

```
# regtest
./test/cmd/create_wallets.sh

# testnet
BITCOIN_NET=testnet ./test/cmd/create_wallets.sh

# mainet
BITCOIN_NET=mainnet ./test/cmd/create_wallets.sh
```

### Deposite Funds and Fees
Alice

```
dlccli wallets addresses create \
    --conf ./conf/bitcoin.regtest.conf \
    --walletdir ./wallets/regtest \
    --walletname "alice" \
    --pubpass "pub_alice"

bcrt1q9679haanl0tax3wmylsdr62ft3xfc2yu9g74a4
```
```
% bitcoin-cli -conf=`pwd`/conf/bitcoin.regtest.conf \
	sendtoaddress bcrt1q9679haanl0tax3wmylsdr62ft3xfc2yu9g74a4 0.20022035
	
b4586b8acc967292f760c6217af91d7500cae217a67ec2e890a6e06913cb7986
```
```
% bitcoin-cli -conf=`pwd`/conf/bitcoin.regtest.conf generate 1

[
  "5d9aef6cae467233c6d263d85875e1db10b5dfc561e03b8c2d4924b2a5d7bd32"
]
```
```
dlccli wallets balance \
	--conf ./conf/bitcoin.regtest.conf \
	--walletdir ./wallets/regtest \
	--walletname "alice" \
	--pubpass "pub_alice"
	
0.20022035
```

Bob

```
dlccli wallets addresses create \
    --conf ./conf/bitcoin.regtest.conf \
    --walletdir ./wallets/regtest \
    --walletname "bob" \
    --pubpass "pub_bob"

bcrt1q08q2ks538wqzu2z766rpfh0cttvckz6jsluszm
```
```
bitcoin-cli -conf=`pwd`/conf/bitcoin.regtest.conf \
	sendtoaddress bcrt1q08q2ks538wqzu2z766rpfh0cttvckz6jsluszm 0.33355368
	
4181f0d52966edf374b7319d3975c41b61bdd53861b319a0e836f9c8df6f5e01
```
```
bitcoin-cli -conf=`pwd`/conf/bitcoin.regtest.conf generate 1

[
  "6f176d42b59de844d0f334ac06aafb348c20fec68ddbc0d6034de94f09fe9d71"
]
```
```
dlccli wallets balance \
	--conf ./conf/bitcoin.regtest.conf \
	--walletdir ./wallets/regtest \
	--walletname "bob" \
	--pubpass "pub_bob"
	
0.33355368
```

or run

```
# regtest
./test/cmd/deposit_funds.sh
./test/cmd/check_balances.sh


# testnet
BITCOIN_NET=testnet ./test/cmd/deposit_funds.sh
BITCOIN_NET=testnet ./test/cmd/check_balances.sh

# mainet
BITCOIN_NET=mainnet ./test/cmd/deposit_funds.sh
BITCOIN_NET=mainnet ./test/cmd/check_balances.sh
```

### Prepare Deals

Prepare deals in csv format

```
# value, distribution1(satoshi), distribution2(satoshi)
% head ./test/cmd/deals.csv
30001,11110,98028100
30002,22220,98016990
30003,33330,98005880
30004,44430,97994770
30005,55540,97983660
30006,66650,97972560
30007,77750,97961450
30008,88860,97950350
30009,99970,97939240
30010,111070,97928140
```

### Create DLC
* Alice creates 2 addresses (transfer and change)
* Bob creates 2 addresses (transfer and change)
* Alice or Bob gets Oracle's pubkey

```
% dlccli oracle rpoints \
    --conf ./conf/bitcoin.regtest.conf \
    --oraclename "olivia" \
    --rpoints 4 \
    --fixingtime "2019-03-30T12:00:00Z" \
> opub.json && cat opub.json

{
  "pubkey": "0268cc166f01171b411e0929243a26e40eb43cc0471b454624bdf76c6b5b2e678d",
  "rpoints": [
    "0246920d67ca40fbdad84f48a7a23ac763c45d3d990fef26a1b55df6a4b1b06098",
    "020de6c68a99ce542bbd256070e0b00b8529ac7037b2ee72929b89cc0efc51f1ab",
    "024afa3590955cb0d84b3f3a8a0ef2f419565d83a3d117a1421ed775ccaea6e716",
    "024744439b417fb75f03b623dd02f5ffc6f4d1e9fe29e10ab1b8b25181e92baaae",
    "02281c1b8e86f4845125584325f5cb3064b0bbf0aa6019a49d14833bc20219eeef"
  ]
}
```

* Alice and Bob create DLC

```
dlccli contracts create \
	--conf ./bitcoind/bitcoin.regtest.conf \
	--oracle_pubkey ./opub.json \
	--fixingtime "2019-03-30T12:00:00Z" \
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

or run

```
# regtest
./test/cmd/create_dlc.sh

# testnet
BITCOIN_NET=testnet ./test/cmd/create_dlc.sh

# mainet
BITCOIN_NET=mainnet ./test/cmd/create_dlc.sh
```


### Confirm Created Transactions

Fund Tx

```
% bitcoin-cli -conf=`pwd`/conf/bitcoin.regtest.conf decoderawtransaction \
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

```
% bitcoin-cli -conf=`pwd`/conf/bitcoin.regtest.conf decoderawtransaction \
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

```
% bitcoin-cli -conf=`pwd`/conf/bitcoin.regtest.conf sendrawtransaction \
020000000001028679cb1369e0a690e8c27ea617e2ca00751df97a21c660f7927296cc8a6b58b40000000000ffffffff015e6fdfc8f936e8a019b36138d5bd611bc475399d31b774f3ed6629d5f081410000000000ffffffff03972bd805000000002200208fc69883a973c042c3d636ba6bf16c6d6fcb8e0645b17c49593d4a81b032c0e351d8c338000000001600149f96e83381d6cdd2547c4ba356799f3dd98e8aa5701a993800000000160014b2148ad2f69392741d6aea79c3e9a5b66a95527f02483045022100960ec31876d5250b5c06ec674149ba5229df35553b9dff901195bde65eb572ea02202d2f2f1ef937be7c654d4dfb9942adcf658111adf6b84d9b84d5465ba5d2a2bc012103cb7a859433bbda5ab9d5a956e05e3c0bc32c93ef89333354dd26ab0ff5a4771602473044022043d3f1f92a68ba7cd79903e150502452a5ea7100c1e70f527b649479f070c0e902202588b70052348e5180bcc0491e4b4dc80f8001a7a176a3d7af09e2d98dfad638012103eb9a227fbae493c2d310cccdd9151d0b372be81ccc9665842279c17fd8cf50a800000000

50acb1c622218daf35919cbd667843ec7c070c43318409b43eb58d751e2aa874
```


### Fix Message

```
% dlccli oracle messages fix \
	--conf ./conf/bitcoin.regtest.conf \
	--oraclename "olivia" \
	--rpoints 5 \
	--fixingtime "2019-03-30T12:00:00Z" \
	--fixingvalue 30099 \
> osig.json && cat osig.json

{
  "sigs": [
    "202654f589c18a0672fc2501af8cc70be7bb5d2afdf40c8e3a7f41981df4e850",
    "409b0f0a358f83fef5ebf546e549b4fb6634e2e4ced96e00336bf99249cb8d86",
    "fdc0ab07f575b4d6e4c7dbf3ad6d2cf853f26cad69a4c1374895c63a19907112",
    "03fc4728ab8268c97a4aaf320bad3f78d6ee67ad42e37f861b759d5d6173a0ab",
    "fb9de05bf985a95ed99199004cb54bcc5b1504c41f5e77dd477bd662c30464e9"
  ],
  "value": 30099
}
```

or run

```
# regtest
./test/cmd/fix_message.sh 3500

# testnet
BITCOIN_NET=testnet ./test/cmd/fix_message.sh 3500

# mainet
BITCOIN_NET=mainnet ./test/cmd/fix_message.sh 3500
```

### Execute Contract
 
```
% dlccli contracts deals fix \
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

or run

```
# regtest
./test/cmd/fix_deal.sh 68a0c4026c76800c33bd5614fec7b3402bf55067dc2670576f146ac26a98b692

# testnet
BITCOIN_NET=testnet ./test/cmd/fix_deal.sh 68a0c4026c76800c33bd5614fec7b3402bf55067dc2670576f146ac26a98b692

# mainet
BITCOIN_NET=mainnet ./test/cmd/fix_deal.sh 68a0c4026c76800c33bd5614fec7b3402bf55067dc2670576f146ac26a98b692
```

send the created CETx and ClosingTx to the network using bitcoin-cli as did for fund tx
