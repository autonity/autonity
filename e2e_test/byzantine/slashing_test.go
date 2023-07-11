package byzantine

// *** e2e tests to write ***:
// Part 1: Slashing tests
// - Staking with penalty absorbing stake
// - Funds moved to treasury
// - Reward redistribution
// - Jail
// - Silence consensus

// Part 2: Accusation flow tests:

//  Validator is accused, submit proof of innocence,
//  Validator is accused, do not submit proof of innocence,
//  Validator is accused and accused again
//  Validators is accused and someone sent direct proof of misbehavior

//  Need to test canAccuse/canSlash for each scenario

// Part 3: Fuzz tests on event handler
