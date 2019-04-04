#!/bin/sh

# unlock all the things! (addresses)
# set coinbase address
for i in 1 2 3 4
do
  IDX=$(($i - 1))
  ADDRESS="http://172.25.0.1$i:8545"
  UNLOCKED=$(autonity attach $ADDRESS --exec "personal.unlockAccount(eth.accounts[$IDX],'test')")
  # IS_COINBASE_SET=$(autonity attach $ADDRESS --exec "miner.setEtherbase(eth.accounts[$IDX])")
  COINBASE=$(autonity attach $ADDRESS --exec "eth.coinbase")
  echo "Node $i $ADDRESS Account: $COINBASE Unlocked: $UNLOCKED"
  # mine all the blocks! (start mining node)
  IS_MINING=$(autonity attach $ADDRESS --exec "miner.start()")
  echo "Node $i is mining"
done

# for i in {1..5}; do ./autonity attach http://0.0.0.0:854$i --exec '[eth.coinbase, eth.getBlock("latest").number, eth.getBlock("latest").hash, eth.mining]'; done
