package message

/* //TODO(lorenzo) restore this and add more tests
func TestAddsVotingPower(t *testing.T) {
	aggregator := new(big.Int)
	core := new(big.Int)

	aggregator.SetBit(aggregator, 0, 1)
	require.True(t, AddsVotingPower(aggregator, core))

	core.SetBit(core, 0, 1)
	require.False(t, AddsVotingPower(aggregator, core))

	aggregator.SetBit(aggregator, 1000000000, 1)
	require.True(t, AddsVotingPower(aggregator, core))

	core.SetBit(core, 1000000001, 1)
	require.True(t, AddsVotingPower(aggregator, core))

	core.SetBit(core, 1000000000, 1)
	require.False(t, AddsVotingPower(aggregator, core))

	aggregator.SetBit(aggregator, 1, 1)
	core.SetBit(core, 4, 1)
	core.SetBit(core, 107, 1)
	core.SetBit(core, 64, 1)
	require.True(t, AddsVotingPower(aggregator, core))

	core.SetBit(core, 1, 1)
	require.False(t, AddsVotingPower(aggregator, core))
}*/
