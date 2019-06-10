# DAG create sign raw transaction hash script

## The Step

+ need a account with address and private key
```bash
$ ./script --network test --action generate-new-address
```   
+ start a `node` with modify some config up to machine config

+ [mining](https://github.com/HLCTest/hlc-miner) 120 `blocks` for getting `coinbase` reward

+ run the command to cost target block for creating many signed transactions

+ use the shell script to batch generate more signed transactions


## Requirements

[Go](http://golang.org) 1.11 or newer.

## Compile
    
```bash
$ git clone git@github.com:HLCTest/Nox-DAG-test.git
$ cd Nox-DAG-test/script
$ go build
```    
    
    
## Run
- create a new address with private key
```bash
$ ./script --network test  --action generate-new-address
```
- create 999 signed raw transactions
```bash
$ ./script --faddress=(your address which has much money) \ 
--privkey=xxxxxx -s (your node rpc server) \
-u (your node rpc user) -P (your node rpc pass) \
--addressfile=/tmp/address.csv --txfile=/tmp/tx.csv \
--notls --network=test --height=1
```
- the param `--send=true` will auto send the transactions to `node`
    
- create more transactions
    - first param is start `coinbase` block 
    - second param is end `coinbase` block 
    - third param is network
```bash
$ ./batch 1 30 test
```