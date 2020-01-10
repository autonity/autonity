package test

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/clearmatics/autonity/common/fdlimit"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/p2p/enode"
	"go.uber.org/goleak"

	"github.com/clearmatics/autonity/common"
	tendermintCore "github.com/clearmatics/autonity/consensus/tendermint/core"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/log"
	"golang.org/x/sync/errgroup"
)

type testCase struct {
	name                   string
	isSkipped              bool
	numPeers               int
	numBlocks              int
	txPerPeer              int
	validatorsCanBeStopped *int64

	maliciousPeers          map[int]injectors
	removedPeers            map[common.Address]uint64
	addedValidatorsBlocks   map[common.Hash]uint64
	removedValidatorsBlocks map[common.Hash]uint64 //nolint: unused, structcheck
	changedValidators       tendermintCore.Changes //nolint: unused,structcheck

	networkRates         map[int]networkRate //map[validatorIndex]networkRate
	beforeHooks          map[int]hook        //map[validatorIndex]beforeHook
	afterHooks           map[int]hook        //map[validatorIndex]afterHook
	sendTransactionHooks map[int]func(validator *testNode, fromAddr common.Address, toAddr common.Address) (bool, *types.Transaction, error)
	finalAssert          func(t *testing.T, validators []*testNode)
	stopTime             map[int]time.Time
	genesisHook          func(g *core.Genesis) *core.Genesis
	mu                   sync.RWMutex
	noQuorumAfterBlock   uint64
	noQuorumTimeout      time.Duration
}

type injectors struct {
	cons  func(basic consensus.Engine) consensus.Engine
	backs func(basic tendermintCore.Backend) tendermintCore.Backend
}

func (test *testCase) getBeforeHook(index int) hook {
	test.mu.Lock()
	defer test.mu.Unlock()

	if test.beforeHooks == nil {
		return nil
	}

	validatorHook, ok := test.beforeHooks[index]
	if !ok || validatorHook == nil {
		return nil
	}

	return validatorHook
}

func (test *testCase) getAfterHook(index int) hook {
	test.mu.Lock()
	defer test.mu.Unlock()

	if test.afterHooks == nil {
		return nil
	}

	validatorHook, ok := test.afterHooks[index]
	if !ok || validatorHook == nil {
		return nil
	}

	return validatorHook
}

func (test *testCase) setStopTime(index int, stopTime time.Time) {
	test.mu.Lock()
	test.stopTime[index] = stopTime
	test.mu.Unlock()
}

func (test *testCase) getStopTime(index int) time.Time {
	test.mu.RLock()
	currentTime := test.stopTime[index]
	test.mu.RUnlock()

	return currentTime
}

type hook func(block *types.Block, validator *testNode, tCase *testCase, currentTime time.Time) error

func runTest(t *testing.T, test *testCase) {
	if test.isSkipped {
		t.SkipNow()
	}

	defer goleak.VerifyNone(t)

	log.Root().SetHandler(log.LvlFilterHandler(log.LvlError, log.StreamHandler(os.Stderr, log.TerminalFormat(true))))
	_, err := fdlimit.Raise(512 * uint64(test.numPeers))
	if err != nil {
		t.Log("can't rise file description limit. errors are possible")
	}

	// Generate a batch of accounts to seal and fund with
	validators := make([]*testNode, test.numPeers)

	for i := range validators {
		validators[i] = new(testNode)
		validators[i].privateKey, err = crypto.GenerateKey()
		if err != nil {
			t.Fatal("cant make pk", err)
		}
	}

	for i := range validators {
		//port
		listener, innerErr := net.Listen("tcp", "127.0.0.1:0")
		if innerErr != nil {
			panic(innerErr)
		}
		validators[i].listener = append(validators[i].listener, listener)

		//rpc port
		listener, innerErr = net.Listen("tcp", "127.0.0.1:0")
		if innerErr != nil {
			panic(innerErr)
		}
		validators[i].listener = append(validators[i].listener, listener)
	}

	for i, validator := range validators {
		listener := validator.listener[0]
		validator.address = listener.Addr().String()
		port := strings.Split(listener.Addr().String(), ":")[1]
		validator.port, _ = strconv.Atoi(port)

		rpcListener := validator.listener[1]
		rpcPort, innerErr := strconv.Atoi(strings.Split(rpcListener.Addr().String(), ":")[1])
		if innerErr != nil {
			t.Fatal("incorrect rpc port ", innerErr)
		}

		validator.rpcPort = rpcPort

		if validator.port == 0 || validator.rpcPort == 0 {
			t.Fatal("On validator", i, "port equals 0")
		}

		validator.url = enode.V4URL(
			validator.privateKey.PublicKey,
			net.IPv4(127, 0, 0, 1),
			validator.port,
			validator.port,
		)
	}

	genesis := makeGenesis(validators)
	if test.genesisHook != nil {
		genesis = test.genesisHook(genesis)
	}
	for i, validator := range validators {
		var engineConstructor func(basic consensus.Engine) consensus.Engine
		var backendConstructor func(basic tendermintCore.Backend) tendermintCore.Backend
		if test.maliciousPeers != nil {
			engineConstructor = test.maliciousPeers[i].cons
			backendConstructor = test.maliciousPeers[i].backs
		}

		validator.listener[0].Close()
		validator.listener[1].Close()

		rates := test.networkRates[i]

		validator.node, err = makeValidator(genesis, validator.privateKey, validator.address, validator.rpcPort, rates.in, rates.out, engineConstructor, backendConstructor) //делаем валидатор. Он просит переменную, которая типа функция engineConstructor
		if err != nil {
			t.Fatal("cant make a validator", i, err)
		}
	}

	wg := &errgroup.Group{}
	for _, validator := range validators {
		validator := validator

		wg.Go(func() error {
			return validator.startNode()
		})
	}
	err = wg.Wait()
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		wgClose := &errgroup.Group{}
		for _, validator := range validators {
			validatorInner := validator
			wgClose.Go(func() error {
				if !validatorInner.isRunning {
					return nil
				}

				errInner := validatorInner.node.Close()
				if errInner != nil {
					return fmt.Errorf("error on node close %v", err)
				}

				validatorInner.node.Wait()

				return nil
			})
		}

		err = wgClose.Wait()
		if err != nil {
			t.Fatal(err)
		}

		time.Sleep(time.Second) //level DB needs a second to close
	}()

	wg = &errgroup.Group{}
	for _, validator := range validators {
		validator := validator

		wg.Go(func() error {
			return validator.startService()
		})
	}
	err = wg.Wait()
	if err != nil {
		t.Fatal(err)
	}

	wg = &errgroup.Group{}
	for i, validator := range validators {
		validator := validator
		i := i

		wg.Go(func() error {
			log.Debug("peers", "i", i,
				"peers", len(validator.node.Server().Peers()),
				"nodes", len(validators))
			return nil
		})
	}
	err = wg.Wait()
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		for _, validator := range validators {
			validator.subscription.Unsubscribe()
		}
	}()

	// each peer sends one tx per block
	sendTransactions(t, test, validators, test.txPerPeer, true)
	if test.finalAssert != nil {
		test.finalAssert(t, validators)
	}

	if len(test.maliciousPeers) != 0 {
		maliciousTest(t, test, validators)
	}
}
