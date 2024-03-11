package types

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/crypto/blst"
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
	posHeaderHash := common.HexToHash("0xc43c096ec826fcaad1c48f5c3418a424fc8ff1def6debc1fa5b8cd43bdacdf86")

	seals := make(Signatures)
	testKey, _ := blst.SecretKeyFromHex("667e85b8b64622c4b8deadf59964e4c6ae38768a54dbbbc8bbd926777b896584")
	seals[common.HexToAddress("0xcafe")] = testKey.Sign([]byte("0xcafe")).(*blst.BlsSignature)

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
				CommittedSeals: seals,
			}),
			posHeaderHash,
		},
		{
			setExtra(PosHeader, headerExtra{
				Committee: Committee{
					{
						Address:           common.HexToAddress("0x1234566"),
						VotingPower:       new(big.Int).SetUint64(12),
						ConsensusKeyBytes: testKey.PublicKey().Marshal(),
						ConsensusKey:      testKey.PublicKey(),
					},
					{
						Address:           common.HexToAddress("0x13371337"),
						VotingPower:       new(big.Int).SetUint64(1337),
						ConsensusKeyBytes: testKey.PublicKey().Marshal(),
						ConsensusKey:      testKey.PublicKey(),
					},
				},
			}),
			common.HexToHash("0x6917cd41cf89cb7a1ad33f11f543d46cc3420d365d5b2d38167af605d8992525"),
		},
		{
			setExtra(PosHeader, headerExtra{
				ProposerSeal: common.Hex2Bytes("0xbebedead"),
			}),
			common.HexToHash("0xc43c096ec826fcaad1c48f5c3418a424fc8ff1def6debc1fa5b8cd43bdacdf86"),
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
