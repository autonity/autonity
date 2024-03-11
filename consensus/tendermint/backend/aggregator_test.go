package backend

//TODO(lorenzo) write tests for:
// - aggregator receives a valid complex aggregate (carrying quorum) for a future round (IMPORTANT!!!)
// - the ignoring behaviour when a peer is disconnected (votesFrom and toIgnore)
// - aggregating a batch of signatures containing an invalid one
//		- explicitly check that the loop to exclude invalid indexes from the validVotes set works as intended
//		- ref: https://github.com/autonity/autonity/blob/a7fe1a0dc98a7668be4b7d2bbd23ea3c84c3868a/consensus/tendermint/backend/aggregator.go#L146
// - aggregation of maliciously crafted batches to try to have signatures with high coefficient for a single validator (e.g. (A^255,B^1,C^1))
