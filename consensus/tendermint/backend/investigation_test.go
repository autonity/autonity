package backend

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	ethereum "github.com/autonity/autonity"
	"github.com/autonity/autonity/accounts/abi"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/ethclient"
	"github.com/autonity/autonity/params"
	"github.com/autonity/autonity/params/generated"
	"github.com/stretchr/testify/require"
)

const rpcUri = "https://rpc-internal-1.piccadilly.autonity.org"

func TestChainHalt(t *testing.T) {
	client, err := ethclient.Dial(rpcUri)
	require.NoError(t, err)

	block, err := client.BlockByNumber(context.Background(), new(big.Int).SetUint64(70352))
	require.NoError(t, err)

	for _, tx := range block.Transactions() {
		if *tx.To() != params.OracleContractAddress {
			t.Log("Skip non-oracle tx")
			continue
		}
		method, err := generated.OracleAbi.MethodById(tx.Data())
		require.NoError(t, err)
		if method.Name != "vote" {
			t.Log("Skip non-vote tx")
			continue
		}
		decoded := make(map[string]interface{})
		err = method.Inputs.UnpackIntoMap(decoded, tx.Data()[4:])
		require.NoError(t, err)
		t.Log(decoded)
	}
}

func getCalledMethod(calldata []byte) string {
	method, err := generated.OracleAbi.MethodById(calldata)
	if err != nil {
		panic(err)
	}
	return method.Name
}

// converts addresses of protocol contracts to easily identifiable strings
func addressToString(addr common.Address) string {
	switch addr {
	case params.AutonityContractAddress:
		return "Autonity"
	case params.AccountabilityContractAddress:
		return "Accountability"
	case params.OracleContractAddress:
		return "Oracle"
	case params.ACUContractAddress:
		return "ACU"
	case params.SupplyControlContractAddress:
		return "SupplyControl"
	case params.StabilizationContractAddress:
		return "Stabilization"
	case params.UpgradeManagerContractAddress:
		return "UpgradeManager"
	case params.InflationControllerContractAddress:
		return "InflationController"
	case params.StakeableVestingManagerContractAddress:
		return "StakeableVestingManager"
	case params.NonStakeableVestingContractAddress:
		return "NonStakeableVesting"
	case params.OmissionAccountabilityContractAddress:
		return "OmissionAccountability"
	default:
		return "Not a protocol contract: " + addr.String()
	}
}

func getEventName(abi abi.ABI, log types.Log) string {
	event, err := abi.EventByID(log.Topics[0])
	if err != nil {
		panic(err)
	}
	return event.Name
}

type oracleData struct {
	outliers uint64
	moniker  string
	address  common.Address
}

func (od *oracleData) String() string {
	return fmt.Sprintf("address: %s, moniker: %s, outliers: %d", od.address, od.moniker, od.outliers)
}

