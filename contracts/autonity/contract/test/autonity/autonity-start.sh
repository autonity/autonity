#!/bin/sh

AUTONITY=../../../../../build/bin/autonity

DATADIR=data
KEYSTORE=keystore
NODEKEY=nodekey1
WS_PORT=8645
WS_ADDR=127.0.0.1
WS_API="tendermint,console,eth,web3,admin,debug,miner,personal,txpool,net"
RPC_PORT=8545
RPC_ADDR=127.0.0.1
RPC_API="tendermint,console,eth,web3,admin,debug,miner,personal,txpool,net"

# start the node with the keystore and nodekey
echo "Autonity START"
$AUTONITY \
  --genesis genesis-tendermint.json \
  --datadir $DATADIR \
  --nodekey $NODEKEY \
  --keystore $KEYSTORE \
  --ws \
  --wsaddr $WS_ADDR \
  --wsport $WS_PORT \
  --wsapi "$WS_API" \
  --rpc \
  --rpcaddr $RPC_ADDR \
  --rpcport $RPC_PORT \
  --rpcapi "$RPC_API" \
  --rpccorsdomain "*" \
  --syncmode "full" \
  --mine \
  --allow-insecure-unlock \
  --miner.threads 1 \
  --verbosity 1
