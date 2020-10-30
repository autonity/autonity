#!/bin/sh

AUTONITY=autonity

DATADIR=autonity-data
GENESIS=genesis.json
KEYSTORE=keystore
NODEKEY=nodekey
WS_PORT=8645
WS_ADDR=$(awk 'END{print $1}' /etc/hosts)
WS_API="tendermint,console,eth,web3,admin,debug,miner,personal,txpool,net"
RPC_PORT=8545
RPC_ADDR=$(awk 'END{print $1}' /etc/hosts)
RPC_API="tendermint,console,eth,web3,admin,debug,miner,personal,txpool,net"

# start the node with the keystore and nodekey
echo "Autonity START"
$AUTONITY \
  --datadir $DATADIR \
  --genesis $GENESIS \
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
  --verbosity 4 \
  --debug
