// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package params

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"net"

	"github.com/autonity/autonity/crypto/blst"

	"github.com/autonity/autonity/common/math"
	"github.com/autonity/autonity/p2p/enode"

	"golang.org/x/crypto/sha3"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/crypto"
)

// Genesis hashes to enforce below configs on.
var (
	MainnetGenesisHash = common.HexToHash("0xd4e56740f876aef8c010b86a40d5f56745a118d0906a34e69aec8c0db1cb8fa3")
	RopstenGenesisHash = common.HexToHash("0x41941023680923e0fe4d74a34bdac8141f2540e3ae90623718e47d66d1ca4a2d")
	SepoliaGenesisHash = common.HexToHash("0x25a5cc106eea7138acab33231d7160d69cb777ee0c2c553fcddf5138993e6dd9")
	RinkebyGenesisHash = common.HexToHash("0x6341fd3daf94b748c72ced5a5b26028f2474f5f00d824504e4fa37a75767e177")
	GoerliGenesisHash  = common.HexToHash("0xbf7e331f7f7c1dd2e05159666b3bf8bc7a8a3a9eb1d518969eab529dd9b88c1a")
)

// TrustedCheckpoints associates each known checkpoint with the genesis hash of
// the chain it belongs to.
var TrustedCheckpoints = map[common.Hash]*TrustedCheckpoint{
	MainnetGenesisHash: MainnetTrustedCheckpoint,
	RopstenGenesisHash: RopstenTrustedCheckpoint,
	SepoliaGenesisHash: SepoliaTrustedCheckpoint,
	RinkebyGenesisHash: RinkebyTrustedCheckpoint,
	GoerliGenesisHash:  GoerliTrustedCheckpoint,
}

// CheckpointOracles associates each known checkpoint oracles with the genesis hash of
// the chain it belongs to.
var CheckpointOracles = map[common.Hash]*CheckpointOracleConfig{
	MainnetGenesisHash: MainnetCheckpointOracle,
	RopstenGenesisHash: RopstenCheckpointOracle,
	RinkebyGenesisHash: RinkebyCheckpointOracle,
	GoerliGenesisHash:  GoerliCheckpointOracle,
}

