#!/bin/sh

AUTONITY=autonity

DATADIR=autonity-data
KEYSTORE=keystore
NODEKEY=nodekey
RPC_PORT=8545
RPC_ADDR=$(awk 'END{print $1}' /etc/hosts)
RPC_API="tendermint,clique,console,eth,web3,admin,debug,miner,personal,txpool,net"

# init the data directory
if [ ! -d "$DATADIR/autonity/chaindata" ]; then
  echo "Autonity INIT $RPC_ADDR $DATADIR\autonity\chaindata"
  $AUTONITY init --datadir $DATADIR genesis.json --verbosity 4
fi

# start the node with the keystore and nodekey
echo "Autonity START"
$AUTONITY \
  --datadir $DATADIR \
  --ethash.cachedir "$DATADIR/cache" \
  --ethash.dagdir "$DATADIR/.ethash" \
  --nodekey $NODEKEY \
  --keystore $KEYSTORE \
  --rpc \
  --rpcaddr $RPC_ADDR \
  --rpcport $RPC_PORT \
  --rpcapi "$RPC_API" \
  --rpccorsdomain "*" \
  --syncmode "full" \
  --mine \
  --miner.threads 1 \
  --verbosity 4 \
  --debug 2>&1 | tee logs.log
