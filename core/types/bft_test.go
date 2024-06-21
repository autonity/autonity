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
	posHeaderHash := common.HexToHash("0xc1baba79e0c591d03dff86ed1a2e1a612bba7dde7c4eaa588a6a106fb0bb9b60")

	quorumCertificate := AggregateSignature{}
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

	epoch := &Epoch{PreviousEpochBlock: common.Big0, NextEpochBlock: common.Big256, Committee: c}
	signature := testKey.Sign(testKey.PublicKey().Marshal())
	proposerSeal := signature.Marshal()

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
			common.HexToHash("0x4429e362109095b8b26dc2b562c20baf0ac9f0ae7691131a3ca3dd1e79799239"),
		},
		{
			setExtra(PosHeader, headerExtra{
				ProposerSeal: proposerSeal,
			}),
			common.HexToHash("0x8e89c040c96dee171875b4127032b51f32b8d92be8144befbad7798bafb42e0b"),
		},
		{
			setExtra(PosHeader, headerExtra{
				ActivityProof: activityProof,
			}),
			common.HexToHash("0xdf76b7d1c578314a7dc4e5a6dd93e17a5d754bc7ebbb5deece3e2c843de81cbe"),
		},
		{
			setExtra(PosHeader, headerExtra{
				ActivityProofRound: uint64(7),
			}),
			common.HexToHash("0x2eaeb68590ba4b86564e52849d66a9bae61104ab1d1c547a0d28e3a4103d0cb9"),
		},
		{
			setExtra(PosHeader, headerExtra{
				ActivityProof:      activityProof,
				ActivityProofRound: uint64(7),
			}),
			common.HexToHash("0x3b093edc9a33bb5d22de9252861221880fd260cc889b461f02b7ac02116a199d"),
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