var (
	NtnPrecision = big.NewInt(1_000_000_000_000_000_000)
	Ntn1         = new(big.Int).Mul(big.NewInt(1), NtnPrecision)
	Ntn10000     = new(big.Int).Mul(big.NewInt(10_000), NtnPrecision)

	// PiccadillyChainConfig contains the chain parameters to run a node on the Piccaddilly test network.
	PiccadillyChainConfig = &ChainConfig{
		ChainID:                 big.NewInt(65_100_002),
		HomesteadBlock:          common.Big0,
		DAOForkBlock:            common.Big0,
		DAOForkSupport:          true,
		EIP150Block:             common.Big0,
		EIP150Hash:              common.Hash{},
		EIP155Block:             common.Big0,
		EIP158Block:             common.Big0,
		ByzantiumBlock:          common.Big0,
		ConstantinopleBlock:     common.Big0,
		PetersburgBlock:         common.Big0,
		IstanbulBlock:           common.Big0,
		MuirGlacierBlock:        common.Big0,
		BerlinBlock:             common.Big0,
		LondonBlock:             common.Big0,
		ArrowGlacierBlock:       common.Big0,
		MergeForkBlock:          nil,
		TerminalTotalDifficulty: nil,
		Ethash:                  nil,
		AutonityContractConfig: &AutonityContractGenesis{
			MinBaseFee:       500_000_000,
			EpochPeriod:      30 * 60,
			UnbondingPeriod:  6 * 60 * 60,
			BlockPeriod:      1,
			MaxCommitteeSize: 9,
			Operator:         common.HexToAddress("0xd32C0812Fa1296F082671D5Be4CbB6bEeedC2397"),
			Treasury:         common.HexToAddress("0xF74c34Fed10cD9518293634C6f7C12638a808Ad5"),
			TreasuryFee:      10_000_000_000_000_000,
			DelegationRate:   1000,
			Validators: []*Validator{
				{
					Treasury:      common.HexToAddress("0x75474aC55768fAb6fE092191eea8016b955072F5"),
					OracleAddress: common.HexToAddress("0x6c5AE53a803796D788E917D1fE919BfC8B56d2E6"),
					ConsensusKey:  common.Hex2Bytes("a3aa75e42e99275f7d7985538fedc06e7f128b138a5311702afc0dc129484763645c40c36fdd97ff0d0293b00a031714"),
					Enode:         "enode://772248dfe1af5f77e0efc0510e83364bfad55cbd6d3e276f3bd0b4ddec6472aa98645655fd80bbf049ba3da18d219ab30a68fcb98da8e06dd42863dd0356cc95@35.242.168.170:30303",
					BondedStake:   Ntn1,
				}, {
					Treasury:      common.HexToAddress("0x821BC352E77D885906B47001863f75e15C114f70"),
					OracleAddress: common.HexToAddress("0x7C056299014D2F6f2e506ef1A4F89c94AAca004e"),
					ConsensusKey:  common.Hex2Bytes("8a7474c5d53279bd21b8e0d0475ca6cd868155ac16d67d22c15eaee75f87101d9b329f4c2e2da52934a845fc4567254a"),
					Enode:         "enode://22f696529d7874ca66d177c2c272600c3d1f2f7111d02140c462a8cbc789f5f8968c2ce57a5aac1373ef17bf3fc67d155b54691d1413516459824067e13750a4@34.92.27.46:30303",
					BondedStake:   Ntn1,
				}, {
					Treasury:      common.HexToAddress("0x59e2EE43e1950a348ab3CC9b6Eb847C019c2CFb1"),
					OracleAddress: common.HexToAddress("0x7615a195832843798DDb75B11CDA87D61C5F794D"),
					ConsensusKey:  common.Hex2Bytes("866610ce1e9fc5871c7b99ee92e0aa3485c7b157bddda85518edf68aea68fe1ec1ccd7bb51064e3b742719c0cf562409"),
					Enode:         "enode://24b2655b0434d1af4e2329cababf38963cab8a154e0b8c9748e75c85d10d7dab5032af7a41f3ec06dd1a7d3d306f1edee5dc46dad7a2858b80ebb56e5fa24925@34.233.111.193:30303",
					BondedStake:   Ntn1,
				}, {
					Treasury:      common.HexToAddress("0xE052e6B778C2b9591548217C3372897513175684"),
					OracleAddress: common.HexToAddress("0xaDF6a07D1f9D52bd9a464C8e3B1B670FD4b2cC55"),
					ConsensusKey:  common.Hex2Bytes("b2eed3ac8ec307e0862cdf7b8435049a98d8c83721a6066308ec57a1d9ec7f9c4e4725872dbfc96a83d102a77b430eb0"),
					Enode:         "enode://a20e27effd92dc11e7340e96a6f2908124ea363e6b68af34cad2a46a9ffdc6f5d4f516acec7f98949cc25955269f7842dc513444902c21239155de7e70b86a87@65.109.160.27:30303",
					BondedStake:   Ntn1,
				}, {
					Treasury:      common.HexToAddress("0xe33cbaA5A4E408560b229bDF1Cb58CDFD33c1C0f"),
					OracleAddress: common.HexToAddress("0xcfC6f79aa1f47b27C219E5Dbb79abA034f068b0a"),
					ConsensusKey:  common.Hex2Bytes("933c2786e202b51b04d64817f2d3d57b99fccf106edda8c53d0c66b9a356b52515e9637ddf5da6a172bb09bc5ffdb764"),
					Enode:         "enode://a2ea938a325381c7b163e7a3ca1a63fcfd927a81cadcf86551ad29f2f3ed05ef06f0b3a5d10ca932d0b85b3cf9a7c7956bf5398a2c9322f941817c92f9f62105@37.252.184.235:30303",
					BondedStake:   Ntn1,
				}, {
					Treasury:      common.HexToAddress("0xa7705415B06e76075982C7e8BaCa4aAF3D6d0f28"),
					OracleAddress: common.HexToAddress("0x469C4A81d86461d1E85be1741b1a3152e5BB9b98"),
					ConsensusKey:  common.Hex2Bytes("86f21aa126b2bae0aca1926ae5ed55f2d0207917ca6995bd48755f809868b8903e56e4cfe48bee22224e854f2e7ffb2e"),
					Enode:         "enode://46f4abe3aeca863ce3a1b4a2b2fce3112476ca75a20039ef4bad78e1a2171ae36404d74b08a0c5a8720e2548d296d37e0b92062c096801b3f6d2d86e4e9da2f2@46.4.32.57:30303",
					BondedStake:   Ntn1,
				}, {
					Treasury:      common.HexToAddress("0xB695CEe3080A58fa021b5e7d721dD1eC391983E5"),
					OracleAddress: common.HexToAddress("0x1e83A988b4fdCD21849a2b16f6db234ab1473bab"),
					ConsensusKey:  common.Hex2Bytes("a71653e690e47a3829971c8d738e54fb37f13e102759175eb04d63a6fa6f1f6873e0c77e168892e387861c3b277ee3f6"),
					Enode:         "enode://84c9a23b75bcd0252e0b361f6962a9f360d38f4fe5206cfb2d907074de877edbb1b810fd9cecf2fa64aa6ec4f7816a7f238650d489eaa82d68e8660769c6763d@51.91.220.174:30303",
					BondedStake:   Ntn1,
				}, {
					Treasury:      common.HexToAddress("0x18E1a1a0B942519e12cC524d9Acc573C7c56b12a"),
					OracleAddress: common.HexToAddress("0x9A6fc51c24D734dfc78A942d70Cd8575b89ABb9D"),
					ConsensusKey:  common.Hex2Bytes("b1782f20a3935c788352e15aec6ba7afded381170a71e6220f8b6fcad5adff168da5c32eb1bff98315168f26c837e94c"),
					Enode:         "enode://11dd1e9d4a68fb07e4cbd60d225c6ffea45852ac3d4e17df3a086a7d27ee05698922e7474db4dbcef14a11e3dd44bf66a52160610bd43a890fdc1bc8a2f51393@65.109.69.239:30303",
					BondedStake:   Ntn1,
				}, {
					Treasury:      common.HexToAddress("0x2F8cEa90D5d8Eff154fe85841C771Ac33bB6EEb1"),
					OracleAddress: common.HexToAddress("0xf41d3430ab8fb0c06c5bbf9efca587c33cdcfae3"),
					ConsensusKey:  common.Hex2Bytes("93d58b3114509a592c321250cb5a512f42aa19e59dfa29aec90a262fa0dbb9364e138caff57bc4efe5fe27c05bfe6d4a"),
					Enode:         "enode://700ae526623b87a748acf278cee299d970ccde4e4d6e7aa7685f4a550500b6e53b84892e37c2c10516673f45253fcb824d8e1836ee91a92a16b66b85b8000642@93.115.25.90:30303",
					BondedStake:   Ntn1,
				}},
		},
		OracleContractConfig: &OracleContractGenesis{
			VotePeriod: OracleVotePeriod,
			Symbols:    OracleInitialSymbols,
		},
		ASM: AsmConfig{
			ACUContractConfig:           DefaultAcuContractGenesis,
			StabilizationContractConfig: DefaultStabilizationGenesis,
			SupplyControlConfig: &SupplyControlGenesis{
				InitialAllocation: (*math.HexOrDecimal256)(new(big.Int).Sub(
					new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil), // 2^256
					new(big.Int).Mul(big.NewInt(16), big.NewInt(Ether)),   // 16 * 10^18
				)),
			},
		},
		AccountabilityConfig: DefaultAccountabilityConfig,
	}

	// BakerlooChainConfig todo: ask Raj to generate validator key for validators in the BakerlooChainConfig
	// BakerlooChainConfig contains the chain parameters to run a node on the Bakerloo test network.
	BakerlooChainConfig = &ChainConfig{
		ChainID:                 big.NewInt(65_010_002),
		HomesteadBlock:          common.Big0,
		DAOForkBlock:            common.Big0,
		DAOForkSupport:          true,
		EIP150Block:             common.Big0,
		EIP150Hash:              common.Hash{},
		EIP155Block:             common.Big0,
		EIP158Block:             common.Big0,
		ByzantiumBlock:          common.Big0,
		ConstantinopleBlock:     common.Big0,
		PetersburgBlock:         common.Big0,
		IstanbulBlock:           common.Big0,
		MuirGlacierBlock:        common.Big0,
		BerlinBlock:             common.Big0,
		LondonBlock:             common.Big0,
		ArrowGlacierBlock:       common.Big0,
		MergeForkBlock:          nil,
		TerminalTotalDifficulty: nil,
		Ethash:                  nil,
		AutonityContractConfig: &AutonityContractGenesis{
			MinBaseFee:       500_000_000,
			EpochPeriod:      30 * 60,
			UnbondingPeriod:  6 * 60 * 60,
			BlockPeriod:      1,
			MaxCommitteeSize: 50,
			Operator:         common.HexToAddress("0x293039dDC627B1dF9562380c0E5377848F94325A"),
			Treasury:         common.HexToAddress("0x7f1B212dcDc119a395Ec2B245ce86e9eE551043E"),
			TreasuryFee:      10_000_000_000_000_000,
			DelegationRate:   1000,
			Validators: []*Validator{{
				Treasury:      common.HexToAddress("0x3e08FEc6ABaf669BD8Da54abEe30b2B8B5024013"),
				OracleAddress: common.HexToAddress("0x4D8387E38F42084aa24CE7DA137222786fF23A3E"),
				ConsensusKey:  common.Hex2Bytes("89d04918f676554704f833d193b0819ef171784309bc4573cc2eeaaef78e3881bba3b7cc9bc8b8c3ff5678c355626a7d"),
				Enode:         "enode://ebd47c61fa3b0125240654f7b73abeb188dd8c954385e2b75dd0b46efc341b801ae4593be77a954120bcff97a5282ff8f9485306976f8541b2fa971066cd37be@35.246.21.247:30303",
				BondedStake:   Ntn10000,
			}, {
				Treasury:      common.HexToAddress("0xf1859D9feD50514F9D805BeC7a30623d061f40B7"),
				OracleAddress: common.HexToAddress("0x22F1e6eA5d67Bef19C6953bdBCFA03320ECd015A"),
				ConsensusKey:  common.Hex2Bytes("8df1b499f3a7acd89cd36ef96ceb931a2455751c5e676000e378bec03f2008c6a321cd833132204f46e43ee719dac5c4"),
				Enode:         "enode://1cbd580027a2d21c114c75720115c4c04524952231f66607514fd6bcb5979f20181066383a7a15f1382aa637f6470170c67f74bf1f95271f2abc4ae08c8a11fb@35.189.64.207:30303",
				BondedStake:   Ntn10000,
			}, {
				Treasury:      common.HexToAddress("0x1B441084736B80f273e498E646b0bEA86B4eC6AB"),
				OracleAddress: common.HexToAddress("0xC81B686402395A83938452DF8398DA9b2649281A"),
				ConsensusKey:  common.Hex2Bytes("a648a602bf8ad04b07cf2cbffcb70eccb2bfe9500dca99c599998f5b8cae950e99ab0e7247a03c1ca731fe7c6fc257f8"),
				Enode:         "enode://28136dd78f976e8fbcbd2c8222cee3f6ac45ac01ceea4a43946d476c03f4971cea23d8c30f0e1a2cbcfde5310a0117ae36ba15c6028e036f637141f330040f62@34.105.163.137:30303",
				BondedStake:   Ntn10000,
			}, {

				Treasury:      common.HexToAddress("0xB5C49d50470743D8dE43bB6822AC4505E64648Da"),
				OracleAddress: common.HexToAddress("0x4c35544931E2Cea6eD41102197685704917F72C3"),
				ConsensusKey:  common.Hex2Bytes("897bef2c488b7bf696d8f4595790233e7ced50bb10bfb33c3e82d8526ce147a70bc6152842eeeeacbc044ebe43c64b27"),
				Enode:         "enode://e4dacdb0170bc1baaf7a90935c8d35e75aeec639087320b19ffdd6eb8e6f3b967dd696b1c2f3f99a6a1fea003f7c2a0bb9d8228c06b31ae6c5fec863b8745a2d@35.177.8.113:30303",
				BondedStake:   Ntn10000,
			}, {
				Treasury:      common.HexToAddress("0x31e1dE659A26F7638FAaFEfD94D47258FE361823"),
				OracleAddress: common.HexToAddress("0xC0bB231711470a92DE3B57DE2Ca04727048f5580"),
				ConsensusKey:  common.Hex2Bytes("b6c83a5bfd23f6b2043e89bf5a774e41428151b913dd842efac094f04ac41bec8f5110e313c38ccc115f916108769c40"),
				Enode:         "enode://07ea9eaa469d07695d6855089ad1e0fda35933b779f5907c6edbb77365eef2297b56039fa021e31183675c201855bd95275852d4b1ff9251cce558d1a3611240@35.179.46.181:30303",
				BondedStake:   Ntn10000,
			}, {
				Treasury:      common.HexToAddress("0xe22617BD2a4e1Fe3938F84060D8a6be7A18a2ef9"),
				OracleAddress: common.HexToAddress("0x82C3E23Aa626Ca1556938bCA38f52B329A99b9d8"),
				ConsensusKey:  common.Hex2Bytes("80293a218b081704b25d6b1db351246472b111e42e73998d34a9783f1439eba2cdd4fb832cc9d018986c59a398fd0e69"),
				Enode:         "enode://fb6f5556df35da18c1ecac878fd44a84a0cc302952e1172871dfc48c996c9940b182df374586722212ecb92a8771cafcf48b21657832316423e33f7dbc0b4e7f@3.9.98.39:30303",
				BondedStake:   Ntn10000,
			}},
		},
		OracleContractConfig: &OracleContractGenesis{
			VotePeriod: OracleVotePeriod,
			Symbols:    OracleInitialSymbols,
		},
		ASM: AsmConfig{
			ACUContractConfig:           DefaultAcuContractGenesis,
			StabilizationContractConfig: DefaultStabilizationGenesis,
			SupplyControlConfig: &SupplyControlGenesis{
				InitialAllocation: (*math.HexOrDecimal256)(new(big.Int).Sub(
					new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil), // 2^256
					new(big.Int).Mul(big.NewInt(16), big.NewInt(Ether)),   // 16 * 10^18
				)),
			},
		},
		AccountabilityConfig: DefaultAccountabilityConfig,
	}

	// MainnetChainConfig is the chain parameters to run a node on the main network.
	MainnetChainConfig = &ChainConfig{
		ChainID:             big.NewInt(1),
		HomesteadBlock:      big.NewInt(1_150_000),
		DAOForkBlock:        big.NewInt(1_920_000),
		DAOForkSupport:      true,
		EIP150Block:         big.NewInt(2_463_000),
		EIP150Hash:          common.HexToHash("0x2086799aeebeae135c246c65021c82b4e15a2c451340993aacfd2751886514f0"),
		EIP155Block:         big.NewInt(2_675_000),
		EIP158Block:         big.NewInt(2_675_000),
		ByzantiumBlock:      big.NewInt(4_370_000),
		ConstantinopleBlock: big.NewInt(7_280_000),
		PetersburgBlock:     big.NewInt(7_280_000),
		IstanbulBlock:       big.NewInt(9_069_000),
		MuirGlacierBlock:    big.NewInt(9_200_000),
		BerlinBlock:         big.NewInt(12_244_000),
		LondonBlock:         big.NewInt(12_965_000),
		ArrowGlacierBlock:   big.NewInt(13_773_000),
		Ethash:              new(EthashConfig),
	}

	// MainnetTrustedCheckpoint contains the light client trusted checkpoint for the main network.
	MainnetTrustedCheckpoint = &TrustedCheckpoint{
		SectionIndex: 413,
		SectionHead:  common.HexToHash("0x8aa8e64ceadcdc5f23bc41d2acb7295a261a5cf680bb00a34f0e01af08200083"),
		CHTRoot:      common.HexToHash("0x008af584d385a2610706c5a439d39f15ddd4b691c5d42603f65ae576f703f477"),
		BloomRoot:    common.HexToHash("0x5a081af71a588f4d90bced242545b08904ad4fb92f7effff2ceb6e50e6dec157"),
	}

	// MainnetCheckpointOracle contains a set of configs for the main network oracle.
	MainnetCheckpointOracle = &CheckpointOracleConfig{
		Address: common.HexToAddress("0x9a9070028361F7AAbeB3f2F2Dc07F82C4a98A02a"),
		Signers: []common.Address{
			common.HexToAddress("0x1b2C260efc720BE89101890E4Db589b44E950527"), // Peter
			common.HexToAddress("0x78d1aD571A1A09D60D9BBf25894b44e4C8859595"), // Martin
			common.HexToAddress("0x286834935f4A8Cfb4FF4C77D5770C2775aE2b0E7"), // Zsolt
			common.HexToAddress("0xb86e2B0Ab5A4B1373e40c51A7C712c70Ba2f9f8E"), // Gary
			common.HexToAddress("0x0DF8fa387C602AE62559cC4aFa4972A7045d6707"), // Guillaume
		},
		Threshold: 2,
	}

	// RopstenChainConfig contains the chain parameters to run a node on the Ropsten test network.
	RopstenChainConfig = &ChainConfig{
		ChainID:             big.NewInt(3),
		HomesteadBlock:      big.NewInt(0),
		DAOForkBlock:        nil,
		DAOForkSupport:      true,
		EIP150Block:         big.NewInt(0),
		EIP150Hash:          common.HexToHash("0x41941023680923e0fe4d74a34bdac8141f2540e3ae90623718e47d66d1ca4a2d"),
		EIP155Block:         big.NewInt(10),
		EIP158Block:         big.NewInt(10),
		ByzantiumBlock:      big.NewInt(1_700_000),
		ConstantinopleBlock: big.NewInt(4_230_000),
		PetersburgBlock:     big.NewInt(4_939_394),
		IstanbulBlock:       big.NewInt(6_485_846),
		MuirGlacierBlock:    big.NewInt(7_117_117),
		BerlinBlock:         big.NewInt(9_812_189),
		LondonBlock:         big.NewInt(10_499_401),
		Ethash:              new(EthashConfig),
	}

	// RopstenTrustedCheckpoint contains the light client trusted checkpoint for the Ropsten test network.
	RopstenTrustedCheckpoint = &TrustedCheckpoint{
		SectionIndex: 346,
		SectionHead:  common.HexToHash("0xafa0384ebd13a751fb7475aaa7fc08ac308925c8b2e2195bca2d4ab1878a7a84"),
		CHTRoot:      common.HexToHash("0x522ae1f334bfa36033b2315d0b9954052780700b69448ecea8d5877e0f7ee477"),
		BloomRoot:    common.HexToHash("0x4093fd53b0d2cc50181dca353fe66f03ae113e7cb65f869a4dfb5905de6a0493"),
	}

	// RopstenCheckpointOracle contains a set of configs for the Ropsten test network oracle.
	RopstenCheckpointOracle = &CheckpointOracleConfig{
		Address: common.HexToAddress("0xEF79475013f154E6A65b54cB2742867791bf0B84"),
		Signers: []common.Address{
			common.HexToAddress("0x32162F3581E88a5f62e8A61892B42C46E2c18f7b"), // Peter
			common.HexToAddress("0x78d1aD571A1A09D60D9BBf25894b44e4C8859595"), // Martin
			common.HexToAddress("0x286834935f4A8Cfb4FF4C77D5770C2775aE2b0E7"), // Zsolt
			common.HexToAddress("0xb86e2B0Ab5A4B1373e40c51A7C712c70Ba2f9f8E"), // Gary
			common.HexToAddress("0x0DF8fa387C602AE62559cC4aFa4972A7045d6707"), // Guillaume
		},
		Threshold: 2,
	}

	// SepoliaChainConfig contains the chain parameters to run a node on the Sepolia test network.
	SepoliaChainConfig = &ChainConfig{
		ChainID:             big.NewInt(11155111),
		HomesteadBlock:      big.NewInt(0),
		DAOForkBlock:        nil,
		DAOForkSupport:      true,
		EIP150Block:         big.NewInt(0),
		EIP155Block:         big.NewInt(0),
		EIP158Block:         big.NewInt(0),
		ByzantiumBlock:      big.NewInt(0),
		ConstantinopleBlock: big.NewInt(0),
		PetersburgBlock:     big.NewInt(0),
		IstanbulBlock:       big.NewInt(0),
		MuirGlacierBlock:    big.NewInt(0),
		BerlinBlock:         big.NewInt(0),
		LondonBlock:         big.NewInt(0),
		Ethash:              new(EthashConfig),
	}

	// SepoliaTrustedCheckpoint contains the light client trusted checkpoint for the Sepolia test network.
	SepoliaTrustedCheckpoint = &TrustedCheckpoint{
		SectionIndex: 1,
		SectionHead:  common.HexToHash("0x5dde65e28745b10ff9e9b86499c3a3edc03587b27a06564a4342baf3a37de869"),
		CHTRoot:      common.HexToHash("0x042a0d914f7baa4f28f14d12291e5f346e88c5b9d95127bf5422a8afeacd27e8"),
		BloomRoot:    common.HexToHash("0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421"),
	}

	// RinkebyChainConfig contains the chain parameters to run a node on the Rinkeby test network.
	RinkebyChainConfig = &ChainConfig{
		ChainID:             big.NewInt(4),
		HomesteadBlock:      big.NewInt(1),
		DAOForkBlock:        nil,
		DAOForkSupport:      true,
		EIP150Block:         big.NewInt(2),
		EIP150Hash:          common.HexToHash("0x9b095b36c15eaf13044373aef8ee0bd3a382a5abb92e402afa44b8249c3a90e9"),
		EIP155Block:         big.NewInt(3),
		EIP158Block:         big.NewInt(3),
		ByzantiumBlock:      big.NewInt(1035301),
		ConstantinopleBlock: big.NewInt(3660663),
		PetersburgBlock:     big.NewInt(4321234),
		IstanbulBlock:       big.NewInt(5435345),
		MuirGlacierBlock:    nil,
		BerlinBlock:         big.NewInt(8_290_928),
		LondonBlock:         big.NewInt(8_897_988),
		ArrowGlacierBlock:   nil,
	}

	// RinkebyTrustedCheckpoint contains the light client trusted checkpoint for the Rinkeby test network.
	RinkebyTrustedCheckpoint = &TrustedCheckpoint{
		SectionIndex: 292,
		SectionHead:  common.HexToHash("0x4185c2f1bb85ecaa04409d1008ff0761092ea2e94e8a71d64b1a5abc37b81414"),
		CHTRoot:      common.HexToHash("0x03b0191e6140effe0b88bb7c97bfb794a275d3543cb3190662fb72d9beea423c"),
		BloomRoot:    common.HexToHash("0x3d5f6edccc87536dcbc0dd3aae97a318205c617dd3957b4261470c71481629e2"),
	}

	// RinkebyCheckpointOracle contains a set of configs for the Rinkeby test network oracle.
	RinkebyCheckpointOracle = &CheckpointOracleConfig{
		Address: common.HexToAddress("0xebe8eFA441B9302A0d7eaECc277c09d20D684540"),
		Signers: []common.Address{
			common.HexToAddress("0xd9c9cd5f6779558b6e0ed4e6acf6b1947e7fa1f3"), // Peter
			common.HexToAddress("0x78d1aD571A1A09D60D9BBf25894b44e4C8859595"), // Martin
			common.HexToAddress("0x286834935f4A8Cfb4FF4C77D5770C2775aE2b0E7"), // Zsolt
			common.HexToAddress("0xb86e2B0Ab5A4B1373e40c51A7C712c70Ba2f9f8E"), // Gary
		},
		Threshold: 2,
	}

	// GoerliChainConfig contains the chain parameters to run a node on the Görli test network.
	GoerliChainConfig = &ChainConfig{
		ChainID:             big.NewInt(5),
		HomesteadBlock:      big.NewInt(0),
		DAOForkBlock:        nil,
		DAOForkSupport:      true,
		EIP150Block:         big.NewInt(0),
		EIP155Block:         big.NewInt(0),
		EIP158Block:         big.NewInt(0),
		ByzantiumBlock:      big.NewInt(0),
		ConstantinopleBlock: big.NewInt(0),
		PetersburgBlock:     big.NewInt(0),
		IstanbulBlock:       big.NewInt(1_561_651),
		MuirGlacierBlock:    nil,
		BerlinBlock:         big.NewInt(4_460_644),
		LondonBlock:         big.NewInt(5_062_605),
		ArrowGlacierBlock:   nil,
	}

	// GoerliTrustedCheckpoint contains the light client trusted checkpoint for the Görli test network.
	GoerliTrustedCheckpoint = &TrustedCheckpoint{
		SectionIndex: 176,
		SectionHead:  common.HexToHash("0x2de018858528434f93adb40b1f03f2304a86d31b4ef2b1f930da0134f5c32427"),
		CHTRoot:      common.HexToHash("0x8c17e497d38088321c147abe4acbdfb3c0cab7d7a2b97e07404540f04d12747e"),
		BloomRoot:    common.HexToHash("0x02a41b6606bd3f741bd6ae88792d75b1ad8cf0ea5e28fbaa03bc8b95cbd20034"),
	}

	// GoerliCheckpointOracle contains a set of configs for the Goerli test network oracle.
	GoerliCheckpointOracle = &CheckpointOracleConfig{
		Address: common.HexToAddress("0x18CA0E045F0D772a851BC7e48357Bcaab0a0795D"),
		Signers: []common.Address{
			common.HexToAddress("0x4769bcaD07e3b938B7f43EB7D278Bc7Cb9efFb38"), // Peter
			common.HexToAddress("0x78d1aD571A1A09D60D9BBf25894b44e4C8859595"), // Martin
			common.HexToAddress("0x286834935f4A8Cfb4FF4C77D5770C2775aE2b0E7"), // Zsolt
			common.HexToAddress("0xb86e2B0Ab5A4B1373e40c51A7C712c70Ba2f9f8E"), // Gary
			common.HexToAddress("0x0DF8fa387C602AE62559cC4aFa4972A7045d6707"), // Guillaume
		},
		Threshold: 2,
	}

	// AllEthashProtocolChanges contains every protocol change (EIPs) introduced
	// and accepted by the Ethereum core developers into the Ethash consensus.
	//
	// This configuration is intentionally not using keyed fields to force anyone
	// adding flags to the config to also have to set these fields.
	AllEthashProtocolChanges = &ChainConfig{big.NewInt(1337), big.NewInt(0), nil, false, big.NewInt(0), common.Hash{}, big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), nil, nil, nil, nil, nil, nil, new(EthashConfig), nil, nil, nil, AsmConfig{}, false}

	TestNodeKeys = []string{
		"b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291",
		"a4b489752489e0f47e410b8e8cbb1ac1b56770d202ffd45b346ca8355c602c91",
		"3cfb95a9d463ee29b8470742a9718ef3298e04b367b7c796fe67cc693587d746",
		"193f20ab2451ea4e4ba0aaf83b4cff335df716247359c98562f8da68e07f7c1d",
	}
	TestConsensusKeys = []string{
		"0afbb1b94ac30db9e145eb30ee6b64d1996a31279e50005b2a470b18dae82bcb",
		"3f0e004faa78fde4627834285760652f71a85942f10b354b67dc55ea494c4e8f",
		"497409f62556016749f7518d154b01baaa0c6a34b1694a3ed55bbffed9a6f30d",
		"5d1a359c9f81b2b199e4cd972990ddf101d03ab5e44d7313b4da06d7dfc06b87",
	}
	TestValidatorBase = Validator{
		BondedStake: Ntn10000,
	}
	TestValidatorConsensusKey, _ = blst.SecretKeyFromHex("0afbb1b94ac30db9e145eb30ee6b64d1996a31279e50005b2a470b18dae82bcb")

	TestAutonityContractConfig = &AutonityContractGenesis{
		MaxCommitteeSize: 21,
		BlockPeriod:      1,
		UnbondingPeriod:  120,
		EpochPeriod:      30,
		DelegationRate:   1200, // 12%
		Treasury:         common.Address{120},
		TreasuryFee:      1500000000000000, // 0.15%,
		MinBaseFee:       InitialBaseFee,
		Operator:         common.Address{},
	}

	TestChainConfig = &ChainConfig{
		big.NewInt(1337),
		big.NewInt(0),
		nil,
		false,
		big.NewInt(0),
		common.Hash{},
		big.NewInt(0),
		big.NewInt(0),
		big.NewInt(0),
		big.NewInt(0),
		big.NewInt(0),
		big.NewInt(0),
		big.NewInt(0),
		big.NewInt(0),
		big.NewInt(0),
		big.NewInt(0),
		nil,
		nil,
		new(EthashConfig),
		TestAutonityContractConfig,
		DefaultAccountabilityConfig,
		DefaultGenesisOracleConfig,
		AsmConfig{
			ACUContractConfig:           DefaultAcuContractGenesis,
			StabilizationContractConfig: DefaultStabilizationGenesis,
			SupplyControlConfig:         DefaultSupplyControlGenesis,
		},
		false,
	}
)

