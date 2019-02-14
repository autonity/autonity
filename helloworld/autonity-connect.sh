#!/bin/sh

# arrays were not cooperating in the docker container (busybox) shell
ENODE_1=$(autonity attach http://172.25.0.11:8545 --exec "admin.nodeInfo.enode.slice(8,136)")
ENODE_2=$(autonity attach http://172.25.0.12:8545 --exec "admin.nodeInfo.enode.slice(8,136)")
ENODE_3=$(autonity attach http://172.25.0.13:8545 --exec "admin.nodeInfo.enode.slice(8,136)")
ENODE_4=$(autonity attach http://172.25.0.14:8545 --exec "admin.nodeInfo.enode.slice(8,136)")
ENODE_5=$(autonity attach http://172.25.0.15:8545 --exec "admin.nodeInfo.enode.slice(8,136)")

ENODE_1="enode://$(echo "${ENODE_1//\"}")@172.25.0.11:30303"
ENODE_2="enode://$(echo "${ENODE_2//\"}")@172.25.0.12:30303"
ENODE_3="enode://$(echo "${ENODE_3//\"}")@172.25.0.13:30303"
ENODE_4="enode://$(echo "${ENODE_4//\"}")@172.25.0.14:30303"
ENODE_5="enode://$(echo "${ENODE_5//\"}")@172.25.0.15:30303"

echo "eNode 1: $ENODE_1"
echo "eNode 2: $ENODE_2"
echo "eNode 3: $ENODE_3"
echo "eNode 4: $ENODE_4"
echo "eNode 5: $ENODE_5"

# connect all the things! (connect nodes between them)
for i in 1 2 3 4 5
do
  ADDRESS="http://172.25.0.1$i:8545"
  $(autonity attach $ADDRESS --exec "admin.addPeer('$ENODE_1')")
  $(autonity attach $ADDRESS --exec "admin.addPeer('$ENODE_2')")
  $(autonity attach $ADDRESS --exec "admin.addPeer('$ENODE_3')")
  $(autonity attach $ADDRESS --exec "admin.addPeer('$ENODE_4')")
  $(autonity attach $ADDRESS --exec "admin.addPeer('$ENODE_5')")
done

# unlock all the things! (addresses)
# set coinbase address
for i in 1 2 3 4 5
do
  IDX=$(($i - 1))
  ADDRESS="http://172.25.0.1$i:8545"
  UNLOCKED=$(autonity attach $ADDRESS --exec "personal.unlockAccount(eth.accounts[$IDX],'test')")
  IS_COINBASE_SET=$(autonity attach $ADDRESS --exec "miner.setEtherbase(eth.accounts[$IDX])")
  COINBASE=$(autonity attach $ADDRESS --exec "eth.coinbase")
  echo "Node $i $ADDRESS Account: $COINBASE Coinbase: $IS_COINBASE_SET Unlocked: $UNLOCKED"
  # mine all the blocks! (start mining node)
  IS_MINING=$(autonity attach $ADDRESS --exec "miner.start()")
  echo "Node $i is mining"
done


# for i in {1..5}; do ./autonity attach http://0.0.0.0:854$i --exec '[eth.coinbase, eth.getBlock("latest").number, eth.getBlock("latest").hash, eth.mining]'; done
