#!/bin/sh


# if there is already an autonity network listening on 8545,
# do not start a new one

lsof -i :8545 | grep autonity
if [ $? -ne 0 ]; then
  AUTONITY=../../../../build/bin/autonity

  DATADIR=data
  KEYSTORE=keystore
  AUTONITYKEYS=autonitykeys1
  WS_PORT=8645
  WS_ADDR=127.0.0.1
  WS_API="admin,debug,eth,net,web3,txpool,miner"
  RPC_PORT=8545
  RPC_ADDR=127.0.0.1
  RPC_API="admin,debug,eth,net,web3,txpool,miner"

  # start the node with the keystore and nodekey
  echo "Autonity START"
  $AUTONITY \
    --genesis genesis-tendermint.json \
    --datadir $DATADIR \
    --autonitykeys $AUTONITYKEYS \
    --keystore $KEYSTORE \
    --ws \
    --ws.addr $WS_ADDR \
    --ws.port $WS_PORT \
    --ws.api "$WS_API" \
    --http \
    --http.addr $RPC_ADDR \
    --http.port $RPC_PORT \
    --http.api "$RPC_API" \
    --http.corsdomain "*" \
    --syncmode "full" \
    --miner.gaslimit 100000000 \
    --password password \
    --mine
fi