func init() {
	// Setup the validators in TestAutonityContractConfig
	for i := range TestNodeKeys {
		validator := TestValidatorBase
		nodeKey, _ := crypto.HexToECDSA(TestNodeKeys[i])
		address := crypto.PubkeyToAddress(nodeKey.PublicKey)
		validator.NodeAddress = &address
		validator.Treasury = address
		validator.OracleAddress = address
		validator.Enode = enode.NewV4(&nodeKey.PublicKey, net.ParseIP("0.0.0.0"), 0, 0).URLv4()
		consensusKey, _ := blst.SecretKeyFromHex(TestConsensusKeys[i])
		validator.ConsensusKey = consensusKey.PublicKey().Marshal()
		TestAutonityContractConfig.Validators = append(TestAutonityContractConfig.Validators, &validator)
	}
	TestAutonityContractConfig.Prepare()
	PiccadillyChainConfig.AutonityContractConfig.Prepare()
	BakerlooChainConfig.AutonityContractConfig.Prepare()
}

// TrustedCheckpoint represents a set of post-processed trie roots (CHT and
// BloomTrie) associated with the appropriate section index and head hash. It is
// used to start light syncing from this checkpoint and avoid downloading the
// entire header chain while still being able to securely access old headers/logs.
type TrustedCheckpoint struct {
	SectionIndex uint64      `json:"sectionIndex"`
	SectionHead  common.Hash `json:"sectionHead"`
	CHTRoot      common.Hash `json:"chtRoot"`
	BloomRoot    common.Hash `json:"bloomRoot"`
}

