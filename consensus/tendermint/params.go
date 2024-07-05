package tendermint

// Wait until the GST + delta blocks to start accounting.
// IMPORTANT: this value must match the value of `Delta` in autonity/solidity/contracts/Autonity.sol
const DeltaBlocks = 10
