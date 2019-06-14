# <font color=Chocolate size=6>Nox-DAG</font>

## Abstract
Introduce how to quickly establish the Nox-DAG Network

---

## What is Nox-DAG?
Nox-DAG is a DAG network based on Nox project. So it's very powerful, and you can even connect your own mining program to it for mining, or you can connect your wallet application to it for transfer transactions.

---



# Manual
## Install
Here's how to install it quickly.
### From docker:
#### ***1.Install docker:***
* First you have to make sure that you have docker installed on your machine.And make sure you are a member of the user group docker.If you're all right, you can ignore this and jump right into ***step 2***.

* Install docker on ubuntu:
```
sudo apt-get update
sudo apt-get install \
    apt-transport-https \
    ca-certificates \
    curl \
    software-properties-common
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
sudo add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
   $(lsb_release -cs) \
   stable"
sudo apt-get update
sudo apt-get install docker-ce
```
* If you are already a root user, you can ignore next step
```
# Add docker user group and add the logged-in user to the docker user group.

sudo groupadd docker
sudo gpasswd -a $USER docker
newgrp docker
docker ps
```
* You can use `docker -v` to test whether the installation is successful or not.
Other systems platforms are similar.You can go [docker](https://www.docker.com/get-started)

#### ***2.Deployment image:***
* First we need to initialize the image, Of course, when you want the latest image, you can also use the following command.
```
docker pull halalchain/nox-dag
```
* <font color=Chocolate size=3>Finally, we can run this every time we start it up.</font>

```
docker run -it -p 18130:18130 -p 18131:18131 halalchain/nox-dag --miningaddr=[Your mining address] --addpeer=[peer1 IP:PORT] [--addpeer=[peer2 IP:PORT]]
```

---

## CLI
CLI is a toolset to interact with server by RPC.

* make alias
```shell
$ alias cli="docker run --rm halalchain/nox-dag cli"
```

* list all commands
```shell
$ cli
```

* If you have some advanced requirements, or just want to test your full node, you can use the `cli` we carefully prepared for you.
```
$ cli [commands]
```
* For example, if you want to get the current total number of blocks, you can use the following command:
```
$ cli block
```
---
* If you want to send a signed transaction, you can do this:
```
docker run --rm halalchain/nox-dag cli sendRawTx [hex]
```
## NX
NX is a toolset to assistant client-side operations.
* make alias
```shell
$ alias nx="docker run --rm halalchain/nox-dag nx"
```

* list all commands
```shell
$ nx
```

* If you want to do something, such as calculate hash, generate HD key and transaction signature etc. You can use `nx` command to do it. 
```
$ nx [commands]
```

# Experiment
firstly, we need create an account to recieve mining rewards.

## Create mining address
entropy->private key->public key->address

```shell
# STATEMENT 1: this command generates the private key which will be used to sign the transaction later
$ nx ec-new $(nx entropy) > miner_key.txt

$ nx ec-to-addr $(nx ec-to-public $(cat miner_key.txt))
TmfaGwUbZiCeqKqrXNBaK5wEUcwcqArqNaW


```

##  Run instance

We  add peers manually by specifying addpeer, we recommend adding at least two peers.

```shell
docker run -it -p 18130:18130 -p 18131:18131 halalchain/nox-dag --miningaddr=TmfaGwUbZiCeqKqrXNBaK5wEUcwcqArqNaW --addpeer=47.103.194.115:18130 --addpeer=42.51.64.58:38130
```

## mining
before we send transactions, we need to get some rewards by mining.

Note: this command is executed by CPU, so it takes more patience to get the result.

```shell
$ cli generate 1
[
  "000000407b57d8f3d27c6da281ca82a42d68abb6be7ab32f8770f1940cb55936"
]
```

get the latest block number
```shell
$ cli block
the lastet block is 300
```

Note: if there are more than one block created, then try reducing the block number to find the corresponding block

Then we inspect the block:
```shell
$ cli block 300
{
  "hash": "000000407b57d8f3d27c6da281ca82a42d68abb6be7ab32f8770f1940cb55936",
  "confirmations": 1,
  "version": 1,
  "height": 300,
  "txRoot": "7a5374dad32ebe125dc7dcdfe2f9510cb97c6b50c8a0f4435502906184810e9c",
  "transactions": [
    "cd1fb199ad3cc58d696cdd1499bc97c3fd9aac38a705db0d9bdd94d09ce1ad3e"
  ],
  "stateRoot": "0000000000000000000000000000000000000000000000000000000000000000",
  "bits": "1e0083d5",
  "difficulty": 503350229,
  "nonce": 12657176,
  "timestamp": "2019-05-12 13:30:43.0000",
  "parents": [
    "00000070b839b4216e658a12f2dfe42bb664932797d653cf3701367844c42664"
  ],
  "children": [
    "null"
  ]
}
```

## send transactions

from the block info, we could infer this only transaction cd1fb199ad3cc58d696cdd1499bc97c3fd9aac38a705db0d9bdd94d09ce1ad3e is the coinbase transaction with rewards
 
inspect the tx
```shell
$ cli tx cd1fb199ad3cc58d696cdd1499bc97c3fd9aac38a705db0d9bdd94d09ce1ad3e
{
  "hex": "01000000010000000000000000000000000000000000000000000000000000000000000000ffffffffffffffff0380b2e60e000000000000000000000000000e6a0c2c010000f646336d5a0c843f80461c86000000001976a914b65a5d4ce219772417459b40af009f0e39646a4b88ac00000000000000000100f902950000000000000000ffffffff0700002f6e6f782f",
  "hexnowit": "01000100010000000000000000000000000000000000000000000000000000000000000000ffffffffffffffff0380b2e60e000000000000000000000000000e6a0c2c010000f646336d5a0c843f80461c86000000001976a914b65a5d4ce219772417459b40af009f0e39646a4b88ac0000000000000000",
  "hexwit": "010002000100f902950000000000000000ffffffff0700002f6e6f782f",
  "txid": "cd1fb199ad3cc58d696cdd1499bc97c3fd9aac38a705db0d9bdd94d09ce1ad3e",
  "txhash": "7a5374dad32ebe125dc7dcdfe2f9510cb97c6b50c8a0f4435502906184810e9c",
  "version": 1,
  "locktime": 0,
  "expire": 0,
  "vin": [
    {
      "amountin": 2500000000,
      "blockheight": 0,
      "txindex": 4294967295,
      "coinbase": "00002f6e6f782f",
      "sequence": 4294967295
    }
  ],
  "vout": [
    {
      "amount": 250000000,
      "scriptPubKey": {
        "asm": "",
        "type": "nonstandard"
      }
    },
    {
      "amount": 0,
      "scriptPubKey": {
        "asm": "OP_RETURN 2c010000f646336d5a0c843f",
        "hex": "6a0c2c010000f646336d5a0c843f",
        "type": "nulldata"
      }
    },
    {
      "amount": 2250000000,
      "scriptPubKey": {
        "asm": "OP_DUP OP_HASH160 b65a5d4ce219772417459b40af009f0e39646a4b OP_EQUALVERIFY OP_CHECKSIG",
        "hex": "76a914b65a5d4ce219772417459b40af009f0e39646a4b88ac",
        "reqSigs": 1,
        "type": "pubkeyhash",
        "addresses": [
          "TmfaGwUbZiCeqKqrXNBaK5wEUcwcqArqNaW"
        ]
      }
    }
  ],
  "blockheight": 300,
  "confirmations": 1
}
```

generate receiving address
```shell
$ nx ec-new $(nx entropy) > receiver_key.txt
$ nx ec-to-addr $(nx ec-to-public $(cat receiver_key.txt))
TmfaGwUbZiCeqKqrXNBaK5wEUcwcqArqNaW
```

```shell
$ nx tx-encode -i cd1fb199ad3cc58d696cdd1499bc97c3fd9aac38a705db0d9bdd94d09ce1ad3e:2 -o TmfJD6F4n4iWLvgtekGCmBoR2MakzSSHm5t:2.5 -o TmfaGwUbZiCeqKqrXNBaK5wEUcwcqArqNaW:20 > tx.txt
# the key is generated in STATEMENT 1
$ nx tx-sign -k $(cat miner_key.txt) $(cat tx.txt)> tx.txt

$  cli sendRawTx $(cat tx.txt)
```

---

## About configuration
### nox
| Field | Explain |
| --- | --- |
| miningaddr | Miner account address |
| debuglevel | Logging level {trace, debug, info, warn, error, critical} |
| addpeer | Add a peer to connect with at startup |
| connect | Connect only to the specified peers at startup |

## Full nodes
| Server Name | IP Address | Describe |
| --- | --- | ---|
| Dagfans | 47.103.194.115:18130 | Shanghai |

| Pool | 47.93.20.102:18130 | Shanghai |
| Google Cloud | 35.246.127.0:18130 | UK |
| ??? | ??? | Xi'an |

<font color=Gray size=3>If you haven't turned on DNS Seed service, you can use "addpeer" to add the above servers manually as your peers.</font>


* For example:
```
docker run -it -p 18130:18130 -p 18131:18131 halalchain/nox-dag --addpeer=47.103.194.115:18130
```

* If you have multiple nodes to add, you can configure them repeatedly.
 ```
docker run -it -p 18130:18130 -p 18131:18131 halalchain/nox-dag --addpeer=x.x.x.x:x --addpeer=x.x.x.x:x
``` 
---


## Remarks
NOTE: make sure the server has at least 2GB memory

---

***Welcome to submit issues to me.***