// HashEqual returns an indicator comparing the itself hash with given one.
func (c *TrustedCheckpoint) HashEqual(hash common.Hash) bool {
	if c.Empty() {
		return hash == common.Hash{}
	}
	return c.Hash() == hash
}

// Hash returns the hash of checkpoint's four key fields(index, sectionHead, chtRoot and bloomTrieRoot).
func (c *TrustedCheckpoint) Hash() common.Hash {
	var sectionIndex [8]byte
	binary.BigEndian.PutUint64(sectionIndex[:], c.SectionIndex)

	w := sha3.NewLegacyKeccak256()
	w.Write(sectionIndex[:])
	w.Write(c.SectionHead[:])
	w.Write(c.CHTRoot[:])
	w.Write(c.BloomRoot[:])

	var h common.Hash
	w.Sum(h[:0])
	return h
}

// Empty returns an indicator whether the checkpoint is regarded as empty.
func (c *TrustedCheckpoint) Empty() bool {
	return c.SectionHead == (common.Hash{}) || c.CHTRoot == (common.Hash{}) || c.BloomRoot == (common.Hash{})
}

// CheckpointOracleConfig represents a set of checkpoint contract(which acts as an oracle)
// config which used for light client checkpoint syncing.
type CheckpointOracleConfig struct {
	Address   common.Address   `json:"address"`
	Signers   []common.Address `json:"signers"`
	Threshold uint64           `json:"threshold"`
}

