#!/bin/sh

AUTONITY=autonity

DATADIR=autonity-data
KEYSTORE=keystore
NODEKEY=nodekey
RPC_PORT=8545
RPC_ADDR=$(awk 'END{print $1}' /etc/hosts)
RPC_API="istanbul,clique,console,eth,web3,admin,debug,miner,personal,txpool,net"

# init the data directory
echo "Autonity INIT $RPC_ADDR"
$AUTONITY init --datadir $DATADIR genesis.json

# start the node with the keystore and nodekey
echo "Autonity START"
cp ./static-nodes.json $DATADIR
$AUTONITY \
  --datadir $DATADIR \
  --nodekey $NODEKEY \
  --keystore $KEYSTORE \
  --rpc \
  --rpcaddr $RPC_ADDR \
  --rpcport $RPC_PORT \
  --rpcapi "$RPC_API" \
  --rpccorsdomain "*" \
  --syncmode "full" \
  --minerthreads 1 \
  --verbosity 4 --debug
