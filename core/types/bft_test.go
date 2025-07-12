package types

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/autonity/autonity/params"
	"github.com/stretchr/testify/require"

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
	posHeaderHash := common.HexToHash("0x5cf94f58b040fca7c695f41a18f447f85955cda27247c98ed24d65bc798cc2f5")

	quorumCertificate := &AggregateSignature{}
	testKey, err := blst.SecretKeyFromHex("667e85b8b64622c4b8deadf59964e4c6ae38768a54dbbbc8bbd926777b896584")
	require.NoError(t, err)
	quorumCertificate.Signature = testKey.Sign([]byte("0xcafe")).(*blst.BlsSignature)
	quorumCertificate.Signers = NewSigners(1)
	quorumCertificate.Signers.increment(0, common.Big1)

	activityProof := quorumCertificate.Copy()

	// add committee to header's EpochExtra.
	c := &Committee{
		Members: []CommitteeMember{
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
	}

	epoch := &Epoch{
		PreviousEpochBlock: common.Big0,
		NextEpochBlock:     common.Big256,
		Committee:          c,
		OmissionDelta:      common.Big5,
		Eip1559: &Eip1559Params{
			MinBaseFee:               new(big.Int).SetUint64(params.TestMinBaseFee),
			BaseFeeChangeDenominator: new(big.Int).SetUint64(8),
			ElasticityMultiplier:     common.Big2,
			GasLimitBoundDivisor:     new(big.Int).SetUint64(1024),
		},
	}
	signature := testKey.Sign(testKey.PublicKey().Marshal())
	proposerSeal := signature.Marshal()

	epoch2 := epoch.Copy()
	epoch2.OmissionDelta = common.Big2

	testCases := []struct {
		header Header
		hash   common.Hash
	}{
		// Non-BFT header tests, PoS fields should not be taken into account.
		{
			Header{},
			common.HexToHash("0x31775027428c8e2ea15ba3512c7dc42bc356fe2cba2a82b537b3633f7e1907ab"),
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
				QuorumCertificate: quorumCertificate,
			}),
			posHeaderHash,
		},
		{
			setExtra(PosHeader, headerExtra{
				Epoch: epoch,
			}),
			common.HexToHash("0x2b2a4d57f7c0b8cce80df51f55972681cbc17307da92a25503658a7c99a5b86f"),
		},
		{
			setExtra(PosHeader, headerExtra{
				Epoch: epoch2,
			}),
			common.HexToHash("0x20a218b23a9f927d6c53d217e1a6d3ff2e26c00fa95cc1adc03946e45a071d13"),
		},
		{
			setExtra(PosHeader, headerExtra{
				ProposerSeal: proposerSeal,
			}),
			common.HexToHash("0xd327f4e2e84d68da696a1aab40ae628c63256d186b2201249e06df5898677d2b"),
		},
		{
			setExtra(PosHeader, headerExtra{
				ActivityProof: activityProof,
			}),
			common.HexToHash("0xeed7caad41663db464b20656dc1926fac12cacec0f1bd9634d265ecb393e8766"),
		},
		{
			setExtra(PosHeader, headerExtra{
				ActivityProofRound: uint64(7),
			}),
			common.HexToHash("0x8f67bc27393f87eca2aab3d6c67de572a3dfbc0f0f7e6d96a977fb7285230c2a"),
		},
		{
			setExtra(PosHeader, headerExtra{
				ActivityProof:      activityProof,
				ActivityProofRound: uint64(7),
			}),
			common.HexToHash("0xe87d18029c1cc7e01a500ffde4c8357dad73a7e23f702d045d403bd8d06a399c"),
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
	h.ProposerSeal = hExtra.ProposerSeal
	h.Round = hExtra.Round
	h.QuorumCertificate = hExtra.QuorumCertificate
	h.Epoch = hExtra.Epoch
	h.ActivityProof = hExtra.ActivityProof
	h.ActivityProofRound = hExtra.ActivityProofRound
	return h
}