// ChainConfig is the core config which determines the blockchain settings.
//
// ChainConfig is stored in the database on a per block basis. This means
// that any network, identified by its genesis block, can have its own
// set of configuration options.
type ChainConfig struct {
	ChainID *big.Int `json:"chainId"` // chainId identifies the current chain and is used for replay protection

	HomesteadBlock *big.Int `json:"homesteadBlock,omitempty"` // Homestead switch block (nil = no fork, 0 = already homestead)

	DAOForkBlock   *big.Int `json:"daoForkBlock,omitempty"`   // TheDAO hard-fork switch block (nil = no fork)
	DAOForkSupport bool     `json:"daoForkSupport,omitempty"` // Whether the nodes supports or opposes the DAO hard-fork

	// EIP150 implements the Gas price changes (https://github.com/ethereum/EIPs/issues/150)
	EIP150Block *big.Int    `json:"eip150Block,omitempty"` // EIP150 HF block (nil = no fork)
	EIP150Hash  common.Hash `json:"eip150Hash,omitempty"`  // EIP150 HF hash (needed for header only clients as only gas pricing changed)

	EIP155Block *big.Int `json:"eip155Block,omitempty"` // EIP155 HF block
	EIP158Block *big.Int `json:"eip158Block,omitempty"` // EIP158 HF block

	ByzantiumBlock      *big.Int `json:"byzantiumBlock,omitempty"`      // Byzantium switch block (nil = no fork, 0 = already on byzantium)
	ConstantinopleBlock *big.Int `json:"constantinopleBlock,omitempty"` // Constantinople switch block (nil = no fork, 0 = already activated)
	PetersburgBlock     *big.Int `json:"petersburgBlock,omitempty"`     // Petersburg switch block (nil = same as Constantinople)
	IstanbulBlock       *big.Int `json:"istanbulBlock,omitempty"`       // Istanbul switch block (nil = no fork, 0 = already on istanbul)
	MuirGlacierBlock    *big.Int `json:"muirGlacierBlock,omitempty"`    // Eip-2384 (bomb delay) switch block (nil = no fork, 0 = already activated)
	BerlinBlock         *big.Int `json:"berlinBlock,omitempty"`         // Berlin switch block (nil = no fork, 0 = already on berlin)
	LondonBlock         *big.Int `json:"londonBlock,omitempty"`         // London switch block (nil = no fork, 0 = already on london)
	ArrowGlacierBlock   *big.Int `json:"arrowGlacierBlock,omitempty"`   // Eip-4345 (bomb delay) switch block (nil = no fork, 0 = already activated)
	MergeForkBlock      *big.Int `json:"mergeForkBlock,omitempty"`      // EIP-3675 (TheMerge) switch block (nil = no fork, 0 = already in merge proceedings)

	// TerminalTotalDifficulty is the amount of total difficulty reached by
	// the network that triggers the consensus upgrade.
	TerminalTotalDifficulty *big.Int `json:"terminalTotalDifficulty,omitempty"`

	// Various consensus engines
	Ethash                 *EthashConfig            `json:"ethash,omitempty"`
	AutonityContractConfig *AutonityContractGenesis `json:"autonity,omitempty"`
	AccountabilityConfig   *AccountabilityGenesis   `json:"accountability,omitempty"`
	OracleContractConfig   *OracleContractGenesis   `json:"oracle,omitempty"`

	ASM AsmConfig `json:"asm,omitempty"`

	// true if run in testmode, false by default
	TestMode bool `json:"testMode,omitempty"`
}

