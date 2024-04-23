package types

import (
	"github.com/autonity/autonity/crypto"
	"github.com/stretchr/testify/require"
	"log"
	"math/big"
	"reflect"
	"testing"

	"github.com/autonity/autonity/common"
)

func TestHeaderHash(t *testing.T) {
	originalHeader := Header{
		ParentHash:  common.HexToHash("0000H45H"),
		UncleHash:   common.HexToHash("0000H45H"),
		Coinbase:    common.HexToAddress("0000H45H"),
		Root:        common.HexToHash("0000H00H"),
		TxHash:      common.HexToHash("0000H45H"),
		ReceiptHash: common.HexToHash("0000H45H"),
		Difficulty:  big.NewInt(1337),
		Number:      big.NewInt(1337),
		GasLimit:    1338,
		GasUsed:     1338,
		Time:        1338,
		BaseFee:     big.NewInt(0),
		Extra:       []byte("Extra data Extra data Extra data  Extra data  Extra data  Extra data  Extra data Extra data"),
		MixDigest:   common.HexToHash("0x0000H45H"),
	}
	PosHeader := originalHeader
	PosHeader.MixDigest = BFTDigest

	originalHeaderHash := common.HexToHash("0xda0ef4df9161184d34a5af7e80b181626f197781e1c51557522047b0eaa63605")
	posHeaderHash := common.HexToHash("0x013d47cf6a2874b7396ec9cb79acbff51c932f487380db8c66c363e22a2f9001")
	committee := new(Committee)
	committee.Members = append(committee.Members, &CommitteeMember{
		Address:     common.HexToAddress("0x1234566"),
		VotingPower: new(big.Int).SetUint64(12),
	}, &CommitteeMember{
		Address:     common.HexToAddress("0x13371337"),
		VotingPower: new(big.Int).SetUint64(1337),
	})

	testCases := []struct {
		header Header
		hash   common.Hash
	}{
		// Non-BFT header tests, PoS fields should not be taken into account.
		{
			Header{},
			common.HexToHash("0xc3bd2d00745c03048a5616146a96f5ff78e54efb9e5b04af208cdaff6f3830ee"),
		},
		{
			originalHeader,
			originalHeaderHash,
		},
		{
			setExtra(originalHeader, headerExtra{}),
			originalHeaderHash,
		},

		// BFT header tests
		{
			PosHeader, // test 3
			posHeaderHash,
		},
		{
			setExtra(PosHeader, headerExtra{
				CommittedSeals: [][]byte{common.Hex2Bytes("0xfacebooc"), common.Hex2Bytes("0xbabababa")},
			}),
			posHeaderHash,
		},
		{
			setExtra(PosHeader, headerExtra{
				CommittedSeals: [][]byte{common.Hex2Bytes("0x123456"), common.Hex2Bytes("0x777777"), common.Hex2Bytes("0xaaaaaaa")},
			}),
			posHeaderHash,
		},
		{
			setExtra(PosHeader, headerExtra{
				Committee: committee,
			}),

			common.HexToHash("0x1660f82e07fda09ac04125e46d6bddca0249bcc790ccbb0cb0d1672a6936054f"),
		},
		{
			setExtra(PosHeader, headerExtra{
				ProposerSeal: common.Hex2Bytes("0xbebedead"),
			}),
			posHeaderHash,
		},
		{
			setExtra(PosHeader, headerExtra{
				Round: 1997,
			}),
			posHeaderHash,
		},
		{
			setExtra(PosHeader, headerExtra{
				Round: 3,
			}),
			posHeaderHash,
		},
		{
			setExtra(PosHeader, headerExtra{
				Round: 0,
			}),
			posHeaderHash,
		},
	}
	for i := range testCases {
		if !reflect.DeepEqual(testCases[i].hash, testCases[i].header.Hash()) {
			t.Errorf("test %d, expected: %v, but got: %v", i, testCases[i].hash.Hex(), testCases[i].header.Hash().Hex())
		}
	}
}

func setExtra(h Header, hExtra headerExtra) Header {
	h.Committee = hExtra.Committee
	h.ProposerSeal = hExtra.ProposerSeal
	h.Round = hExtra.Round
	h.CommittedSeals = hExtra.CommittedSeals
	return h
}

func TestElectProposer(t *testing.T) {
	samePowers := []int{100, 100, 100, 100}
	linearPowers := []int{100, 200, 400, 800}
	t.Run("Proposer election should be deterministic", func(t *testing.T) {
		committee := generateCommittee(samePowers)
		for h := uint64(0); h < uint64(100); h++ {
			for r := int64(0); r <= int64(3); r++ {
				proposer1 := committee.Proposer(h, r)
				proposer2 := committee.Proposer(h, r)
				require.Equal(t, proposer1, proposer2)
			}
		}
	})

	t.Run("Proposer selection, print and compare the scheduling rate with same stake", func(t *testing.T) {
		committee := generateCommittee(samePowers)
		maxHeight := uint64(10000)
		maxRound := int64(4)
		//expectedRatioDelta := float64(0.01)
		counterMap := make(map[common.Address]int)
		counterMap[common.Address{}] = 1
		for h := uint64(0); h < maxHeight; h++ {
			for round := int64(0); round < maxRound; round++ {
				proposer := committee.Proposer(h, round)
				_, ok := counterMap[proposer]
				if ok {
					counterMap[proposer]++
				} else {
					counterMap[proposer] = 1
				}
			}
		}

		totalStake := 0
		for _, s := range samePowers {
			totalStake += s
		}

		for i, c := range committee.Members {
			stake := samePowers[i]
			scheduled := counterMap[c.Address]
			log.Print("electing ", "proposer: ", c.Address.String(), " stake: ", stake, " scheduled: ", scheduled)
		}
	})

	t.Run("Proposer selection, print and compare the scheduling rate with liner increasing stake", func(t *testing.T) {
		committee := generateCommittee(linearPowers)
		maxHeight := uint64(1000000)
		maxRound := int64(4)
		//expectedRatioDelta := float64(0.01)
		counterMap := make(map[common.Address]int)
		counterMap[common.Address{}] = 1
		for h := uint64(0); h < maxHeight; h++ {
			for round := int64(0); round < maxRound; round++ {
				proposer := committee.Proposer(h, round)
				_, ok := counterMap[proposer]
				if ok {
					counterMap[proposer]++
				} else {
					counterMap[proposer] = 1
				}
			}
		}

		totalStake := 0
		for _, s := range samePowers {
			totalStake += s
		}

		for _, c := range committee.Members {
			stake := c.VotingPower.Uint64()
			scheduled := counterMap[c.Address]
			log.Print("electing ", "proposer: ", c.Address.String(), " stake: ", stake, " scheduled: ", scheduled)
		}
	})
}
func generateCommittee(powers []int) *Committee {
	c := new(Committee)
	for _, p := range powers {
		privateKey, _ := crypto.GenerateKey()
		c.Members = append(c.Members, &CommitteeMember{
			Address:     crypto.PubkeyToAddress(privateKey.PublicKey),
			VotingPower: new(big.Int).SetInt64(int64(p)),
		})
	}

	c.Sort()
	return c
}