func TestSlashingInvestigation(t *testing.T) {
	client, err := ethclient.Dial(rpcUri)
	require.NoError(t, err)

	from := 82200
	to := 82900
	t.Logf("Checking blocks from %d to %d\n", from, to)

	stakeflowOracle := common.HexToAddress("0xC9C7d1fA7adFb2856bB2F77f9F80fEa5F5a5073d")
	lesnikOracle := common.HexToAddress("0x469C4A81d86461d1E85be1741b1a3152e5BB9b98")
	decentrioOracle := common.HexToAddress("0xa32a0cad4A19AF125dB862C28382c8C83939747c")

	// maps oracle address to number of outliers
	oracles := make(map[common.Address]*oracleData)
	oracles[stakeflowOracle] = &oracleData{outliers: 0, moniker: "stakeflow", address: stakeflowOracle}
	oracles[lesnikOracle] = &oracleData{outliers: 0, moniker: "lesnik", address: lesnikOracle}
	oracles[decentrioOracle] = &oracleData{outliers: 0, moniker: "decentrio", address: decentrioOracle}

	// check why jailed
	//otherNodeAddress := common.HexToAddress("0x100E38f7BCEc53937BDd79ADE46F34362470577B")

	oracleVotes := 0
	nonOracleVotes := 0
	epochHeaders := 0
	skippedACUEvents := 0
	skippedAutonityEvents := 0
	skippedActivateValidatorEvent := 0
	skippedMintedStakeEvent := 0
	skippedRewardedEvent := 0
	skippedNewEpochEvent := 0

	for number := from; number <= to; number++ {
		block, err := client.BlockByNumber(context.Background(), new(big.Int).SetUint64(uint64(number)))
		require.NoError(t, err)

		txs := block.Transactions()
		t.Logf("block %d has %d transactions\n", number, len(txs))

		if block.IsEpochHead() {
			t.Logf("block %d is an epoch header\n", number)
			epochHeaders++
		}
		for i, tx := range txs {
			if *tx.To() == params.OracleContractAddress && getCalledMethod(tx.Data()) == "vote" {
				t.Logf("tx %d in block %d is an oracle vote\n", i, number)
				oracleVotes++
			} else {
				t.Logf("tx %d in block %d is not an oracle vote\n", i, number)
				nonOracleVotes++
			}
		}
		// retrieve logs of the block
		blockHash := block.Hash()
		logs, err := client.FilterLogs(context.Background(), ethereum.FilterQuery{BlockHash: &blockHash})
		require.NoError(t, err)
		for _, log := range logs {
			switch addressToString(log.Address) {
			case "Oracle":
				event, err := generated.OracleAbi.EventByID(log.Topics[0])
				require.NoError(t, err)
				decodedLog := make(map[string]interface{})
				err = generated.OracleAbi.UnpackIntoMap(decodedLog, event.Name, log.Data)
				require.NoError(t, err)
				if event.Name == "Penalized" {
					// print the offender and update counters
					offender := common.BytesToAddress(log.Topics[1].Bytes())
					moniker := "unknown"
					oracleMetadata, ok := oracles[offender]
					if ok {
						moniker = oracleMetadata.moniker
					}
					t.Logf("Penalized event detected, offender: %s - %s", moniker, offender.Hex())
					oracleMetadata.outliers++
				}
				t.Log(decodedLog)
			case "ACU":
				t.Log("Skipping analysis of ACU event")
				skippedACUEvents++
			case "Autonity":
				t.Log("Skipping analysis of Autonity event")
				skippedAutonityEvents++
				event, err := generated.AutonityAbi.EventByID(log.Topics[0])
				require.NoError(t, err)
				if event.Name == "ActivatedValidator" {
					t.Log("Skipping analysis of activated validator event")
					skippedActivateValidatorEvent++
					break
				}
				if event.Name == "MintedStake" {
					t.Log("Skipping analysis of minted stake event")
					skippedMintedStakeEvent++
					break
				}
				if event.Name == "Rewarded" {
					t.Log("Skipping analysis of rewarded event")
					skippedRewardedEvent++
					break
				}
				if event.Name == "NewEpoch" {
					t.Log("Skipping analysis of new epoch event")
					skippedNewEpochEvent++
					break
				}
				panic("un-handled event: " + event.Name)
			default:
				panic("un-handled contract: " + addressToString(log.Address))
			}
		}
		/* Cannot do this since Autonity does not support fetching the finalize receipt yet :(
		finalizeHash := common.ACHash(block.Number())
		receipt, err := client.TransactionReceipt(context.Background(), finalizeHash)
		require.NoError(t, err)
		fmt.Println(receipt)
		*/

	}
	t.Logf("Number of oracles votes txs: %d\n", oracleVotes)
	t.Logf("Number of non oracles votes txs: %d\n", nonOracleVotes)
	t.Logf("Number of epoch headers: %d\n", epochHeaders)
	t.Logf("Number of skipped ACU events: %d\n", skippedACUEvents)
	t.Logf("Number of skipped Autonity events: %d\n", skippedAutonityEvents)
	t.Logf("\tof which %d were ActivatedValidator event\n", skippedActivateValidatorEvent)
	t.Logf("\tof which %d were MintedStake event\n", skippedMintedStakeEvent)
	t.Logf("\tof which %d were Rewarded event\n", skippedRewardedEvent)
	t.Logf("\tof which %d were NewEpoch event\n", skippedNewEpochEvent)
	t.Logf("block range: %d - %d\n", from, to)
	t.Log("Oracle data:")
	for _, oracle := range oracles {
		t.Log(oracle.String())
	}
}
func TestOmissionJailing(t *testing.T) {
	client, err := ethclient.Dial(rpcUri)
	require.NoError(t, err)

	omissionLogs, err := client.FilterLogs(context.Background(), ethereum.FilterQuery{Addresses: []common.Address{params.OmissionAccountabilityContractAddress}})
	require.NoError(t, err)

	for _, log := range omissionLogs {
		switch getEventName(generated.OmissionAccountabilityAbi, log) {
		case "InactivityJailingEvent":
			decodedLog := make(map[string]interface{})
			err = generated.OmissionAccountabilityAbi.UnpackIntoMap(decodedLog, "InactivityJailingEvent", log.Data)
			require.NoError(t, err)
			t.Logf("block %d: validator %s jailed for inactivity", log.BlockNumber, decodedLog["validator"])
		default:
			panic("unhandled event: " + getEventName(generated.OmissionAccountabilityAbi, log))
		}
	}
}

