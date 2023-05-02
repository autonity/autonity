#!/bin/sh

AUTONITY=../../../../build/bin/autonity

DATADIR=data
KEYSTORE=keystore
NODEKEY=nodekey1
WS_PORT=8645
WS_ADDR=127.0.0.1
WS_API="tendermint,eth,web3,admin,debug,miner,personal,txpool,net"
RPC_PORT=8545
RPC_ADDR=127.0.0.1
RPC_API="tendermint,eth,web3,admin,debug,miner,personal,txpool,net"

# start the node with the keystore and nodekey
echo "Autonity START"
$AUTONITY \
  --genesis genesis-tendermint.json \
  --datadir $DATADIR \
  --nodekey $NODEKEY \
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
  --allow-insecure-unlock \
  --unlock 0x850c1eb8d190e05845ad7f84ac95a318c8aab07f,0x4ad219b58a5b46a1d9662beaa6a70db9f570dea5,0x4b07239bd581d21aefcdee0c6db38070f9a5fd2d,0xc443c6c6ae98f5110702921138d840e77da67702,0x09428e8674496e2d1e965402f33a9520c5fcbbe2,0x64852003fc0b84d6c49c5cb3dfcd17922affddc1,0x4839950a5f07d6d6cd82f933d1de8574c48d6e74,0x160bc705bf2e5871557722c9332cfa185c02b765,0xe12b43B69E57eD6ACdd8721Eb092BF7c8D41Df41,0xDE03B7806f885Ae79d2aa56568b77caDB0de073E \
  --password password \
  --miner.threads 1 
