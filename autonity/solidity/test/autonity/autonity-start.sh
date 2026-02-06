#!/bin/sh

set -eu

AUTONITY=../../../../build/bin/autonity

DATA1=data1
KEYSTORE=keystore
AUTONITYKEYS1=autonitykeys1

RPC_ADDR=127.0.0.1
RPC_PORT=8545
RPC_API="eth,web3,admin,debug,miner,personal,txpool,net"

WS_ADDR=127.0.0.1
WS_PORT=8645
WS_API="eth,web3,admin,debug,miner,personal,txpool,net"

P2P_PORT1=30303
CONS_PORT1=21203

MAXPEERS=25

echo "Autonity START"

# If there is already an autonity JSON-RPC listening on 8545, do not start a new one.
if curl -sf -H 'Content-Type: application/json' \
  --data '{"jsonrpc":"2.0","method":"web3_clientVersion","params":[],"id":1}' \
  "http://$RPC_ADDR:$RPC_PORT" >/dev/null 2>&1; then
  exit 0
fi

# Clean any stale chaindata (genesis changes between runs).
rm -rf "$DATA1"
mkdir -p "$DATA1"

# Node 1: JSON-RPC for Truffle. `--mine` is required for TestMode (FullFaker)
# networks to produce blocks without peers.
$AUTONITY \
  --genesis genesis-tendermint.json \
  --datadir "$DATA1" \
  --autonitykeys "$AUTONITYKEYS1" \
  --keystore "$KEYSTORE" \
  --verbosity 3 \
  --mine \
  --miner.recommit 200ms \
  --http \
  --http.addr "$RPC_ADDR" \
  --http.port "$RPC_PORT" \
  --http.api "$RPC_API" \
  --http.corsdomain "*" \
  --ws \
  --ws.addr "$WS_ADDR" \
  --ws.port "$WS_PORT" \
  --ws.api "$WS_API" \
  --syncmode "full" \
  --port "$P2P_PORT1" \
  --maxpeers "$MAXPEERS" \
  --nodiscover \
  --nat none \
  --bootnodes "" \
  --discovery.dns "" \
  --consensus.port "$CONS_PORT1" \
  --consensus.nat none \
  --rpc.txfeecap 0 \
  >"$DATA1/autonity.log" 2>&1 &

# Wait for RPC to come up.
i=0
while [ $i -lt 400 ]; do
  if curl -sf -H 'Content-Type: application/json' \
    --data '{"jsonrpc":"2.0","method":"web3_clientVersion","params":[],"id":1}' \
    "http://$RPC_ADDR:$RPC_PORT" >/dev/null 2>&1; then
    break
  fi
  i=$((i+1))
  sleep 0.05
done

# Wait until tx indexing is finished (avoids "transaction indexing is in progress" during receipt polling).
i=0
while [ $i -lt 400 ]; do
  SYNCING="$(curl -sf -H 'Content-Type: application/json' \
    --data '{"jsonrpc":"2.0","method":"eth_syncing","params":[],"id":1}' \
    "http://$RPC_ADDR:$RPC_PORT" | tr -d '\n' | tr -d '\r')"
  echo "$SYNCING" | grep -q '"result":false' && break
  i=$((i+1))
  sleep 0.05
done
