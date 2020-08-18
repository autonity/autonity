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

# init the data directory
echo "Autonity INIT $RPC_ADDR"
$AUTONITY init --datadir $DATADIR genesis-tendermint.json

# start the node with the keystore and nodekey
echo "Autonity START"
$AUTONITY \
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
  --unlock 0x850c1eb8d190e05845ad7f84ac95a318c8aab07f \
  --password password \
  --miner.threads 1 \
  --verbosity 1