type AsmConfig struct {
	ACUContractConfig           *AcuContractGenesis           `json:"acu,omitempty"`
	StabilizationContractConfig *StabilizationContractGenesis `json:"stabilization,omitempty"`
	SupplyControlConfig         *SupplyControlGenesis         `json:"supplyControl,omitempty"`
}

// EthashConfig is the consensus engine configs for proof-of-work based sealing.
type EthashConfig struct{}

// String implements the stringer interface, returning the consensus engine details.
func (c *EthashConfig) String() string {
	return "ethash"
}

// String implements the fmt.Stringer interface.
func (c *ChainConfig) String() string {
	var engine interface{}
	switch {
	case c.Ethash != nil:
		engine = c.Ethash
	default:
		engine = "Tendermint"
	}
	return fmt.Sprintf("{ChainID: %v Homestead: %v DAO: %v DAOSupport: %v EIP150: %v EIP155: %v EIP158: %v Byzantium: %v Constantinople: %v Petersburg: %v Istanbul: %v, Muir Glacier: %v, Berlin: %v, London: %v, Arrow Glacier: %v, MergeFork: %v, Engine: %v}",
		c.ChainID,
		c.HomesteadBlock,
		c.DAOForkBlock,
		c.DAOForkSupport,
		c.EIP150Block,
		c.EIP155Block,
		c.EIP158Block,
		c.ByzantiumBlock,
		c.ConstantinopleBlock,
		c.PetersburgBlock,
		c.IstanbulBlock,
		c.MuirGlacierBlock,
		c.BerlinBlock,
		c.LondonBlock,
		c.ArrowGlacierBlock,
		c.MergeForkBlock,
		engine,
	)
}

// IsHomestead returns whether num is either equal to the homestead block or greater.
func (c *ChainConfig) IsHomestead(num *big.Int) bool {
	return isForked(c.HomesteadBlock, num)
}