func TestOraclePenalized(t *testing.T) {
	client, err := ethclient.Dial(rpcUri)
	require.NoError(t, err)

	var from int64 = 591870 - 24*3600
	var to int64 = 591870

	for blockNumber := from; blockNumber <= to; blockNumber += 30 {
		oracleLogs, err := client.FilterLogs(context.Background(), ethereum.FilterQuery{
			Addresses: []common.Address{params.OracleContractAddress},
			FromBlock: big.NewInt(blockNumber),
			ToBlock:   big.NewInt(blockNumber),
		})
		require.NoError(t, err)

		for _, log := range oracleLogs {
			switch getEventName(generated.OracleAbi, log) {
			case "Penalized":
				offender := common.BytesToAddress(log.Topics[1].Bytes())
				if offender == common.HexToAddress("0xc7dD8483e3dFAE05D0D86302382cc98E7F2706f0") {
					t.Logf("block %d: validator %s penalized as outlier by oracle contract", log.BlockNumber, offender)
				}
			case "NewRound":
				// do nothing
			default:
				panic("unhandled event: " + getEventName(generated.OracleAbi, log))
			}
		}
	}
}

func TestAccountabilitySlashing(t *testing.T) {
	client, err := ethclient.Dial(rpcUri)
	require.NoError(t, err)

	accountabilityLogs, err := client.FilterLogs(context.Background(), ethereum.FilterQuery{Addresses: []common.Address{params.AccountabilityContractAddress}})
	require.NoError(t, err)

	for _, log := range accountabilityLogs {
		switch getEventName(generated.AccountabilityAbi, log) {
		default:
			panic("unhandled event: " + getEventName(generated.AccountabilityAbi, log))
		}
	}
}

func TestEpochFetch(t *testing.T) {
	epochPeriod := params.PiccadillyChainConfig.AutonityContractConfig.EpochPeriod
	blockNumber := epochPeriod * (515250 / epochPeriod)

	client, err := ethclient.Dial(rpcUri)
	require.NoError(t, err)

	block, err := client.BlockByNumber(context.Background(), new(big.Int).SetUint64(blockNumber))
	require.NoError(t, err)

	require.True(t, block.IsEpochHead())

	epoch := block.Header().Epoch

	for i, v := range epoch.Committee.Members {
		fmt.Printf("validator %d: addr %v, voting power %v\n", i, v.Address, v.VotingPower)
	}

	require.Fail(t, "fails")
}
