# hlc-wallet command line tools
----
## Abstract
Introduce how to use the wallet command line tools for address generation and trading operations

## Get command line tools
This tools supports windows,linux and macos, please open the following page to download.
[wallet-v0.1-alpha](https://github.com/hlcfans/HLCTest/releases/tag/wallet-v0.1-alpha)

## Usage
1. Unzip the downloaded package
   - There will be two runnable programs ‘cli’ and ‘nx’
   - The ‘cli’ is an interactive command tool, 'nx' is used to encode,decode,sign etc, it's referenced internally by the cli. So, we can just use cli.
2. Create a wallet address
   - Run the cli directly in the shell or command dialog, then you will enter interactive mode.
   - If you are running the program for the first time, it will automatically generate a wallet address, typing the following command to show the address.
     ```
     [wallet-cli]: addr
     ```
   - You can also type 'help' to show usage, as following
     ```
     [wallet-cli]: help
     ```
3. Get some tokens from miner
   - If you are not miner, the wallet address does not have any tokens
   - Use the following command to get tokens from miner
     ```
     [wallet-cli]: get 10
     ```
4. Show outputs of this wallet address as following commad
   - Get unspent transaction output
     ```
     [wallet-cli]: out -u
     ```
   - Get spent transaction output
     ```
     [wallet-cli]: out -s
     ```
   - Get all transaction output
     ```
     [wallet-cli]: out -a
     ```
5. Transfer some token to the other address
   ```
   [wallet-cli]: transfer [address] 10
   ```
6. Show the status of block dag
   ```
   [wallet-cli]: status
   ```