// IsDAOFork returns whether num is either equal to the DAO fork block or greater.
func (c *ChainConfig) IsDAOFork(num *big.Int) bool {
	return isForked(c.DAOForkBlock, num)
}

// IsEIP150 returns whether num is either equal to the EIP150 fork block or greater.
func (c *ChainConfig) IsEIP150(num *big.Int) bool {
	return isForked(c.EIP150Block, num)
}

// IsEIP155 returns whether num is either equal to the EIP155 fork block or greater.
func (c *ChainConfig) IsEIP155(num *big.Int) bool {
	return isForked(c.EIP155Block, num)
}

// IsEIP158 returns whether num is either equal to the EIP158 fork block or greater.
func (c *ChainConfig) IsEIP158(num *big.Int) bool {
	return isForked(c.EIP158Block, num)
}

// IsByzantium returns whether num is either equal to the Byzantium fork block or greater.
func (c *ChainConfig) IsByzantium(num *big.Int) bool {
	return isForked(c.ByzantiumBlock, num)
}

// IsConstantinople returns whether num is either equal to the Constantinople fork block or greater.
func (c *ChainConfig) IsConstantinople(num *big.Int) bool {
	return isForked(c.ConstantinopleBlock, num)
}

// IsMuirGlacier returns whether num is either equal to the Muir Glacier (EIP-2384) fork block or greater.
func (c *ChainConfig) IsMuirGlacier(num *big.Int) bool {
	return isForked(c.MuirGlacierBlock, num)
}

// IsPetersburg returns whether num is either
// - equal to or greater than the PetersburgBlock fork block,
// - OR is nil, and Constantinople is active
func (c *ChainConfig) IsPetersburg(num *big.Int) bool {
	return isForked(c.PetersburgBlock, num) || c.PetersburgBlock == nil && isForked(c.ConstantinopleBlock, num)
}

// IsIstanbul returns whether num is either equal to the Istanbul fork block or greater.
func (c *ChainConfig) IsIstanbul(num *big.Int) bool {
	return isForked(c.IstanbulBlock, num)
}

// IsBerlin returns whether num is either equal to the Berlin fork block or greater.
func (c *ChainConfig) IsBerlin(num *big.Int) bool {
	return isForked(c.BerlinBlock, num)
}

// IsLondon returns whether num is either equal to the London fork block or greater.
func (c *ChainConfig) IsLondon(num *big.Int) bool {
	return isForked(c.LondonBlock, num)
}

// IsArrowGlacier returns whether num is either equal to the Arrow Glacier (EIP-4345) fork block or greater.
func (c *ChainConfig) IsArrowGlacier(num *big.Int) bool {
	return isForked(c.ArrowGlacierBlock, num)
}

// CheckCompatible checks whether scheduled fork transitions have been imported
// with a mismatching chain configuration.
func (c *ChainConfig) CheckCompatible(newcfg *ChainConfig, height uint64) *ConfigCompatError {
	bhead := new(big.Int).SetUint64(height)

	// Iterate checkCompatible to find the lowest conflict.
	var lasterr *ConfigCompatError
	for {
		err := c.checkCompatible(newcfg, bhead)
		if err == nil || (lasterr != nil && err.RewindTo == lasterr.RewindTo) {
			break
		}
		lasterr = err
		bhead.SetUint64(err.RewindTo)
	}
	return lasterr
}

// CheckConfigForkOrder checks that we don't "skip" any forks, geth isn't pluggable enough
// to guarantee that forks can be implemented in a different order than on official networks
func (c *ChainConfig) CheckConfigForkOrder() error {
	type fork struct {
		name     string
		block    *big.Int
		optional bool // if true, the fork may be nil and next fork is still allowed
	}
	var lastFork fork
	for _, cur := range []fork{
		{name: "homesteadBlock", block: c.HomesteadBlock},
		{name: "daoForkBlock", block: c.DAOForkBlock, optional: true},
		{name: "eip150Block", block: c.EIP150Block},
		{name: "eip155Block", block: c.EIP155Block},
		{name: "eip158Block", block: c.EIP158Block},
		{name: "byzantiumBlock", block: c.ByzantiumBlock},
		{name: "constantinopleBlock", block: c.ConstantinopleBlock},
		{name: "petersburgBlock", block: c.PetersburgBlock},
		{name: "istanbulBlock", block: c.IstanbulBlock},
		{name: "muirGlacierBlock", block: c.MuirGlacierBlock, optional: true},
		{name: "berlinBlock", block: c.BerlinBlock},
		{name: "londonBlock", block: c.LondonBlock},
		{name: "arrowGlacierBlock", block: c.ArrowGlacierBlock, optional: true},
		{name: "mergeStartBlock", block: c.MergeForkBlock, optional: true},
	} {
		if lastFork.name != "" {
			// Next one must be higher number
			if lastFork.block == nil && cur.block != nil {
				return fmt.Errorf("unsupported fork ordering: %v not enabled, but %v enabled at %v",
					lastFork.name, cur.name, cur.block)
			}
			if lastFork.block != nil && cur.block != nil {
				if lastFork.block.Cmp(cur.block) > 0 {
					return fmt.Errorf("unsupported fork ordering: %v enabled at %v, but %v enabled at %v",
						lastFork.name, lastFork.block, cur.name, cur.block)
				}
			}
		}
		// If it was optional and not set, then ignore it
		if !cur.optional || cur.block != nil {
			lastFork = cur
		}
	}
	return nil
}

