# Performance regression testing
Once per releace we'd like to measure performance metrics on different levels.
## Metrics to control
### Network level

| **Name**                               | **Type**  | **Metrics** | **Description** |
| -------------------------------------- | --------- | ------------ | ----------------- |
| consensus_block_interval               | Histogram | percentiles 50, 75 and 90, max | Time between this and last block (Block.Header.Time) in seconds |
| consensus_rounds                       | Gauge     | percentiles 50, 75 and 90, max | Number of rounds |
| consensus_block_propose_step_interval       | Histogram | percentiles 50, 75 and 90, max | Propose step duration in microseconds |
| consensus_block_prevote_step_interval       | Histogram | percentiles 50, 75 and 90, max | Prevote step duration in microseconds |
| consensus_block_precommit_step_interval       | Histogram | percentiles 50, 75 and 90, max | Precommit step duration in microseconds |
| consensus_block_commit_step_interval       | Histogram | percentiles 50, 75 and 90, max | Commit step duration in microseconds |
| consensus_validator_missed_blocks      | counter   | count | Total amount of blocks missed for the node, if the node is a validator |
| consensus_txs_per_block                | Gauge   | percentiles 50, 75 and 90, max | number of transactions per block |
| consensus_transaction_time             | Histogram | percentiles 50, 75 and 90, max | Time between a transaction was sent to a peer and included into a block in seconds |
| consensus_rounds                       | Gauge     | percentiles 50, 75 and 90, max | Number of rounds |
| consensus_validator_missed_blocks      | counter   | count | Total amount of blocks missed for the node, if the node is a validator |

### Node level
| **Name**                               | **Type**  | **Metrics** | **Description** |
| -------------------------------------- | --------- | ------------ | ----------------- |
| p2p_peer_receive_bytes_total           | counter   | count | number of bytes per channel received to a peer |
| p2p_peer_send_bytes_total              | counter   | count | number of bytes per channel per block sent to a peer |
| p2p_peer_receive_bytes_per_block       | Gauge   | percentiles 50, 75 and 90, max | number of bytes per channel received from a peer |
| p2p_peer_send_bytes_per_block          | Gauge   | percentiles 50, 75 and 90, max | number of bytes per channel ber block sent to a peer |
| mempool_size                           | Gauge     | percentiles 50, 75 and 90, max | Number of uncommitted transactions |
| mempool_failed_txs                     | counter   | count | number of failed transactions |

## Report format
In ideal it should be the same as used in goperf to be able to use golang tooling like benchcmp https://godoc.org/golang.org/x/tools/cmd/benchcmp , that compares different statistics in a correct way (t-test, f-test). Or we should write our own with correct comparacing.

## Test cases
* Python tests
* Melicious tests:
    * 1/3 validators offline
    * melicious validators try to send multiple votes for different or the same blocks: 1 validator, 5%, 10%, 1/3 of validators count
    * transaction spam attack: 10, 100, 1000 times traffic comparing to normal
    * slow nodes: 10, 100, 1000 less network bandwidth than normal nodes

## Resources
Basic go-ethereum and tendermint metrics was used.