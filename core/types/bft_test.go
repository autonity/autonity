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
	posHeaderHash := common.HexToHash("0x5cf94f58b040fca7c695f41a18f447f85955cda27247c98ed24d65bc798cc2f5")

	quorumCertificate := &AggregateSignature{}
	testKey, _ := blst.SecretKeyFromHex("667e85b8b64622c4b8deadf59964e4c6ae38768a54dbbbc8bbd926777b896584")
	quorumCertificate.Signature = testKey.Sign([]byte("0xcafe")).(*blst.BlsSignature)
	quorumCertificate.Signers = NewSigners(1)
	quorumCertificate.Signers.increment(0)

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

	epoch := &Epoch{PreviousEpochBlock: common.Big0, NextEpochBlock: common.Big256, Committee: c, Delta: common.Big5}
	signature := testKey.Sign(testKey.PublicKey().Marshal())
	proposerSeal := signature.Marshal()

	epoch2 := epoch.Copy()
	epoch2.Delta = common.Big2

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
				QuorumCertificate: quorumCertificate,
			}),
			posHeaderHash,
		},
		{
			setExtra(PosHeader, headerExtra{
				Epoch: epoch,
			}),
			common.HexToHash("0x4003fa038541d2ee678ed013f972109f55c9eea389dc1afedea08250e260dd41"),
		},
		{
			setExtra(PosHeader, headerExtra{
				Epoch: epoch2,
			}),
			common.HexToHash("0x4e7ad6d14406b030f430e98d9b4d9950c0de2af715c3f6ffa2d5425a0262e23e"),
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
			common.HexToHash("0xed56b294e28b72c062a85e0d8a6df84f5215fb433d2d66960f27a2c675fa62f4"),
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
			common.HexToHash("0xd5ceac2d0f738b13bd838f5f8ecf4741c25143bac9b01674118c9596b67ed6f7"),
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