func (c *ChainConfig) checkCompatible(newcfg *ChainConfig, head *big.Int) *ConfigCompatError {
	if isForkIncompatible(c.HomesteadBlock, newcfg.HomesteadBlock, head) {
		return newCompatError("Homestead fork block", c.HomesteadBlock, newcfg.HomesteadBlock)
	}
	if isForkIncompatible(c.DAOForkBlock, newcfg.DAOForkBlock, head) {
		return newCompatError("DAO fork block", c.DAOForkBlock, newcfg.DAOForkBlock)
	}
	if c.IsDAOFork(head) && c.DAOForkSupport != newcfg.DAOForkSupport {
		return newCompatError("DAO fork support flag", c.DAOForkBlock, newcfg.DAOForkBlock)
	}
	if isForkIncompatible(c.EIP150Block, newcfg.EIP150Block, head) {
		return newCompatError("EIP150 fork block", c.EIP150Block, newcfg.EIP150Block)
	}
	if isForkIncompatible(c.EIP155Block, newcfg.EIP155Block, head) {
		return newCompatError("EIP155 fork block", c.EIP155Block, newcfg.EIP155Block)
	}
	if isForkIncompatible(c.EIP158Block, newcfg.EIP158Block, head) {
		return newCompatError("EIP158 fork block", c.EIP158Block, newcfg.EIP158Block)
	}
	if c.IsEIP158(head) && !configNumEqual(c.ChainID, newcfg.ChainID) {
		return newCompatError("EIP158 chain ID", c.EIP158Block, newcfg.EIP158Block)
	}
	if isForkIncompatible(c.ByzantiumBlock, newcfg.ByzantiumBlock, head) {
		return newCompatError("Byzantium fork block", c.ByzantiumBlock, newcfg.ByzantiumBlock)
	}
	if isForkIncompatible(c.ConstantinopleBlock, newcfg.ConstantinopleBlock, head) {
		return newCompatError("Constantinople fork block", c.ConstantinopleBlock, newcfg.ConstantinopleBlock)
	}
	if isForkIncompatible(c.PetersburgBlock, newcfg.PetersburgBlock, head) {
		// the only case where we allow Petersburg to be set in the past is if it is equal to Constantinople
		// mainly to satisfy fork ordering requirements which state that Petersburg fork be set if Constantinople fork is set
		if isForkIncompatible(c.ConstantinopleBlock, newcfg.PetersburgBlock, head) {
			return newCompatError("Petersburg fork block", c.PetersburgBlock, newcfg.PetersburgBlock)
		}
	}
	if isForkIncompatible(c.IstanbulBlock, newcfg.IstanbulBlock, head) {
		return newCompatError("Istanbul fork block", c.IstanbulBlock, newcfg.IstanbulBlock)
	}
	if isForkIncompatible(c.MuirGlacierBlock, newcfg.MuirGlacierBlock, head) {
		return newCompatError("Muir Glacier fork block", c.MuirGlacierBlock, newcfg.MuirGlacierBlock)
	}
	if isForkIncompatible(c.BerlinBlock, newcfg.BerlinBlock, head) {
		return newCompatError("Berlin fork block", c.BerlinBlock, newcfg.BerlinBlock)
	}
	if isForkIncompatible(c.LondonBlock, newcfg.LondonBlock, head) {
		return newCompatError("London fork block", c.LondonBlock, newcfg.LondonBlock)
	}
	if isForkIncompatible(c.ArrowGlacierBlock, newcfg.ArrowGlacierBlock, head) {
		return newCompatError("Arrow Glacier fork block", c.ArrowGlacierBlock, newcfg.ArrowGlacierBlock)
	}
	if isForkIncompatible(c.MergeForkBlock, newcfg.MergeForkBlock, head) {
		return newCompatError("Merge Start fork block", c.MergeForkBlock, newcfg.MergeForkBlock)
	}
	return nil
}

func (c *ChainConfig) Copy() *ChainConfig {
	cfg := &ChainConfig{
		DAOForkSupport: c.DAOForkSupport,
		EIP150Hash:     c.EIP150Hash,
	}
	if c.Ethash != nil {
		cfg.Ethash = &(*c.Ethash)
	}
	if c.AutonityContractConfig != nil {
		cfg.AutonityContractConfig = &(*c.AutonityContractConfig)
	}
	if c.OracleContractConfig != nil {
		cfg.OracleContractConfig = c.OracleContractConfig
	}
	if c.ChainID != nil {
		cfg.ChainID = big.NewInt(0).Set(c.ChainID)
	}
	if c.HomesteadBlock != nil {
		cfg.HomesteadBlock = big.NewInt(0).Set(c.HomesteadBlock)
	}
	if c.DAOForkBlock != nil {
		cfg.DAOForkBlock = big.NewInt(0).Set(c.DAOForkBlock)
	}
	if c.EIP150Block != nil {
		cfg.EIP150Block = big.NewInt(0).Set(c.EIP150Block)
	}
	if c.EIP155Block != nil {
		cfg.EIP155Block = big.NewInt(0).Set(c.EIP155Block)
	}
	if c.EIP158Block != nil {
		cfg.EIP158Block = big.NewInt(0).Set(c.EIP158Block)
	}
	if c.ByzantiumBlock != nil {
		cfg.ByzantiumBlock = big.NewInt(0).Set(c.ByzantiumBlock)
	}
	if c.ConstantinopleBlock != nil {
		cfg.ConstantinopleBlock = big.NewInt(0).Set(c.ConstantinopleBlock)
	}
	if c.PetersburgBlock != nil {
		cfg.PetersburgBlock = big.NewInt(0).Set(c.PetersburgBlock)
	}
	return cfg
}

// isForkIncompatible returns true if a fork scheduled at s1 cannot be rescheduled to
// block s2 because head is already past the fork.
func isForkIncompatible(s1, s2, head *big.Int) bool {
	return (isForked(s1, head) || isForked(s2, head)) && !configNumEqual(s1, s2)
}

// isForked returns whether a fork scheduled at block s is active at the given head block.
func isForked(s, head *big.Int) bool {
	if s == nil || head == nil {
		return false
	}
	return s.Cmp(head) <= 0
}

func configNumEqual(x, y *big.Int) bool {
	if x == nil {
		return y == nil
	}
	if y == nil {
		return x == nil
	}
	return x.Cmp(y) == 0
}

// ConfigCompatError is raised if the locally-stored blockchain is initialised with a
// ChainConfig that would alter the past.
type ConfigCompatError struct {
	What string
	// block numbers of the stored and new configurations
	StoredConfig, NewConfig *big.Int
	// the block number to which the local chain must be rewound to correct the error
	RewindTo uint64
}

func newCompatError(what string, storedblock, newblock *big.Int) *ConfigCompatError {
	var rew *big.Int
	switch {
	case storedblock == nil:
		rew = newblock
	case newblock == nil || storedblock.Cmp(newblock) < 0:
		rew = storedblock
	default:
		rew = newblock
	}
	err := &ConfigCompatError{what, storedblock, newblock, 0}
	if rew != nil && rew.Sign() > 0 {
		err.RewindTo = rew.Uint64() - 1
	}
	return err
}

func (err *ConfigCompatError) Error() string {
	return fmt.Sprintf("mismatching %s in database (have %d, want %d, rewindto %d)", err.What, err.StoredConfig, err.NewConfig, err.RewindTo)
}

// Rules wraps ChainConfig and is merely syntactic sugar or can be used for functions
// that do not have or require information about the block.
//
// Rules is a one time interface meaning that it shouldn't be used in between transition
// phases.
type Rules struct {
	ChainID                                                 *big.Int
	IsHomestead, IsEIP150, IsEIP155, IsEIP158               bool
	IsByzantium, IsConstantinople, IsPetersburg, IsIstanbul bool
	IsBerlin, IsLondon                                      bool
	IsMerge                                                 bool
}

// Rules ensures c's ChainID is not nil.
func (c *ChainConfig) Rules(num *big.Int, isMerge bool) Rules {
	chainID := c.ChainID
	if chainID == nil {
		chainID = new(big.Int)
	}
	return Rules{
		ChainID:          new(big.Int).Set(chainID),
		IsHomestead:      c.IsHomestead(num),
		IsEIP150:         c.IsEIP150(num),
		IsEIP155:         c.IsEIP155(num),
		IsEIP158:         c.IsEIP158(num),
		IsByzantium:      c.IsByzantium(num),
		IsConstantinople: c.IsConstantinople(num),
		IsPetersburg:     c.IsPetersburg(num),
		IsIstanbul:       c.IsIstanbul(num),
		IsBerlin:         c.IsBerlin(num),
		IsLondon:         c.IsLondon(num),
		IsMerge:          isMerge,
	}
}
