#!/bin/sh

AUTONITY=autonity

DATADIR=autonity-data
KEYSTORE=keystore
NODEKEY=nodekey
RPC_PORT=8545
RPC_ADDR=$(awk 'END{print $1}' /etc/hosts)
RPC_API="clique,console,eth,web3,admin,debug,miner,personal,txpool,net"

echo "Autonity INIT $RPC_ADDR"
$AUTONITY init --datadir $DATADIR genesis.json


echo "Autonity START"
$AUTONITY \
  --datadir $DATADIR \
  --nodekey $NODEKEY \
  --keystore $KEYSTORE \
  --rpc \
  --rpcaddr $RPC_ADDR \
  --rpcport $RPC_PORT \
  --rpcapi "$RPC_API" \
  --rpccorsdomain "*"
