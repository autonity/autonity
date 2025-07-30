package params

import (
	"encoding/json"
	"log"
	"math/big"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/common/math"
)

// Ideally should be in the GO format directly
func BakerlooSchedules() []Schedule {
	const bakerlooSchedulesJSON = `
[
    {
      "startTime": 1759968000,
      "totalDuration": 63072000,
      "amount": 1000000000000000000000000,
      "vaultAddress": "0x0EB1a6E69ce1b9Cb8540863EC32e483d819704D2"
    },
    {
      "startTime": 1759968000,
      "totalDuration": 63072000,
      "amount": 1000000000000000000000000,
      "vaultAddress": "0x0EB1a6E69ce1b9Cb8540863EC32e483d819704D2"
    },
    {
      "startTime": 1759968000,
      "totalDuration": 63072000,
      "amount": 1000000000000000000000000,
      "vaultAddress": "0x0EB1a6E69ce1b9Cb8540863EC32e483d819704D2"
    },
    {
      "startTime": 1759968000,
      "totalDuration": 63072000,
      "amount": 250000000000000000000000,
      "vaultAddress": "0x0EB1a6E69ce1b9Cb8540863EC32e483d819704D2"
    },
    {
      "startTime": 1766707200,
      "totalDuration": 63072000,
      "amount": 1000000000000000000000000,
      "vaultAddress": "0x0EB1a6E69ce1b9Cb8540863EC32e483d819704D2"
    },
    {
      "startTime": 1766707200,
      "totalDuration": 63072000,
      "amount": 1000000000000000000000000,
      "vaultAddress": "0x0EB1a6E69ce1b9Cb8540863EC32e483d819704D2"
    },
    {
      "startTime": 1766707200,
      "totalDuration": 63072000,
      "amount": 1000000000000000000000000,
      "vaultAddress": "0x0EB1a6E69ce1b9Cb8540863EC32e483d819704D2"
    },
    {
      "startTime": 1766707200,
      "totalDuration": 63072000,
      "amount": 1000000000000000000000000,
      "vaultAddress": "0x0EB1a6E69ce1b9Cb8540863EC32e483d819704D2"
    },
    {
      "startTime": 1766707200,
      "totalDuration": 63072000,
      "amount": 1000000000000000000000000,
      "vaultAddress": "0x0EB1a6E69ce1b9Cb8540863EC32e483d819704D2"
    },
    {
      "startTime": 1766707200,
      "totalDuration": 63072000,
      "amount": 250000000000000000000000,
      "vaultAddress": "0x0EB1a6E69ce1b9Cb8540863EC32e483d819704D2"
    },
    {
      "startTime": 1755561600,
      "totalDuration": 7884000,
      "amount": 1000000000000000000000000,
      "vaultAddress": "0x0EB1a6E69ce1b9Cb8540863EC32e483d819704D2"
    },
    {
      "startTime": 1755561600,
      "totalDuration": 7884000,
      "amount": 500000000000000000000000,
      "vaultAddress": "0x0EB1a6E69ce1b9Cb8540863EC32e483d819704D2"
    },
    {
      "startTime": 1755561600,
      "totalDuration": 7884000,
      "amount": 250000000000000000000000,
      "vaultAddress": "0x0EB1a6E69ce1b9Cb8540863EC32e483d819704D2"
    },
    {
      "startTime": 1764028800,
      "totalDuration": 63072000,
      "amount": 1000000000000000000000000,
      "vaultAddress": "0x0EB1a6E69ce1b9Cb8540863EC32e483d819704D2"
    },
    {
      "startTime": 1764028800,
      "totalDuration": 63072000,
      "amount": 1000000000000000000000000,
      "vaultAddress": "0x0EB1a6E69ce1b9Cb8540863EC32e483d819704D2"
    },
    {
      "startTime": 1764028800,
      "totalDuration": 63072000,
      "amount": 1000000000000000000000000,
      "vaultAddress": "0x0EB1a6E69ce1b9Cb8540863EC32e483d819704D2"
    },
    {
      "startTime": 1764028800,
      "totalDuration": 63072000,
      "amount": 500000000000000000000000,
      "vaultAddress": "0x0EB1a6E69ce1b9Cb8540863EC32e483d819704D2"
    },
    {
      "startTime": 1764028800,
      "totalDuration": 63072000,
      "amount": 250000000000000000000000,
      "vaultAddress": "0x0EB1a6E69ce1b9Cb8540863EC32e483d819704D2"
    },
    {
      "startTime": 1764028800,
      "totalDuration": 63072000,
      "amount": 1000000000000000000000000,
      "vaultAddress": "0x0EB1a6E69ce1b9Cb8540863EC32e483d819704D2"
    },
    {
      "startTime": 1764028800,
      "totalDuration": 63072000,
      "amount": 1000000000000000000000000,
      "vaultAddress": "0x0EB1a6E69ce1b9Cb8540863EC32e483d819704D2"
    },
    {
      "startTime": 1754485200,
      "totalDuration": 60444000,
      "amount": 12985416666666680000000000,
      "vaultAddress": "0x0EB1a6E69ce1b9Cb8540863EC32e483d819704D2"
    },
    {
      "startTime": 1754485200,
      "totalDuration": 94608000,
      "amount": 2420000000000000000000000,
      "vaultAddress": "0x0EB1a6E69ce1b9Cb8540863EC32e483d819704D2"
    }
  ]
`
	var schedules []Schedule
	err := json.Unmarshal([]byte(bakerlooSchedulesJSON), &schedules)
	if err != nil {
		log.Fatal(err)
	}
	return schedules
}

func hexToAddressPtr(s string) *common.Address {
	addr := common.HexToAddress(s)
	return &addr
}

var BakerlooValidators = []*Validator{
	// 1dNXef
	{
		NodeAddress:   hexToAddressPtr("0xf1CA9c18C5B66a8fe56a73bBA142951CB7030137"),
		Treasury:      common.HexToAddress("0x9Dd042B73CfA67499A29D5e0eF1BeE9a20bf6d85"),
		OracleAddress: common.HexToAddress("0x5A3E1d640472F6D8642983e462C13638e5850499"),
		ConsensusKey:  common.Hex2Bytes("8d251ff9173337812bc082bd689b8329c94ee3d9b409cd2a60a332a270ee24bbf5449554d31e50314a218afe96513d91"),
		Enode:         "enode://d53b4fcebcff744d4cf039cbf23f49d175c60e4f678daa1edeb54d7510b154454630e7c21560ef4b66cdf535c5fc492cda507cfdaa430cc6d3e01560222fe042@160.250.106.47:30303",
		BondedStake:   Ntn1,
	},
	// 5AYduf
	{
		NodeAddress:   hexToAddressPtr("0xc8BF6d9A810b19Bd5d9fB0d4F74B0F9A3b5AAFc3"),
		Treasury:      common.HexToAddress("0xF1bB2F9Df6B7a25FFe4cEbec5C4708da67987772"),
		OracleAddress: common.HexToAddress("0x354C4750100C567516d6474Cb1303F04b0C695c4"),
		ConsensusKey:  common.Hex2Bytes("87c01b29913ac2908bdfb692d759a18909f365cbe94ccd2dcb138226e5471057402c0309a68de0ef264ee9965b8aa400"),
		Enode:         "enode://1a1835fa82af5f055cd56b68ecd32fd9a0e63992defa03c10409bc2420b736f25bcaed5d3ff6ec05cb59bf0a478671f05ed7189cdb925956d645cd3ca687f57f@65.21.47.49:30303",
		BondedStake:   Ntn1,
	},
	// 7La9VU
	{
		NodeAddress:   hexToAddressPtr("0xaC6Bc78b6b12520D14dAF2d50e5f2be02c5fd836"),
		Treasury:      common.HexToAddress("0xB95bF8d34a48c8e6D206E60D4d26A8B115d263D3"),
		OracleAddress: common.HexToAddress("0xfe47Da20EbD3576C4cE05f5D7e877c816eE773E7"),
		ConsensusKey:  common.Hex2Bytes("af2bb9b721b54baadd799d04be214dfded5d50ba5e39048605d891af76d0e7ca425c3ece95eafa35a8f41e7f26065c59"),
		Enode:         "enode://0c62c810933e15bf951c4cba4074dc99bc46163aee72bb35c185270c95fbd15f81af84c889af6da4e5ce55ee39038a2a6fa20158f60d00734298384eba99d605@148.251.80.29:30303",
		BondedStake:   Ntn1,
	},
	// BjcS9s
	{
		NodeAddress:   hexToAddressPtr("0x580058bB92250BF6941f8BAb221EFD7567ee7095"),
		Treasury:      common.HexToAddress("0x92C7537990EE7C30d814e3ACE186F3a846BdAbce"),
		OracleAddress: common.HexToAddress("0x195A01D501F25F159b813b83CB03a093630Cf2C1"),
		ConsensusKey:  common.Hex2Bytes("951008b065a8bd61621902b48c71471545c1551becc209fe677bf511013cd8d89e020df41ba9b00a47341a3546f62985"),
		Enode:         "enode://e44b5c1de07553cc67a0f82939ebdaf48b2b1af3c8a02d801f7241afeccd86a8ac1e9db556ad860366e22f94ce0bebf792f2fe7356d0c63f3d19502915e49c5f@65.108.206.118:30303",
		BondedStake:   Ntn1,
	},
	// buSlFR
	{
		NodeAddress:   hexToAddressPtr("0x2b465282F726529c8064E7886ceE5b25ca9F06f2"),
		Treasury:      common.HexToAddress("0xff2d345444E0c15d220f6F82EBCaD06c6BbaA0a9"),
		OracleAddress: common.HexToAddress("0x593e9c05780a32D70fe4F44Dfe8b2596C967e5C3"),
		ConsensusKey:  common.Hex2Bytes("b7653eaf4c1c89d5b8a831d13963838b29bf6d3472d8eb866c2fa42ecc9e3c495fdbbd454409a0607dd35b701a94db2e"),
		Enode:         "enode://7fcc7b3942799e15eecffdc76f4c874ef14554fa8b5fc9d2f905b402e2253fdf98af98113d27dc7828a6130b6aa40f9e962585fe8bede25cc4cd10d8803defd2@173.231.16.18:30303",
		BondedStake:   Ntn1,
	},
	// C5h4Wf
	{
		NodeAddress:   hexToAddressPtr("0x220311995c18c8aB731EFC6d1Df3658Cef26A4FF"),
		Treasury:      common.HexToAddress("0x2892461914521e5b606aAA95e2AFAFaad620631F"),
		OracleAddress: common.HexToAddress("0x2892461914521e5b606aAA95e2AFAFaad620631F"),
		ConsensusKey:  common.Hex2Bytes("a1cd25974563a2b5c6a197a69f40b4c9f49c83a306aee0f1045eb89de377b856a1e62887d8f3ed0cd0a217280af3d14f"),
		Enode:         "enode://e2e5af88ae1a6e3268c7161c97a4d5760a1cbd196cc74a173a42f5cbb5b1e586c896b49643cece12e7b5bad24745c77641bf3e4a98e9cab8f67ef24118281949@213.239.198.181:30304",
		BondedStake:   Ntn1,
	},
	// D0s608
	{
		NodeAddress:   hexToAddressPtr("0x9064dcAdF37f1eafbb6102fAbEA7F8cfd55f917f"),
		Treasury:      common.HexToAddress("0x169d652BBc9C6D913EbEB802e4Ec8A69516C3A7a"),
		OracleAddress: common.HexToAddress("0xa9E86E50Dc0B3aD31130205a67c2239CA3b1c83E"),
		ConsensusKey:  common.Hex2Bytes("8af44adf34efd8d092fae6b3abd37316cd899b0a5da41acc714fd6545af69fc072d9df2cd25290c2178e41f2b96f8443"),
		Enode:         "enode://6ae46570013b2240acd76205b1568390c9da468ac7fcd22d39e066ffe4bf07eb7d4c8379bba1e810e453515c3ccc30d618016e3c20eca3fb982644fb89fae4b8@135.181.180.144:30303",
		BondedStake:   Ntn1,
	},
	// DK29qE
	{
		NodeAddress:   hexToAddressPtr("0x9456155F24bD91E94db0Cb2846dD1D4C25c9E9e6"),
		Treasury:      common.HexToAddress("0x3e08FEc6ABaf669BD8Da54abEe30b2B8B5024013"),
		OracleAddress: common.HexToAddress("0x4D8387E38F42084aa24CE7DA137222786fF23A3E"),
		ConsensusKey:  common.Hex2Bytes("b705f3e14132148c02ce92eed0700945b83150be4bbecf14342f608d2bcf9179434089cb9570a974952f941f77bc17bd"),
		Enode:         "enode://11a803ed8573eab91a174ac52b7047b26721876020bbb970c8ab8abc4799d0f7147e11305d73763def6fe74a69908b20d1716e9d2d707bf6d7afa38bfc427804@34.142.43.147:30303",
		BondedStake:   Ntn1,
	},
	// eO6oqk
	{
		NodeAddress:   hexToAddressPtr("0x893821B7bf645eD95B2c574ECD22145091AacC0a"),
		Treasury:      common.HexToAddress("0x74Ad02697fF0C88A752F11c3B0CE82063bD27132"),
		OracleAddress: common.HexToAddress("0x5Ab50DDe7587804EA7675e163F1E5178a01Ac30a"),
		ConsensusKey:  common.Hex2Bytes("8dce049480a8e4e1abbd341fea4be35231f9ff3c3d4343282e69320fbe3ce9901dfbe793c4d908f4d68216b06fe02497"),
		Enode:         "enode://0bac62f7f692ea9ea6bb4d763715679f36852eb8fdd1cbf21b52b69e8cf30503c19181273426da183d0da93e93076c3b93065518c94f1df14cd19c847bc036c0@210.109.80.48:30303",
		BondedStake:   Ntn1,
	},
	// GxHedr
	{
		NodeAddress:   hexToAddressPtr("0xbBfE1EB7747fBea8Bf71970D85d8406a62c244A7"),
		Treasury:      common.HexToAddress("0x493ca6b6574D8B6d4A8609142eC590cEb1117DaB"),
		OracleAddress: common.HexToAddress("0x709b553e93BAfA1fF4eD3BD469eEd81924488003"),
		ConsensusKey:  common.Hex2Bytes("88df820b856be9fc445e1f06236fbfe6b390164c7c283951c62b85d5efeb679689013260c6d595812b4e51a7ebfbb420"),
		Enode:         "enode://be35f549821d72dbf384363d9b7779dd06f58920a901810501e41addad393cc69a59b19d4e32691941e9ea6d389ace2e85913d93cbdc3321f663f1df812aa548@45.153.35.148:30200",
		BondedStake:   Ntn1,
	},
	// HrLoXr
	{
		NodeAddress:   hexToAddressPtr("0x567DECAef0D8F3Cc599f7021Ad7f53bb9c4d214F"),
		Treasury:      common.HexToAddress("0xba8A40A19fd948348DfB617c80b5eC94d8AEC704"),
		OracleAddress: common.HexToAddress("0xe2F52f5345B319f543B8b610495EfB3C00Ec96E3"),
		ConsensusKey:  common.Hex2Bytes("af0e98065fe51de416617b08a22be4be6c13da5ffcd5893ae882d89b46a485bed41bcaa5595657ba7de364c05891a451"),
		Enode:         "enode://1f17288fe5c004d943d7dcf0e820f7b7bb12f84a67eb075f058e609c7325de5847b0ef4f1b58ef1dde2593a67512d3722cd5fa2620288bff43dadee988f34d04@65.109.154.189:40303",
		BondedStake:   Ntn1,
	},
	// Ilg5ne
	{
		NodeAddress:   hexToAddressPtr("0x0E845E24b33d5eb6FbB7d779408aeD970cf083BD"),
		Treasury:      common.HexToAddress("0xC65A4E9ce7067c6DEF75BD4344CC59444Ea9bcE6"),
		OracleAddress: common.HexToAddress("0x179Db11f32b0CBb326a0018FA0D12335Eaa18380"),
		ConsensusKey:  common.Hex2Bytes("ad4889bba251b82d8345e50ebac4805bd167770da1287ce5b31b3e048a8fe1c4acd580f4fb58b3619acd92235d39f649"),
		Enode:         "enode://f3b1f57246c2f8b68b41b33de78214944774f249f8571b8cc672c1c7efde389b9700b3fb90c18ba1a671702a000b89a484d0fb28665ef72fb2b7cb68fde1683d@149.50.110.220:30303",
		BondedStake:   Ntn1,
	},
	// J4mPq8
	{
		NodeAddress:   hexToAddressPtr("0xA1D0b6c927EB0C4aD6c473A127AA7516283C80d1"),
		Treasury:      common.HexToAddress("0x941091853B35233e17dD47935e44b47e3A277D56"),
		OracleAddress: common.HexToAddress("0x2d05CE1Efaad6C5fe51316eFef696fBE9EE08534"),
		ConsensusKey:  common.Hex2Bytes("a489c57d92feb966d72a160950e7fe156ce055021ccbc216792c1dd83e213ac48841a5b8acc27f315105b10174c3f6ec"),
		Enode:         "enode://c9999c795955bb6949066911276298b90e3b57843293e89311e7bc0848247081db9d1e40720da89ce704ce1279d378657e95a01e0ab833b78116b685fccae335@167.235.177.151:30303",
		BondedStake:   Ntn1,
	},
	// L9bXw2
	{
		NodeAddress:   hexToAddressPtr("0xd16E87cA43814a9f31BDAF6C86f3a3B170C5ab00"),
		Treasury:      common.HexToAddress("0xC4A46a46697404DFE6bDD220ba15e93431852227"),
		OracleAddress: common.HexToAddress("0x6a301D9B8929c5F03b3Dc4bCF938C1866863b692"),
		ConsensusKey:  common.Hex2Bytes("81b92b30f13cce4ed6a55ff7b33f72390451de7e213f04697f36724e3034ca2977426b1b72c2ca8781409f7ea5bcb255"),
		Enode:         "enode://0d91c4abad25a2737e453f0c4d1568c23884be4697cf3a8426738bd6523026011746108a479a8faf7215fdef24bdd46887ccae53508c6f796fd5fffcaeab6787@37.27.127.216:33333",
		BondedStake:   Ntn1,
	},
	// L9RwkS
	{
		NodeAddress:   hexToAddressPtr("0x3e8d115bd4287cD1A2Df2238a4059a8e1348882B"),
		Treasury:      common.HexToAddress("0x3941bDd2cB0B3BdE01FC9139Af6db0F7052e161D"),
		OracleAddress: common.HexToAddress("0xd18d3dCAbCc0C7E075327cd620235333dc9A01EE"),
		ConsensusKey:  common.Hex2Bytes("83f34f59db8e04d9e8617ac10bfb2bf494813676cb16115d8369c3f6f89f1eea3ea5da1f286584d668407588d4c4ad2f"),
		Enode:         "enode://8d64d872b8510e3df4295f402f44d1c9243470c574d34473af717814cc2f687464eae61086f11a3a477ca3cb4d8ac138e82ca29827aecaf2dfea0f6ba86e2431@85.195.103.141:30303",
		BondedStake:   Ntn1,
	},
	// lvlO3g
	{
		NodeAddress:   hexToAddressPtr("0x0ec82995BcEa73b2EdB0f293b918F6938fDb15bF"),
		Treasury:      common.HexToAddress("0x7Fe1abe90074eAFE3CcB895AA63d8bE05112C49a"),
		OracleAddress: common.HexToAddress("0xfc410ed3eAb7c405D4e9D0c1cb21FFe5bF2aF853"),
		ConsensusKey:  common.Hex2Bytes("85ab94be782ef1d3f2e6ca3eaaef1abb339c269a9e162f56ef8442908c3f077b8c87008a359be5b72ee43465b72e52d8"),
		Enode:         "enode://7fdcf1f97a79ab7aa801ae7511cca42c4ca68f1a46d7c26c0930a34e9d8dd66d317ae30330fdad62dd67b2d89a5bd02bf8eb55f027a29c2077f05072d6368100@45.82.65.55:30303",
		BondedStake:   Ntn1,
	},
	// NBXpxO
	{
		NodeAddress:   hexToAddressPtr("0xB709BD9701cC5C6774b98511be4B78f234085d8C"),
		Treasury:      common.HexToAddress("0x40077ae14aD4f5D2A92b19d48Ff258be9caa8B46"),
		OracleAddress: common.HexToAddress("0xb28487510B3bf049B6E0968934A1A2853e933bE2"),
		ConsensusKey:  common.Hex2Bytes("8e34fb90980e937f879551205ecaf7ebd85435c6a4edff5814f85e612ec53e0be6d313b9f92f3fd561126e906bdd3ae8"),
		Enode:         "enode://545a77c10fd3195fbb10c601a2f7f7ed8eb8c12be2ebb23e88ad9cbc68e924d2956ea3400fa70858f072670e1994e2f261bde74cc9e85a161e445b68a3bdbfc4@49.12.84.248:30303",
		BondedStake:   Ntn1,
	},
	// pYcWqU
	{
		NodeAddress:   hexToAddressPtr("0xB7ce027e296E488a9DA6a4c8B2FD41bE068Adfc1"),
		Treasury:      common.HexToAddress("0xfD9972984905aa43a9a5cc27217544A79daAa0ED"),
		OracleAddress: common.HexToAddress("0xb26323FF2F462a49e7c74BBb982eeAA95832fD30"),
		ConsensusKey:  common.Hex2Bytes("952cc206527076ce663cb87749193108ef511430006dc39458caefd601cb748a226a7761250ca3783a64fa1501f59bfc"),
		Enode:         "enode://fc2884d08b4cba9546a695a18bd57696b847cb2f2f76288d07dbc4250393e784b8e826c8a5b900073f0245609635fbc70de39df0cc63adb36b2b617f8965e88c@37.252.186.239:30303",
		BondedStake:   Ntn1,
	},
	// q0QnDh
	{
		NodeAddress:   hexToAddressPtr("0x36b825D611930EdA4E3eA61243eC5dc1fBD40b5c"),
		Treasury:      common.HexToAddress("0xDc9a47d6Ce6A35a8Cd2AF94C710d6D0a754594ca"),
		OracleAddress: common.HexToAddress("0x062CEf0f81DD65BC88C3C46F27BC9ED439c6D37C"),
		ConsensusKey:  common.Hex2Bytes("888bd408be35657d7377d9c009c75f13441d42d98f6f1fa022ff026f839f72ffb5c9e91714cbade5dadbc95863669c95"),
		Enode:         "enode://3f561a9041fa5b351d3d6258285f794716518e500e31ef9099f27583845eea464d8faf36f74fdea3ee6fef12acb3415dfd22c7b6aab87f37ad96818508b55f5c@207.121.21.148:30303",
		BondedStake:   Ntn1,
	},
	// q7Z1nS
	{
		NodeAddress:   hexToAddressPtr("0x8787393A67db15c10b0c6689FC717057b967c0b0"),
		Treasury:      common.HexToAddress("0xf8f97ea875bB3348cc1A55E56Bbe92ebC9cD682D"),
		OracleAddress: common.HexToAddress("0x0751Da7041f99F8310B13c56E69aF5c98589F9A5"),
		ConsensusKey:  common.Hex2Bytes("a082362f0299e521d4aaa457a9b1c47183dceb29424d7d871b0c940c772a33cf2736e430f6b2e802bbf38fc8efc86ab8"),
		Enode:         "enode://0ab5b19d89bcc9571cb8868bd79d75badf6c81246c59093e3f5e51e5b9f1ba28ed72c63796909425e8def94d90263b8f4bb2ff68876c71ca14eacede5e760074@65.108.106.168:30303",
		BondedStake:   Ntn1,
	},
	// Qlb0W3
	{
		NodeAddress:   hexToAddressPtr("0x4237080d57e2bB59b94F06ae49557A8277b2374e"),
		Treasury:      common.HexToAddress("0xB8a67f8d7C3Ad3d6a99dAE7F28eFf0A640A30512"),
		OracleAddress: common.HexToAddress("0xB1300c6f35c85149683E756a2D6Ae6f79DB34E95"),
		ConsensusKey:  common.Hex2Bytes("8c8ed014b61090b3644f8632642bd1c7129dd38140438f63fb9782205aff4a1961f82cfeb8364026f1ce8c7b916919a8"),
		Enode:         "enode://dedb1626a0276d9c31b76b3487fe35d2a7b8fb278c4687c2dd7f3daa7de4786b5357a611534783b7fe0b3343f896b3c4e90556810824397d04a5857706c5a55f@185.100.10.164:30310",
		BondedStake:   Ntn1,
	},
	// qOe4W1
	{
		NodeAddress:   hexToAddressPtr("0xD6F662Ef7C486d49E2F7356a69B1c4664C2Ca640"),
		Treasury:      common.HexToAddress("0x531bbA0AaAa7AA9C5811865192CF2639bC96C2D0"),
		OracleAddress: common.HexToAddress("0x852Ef92b11dC8eb4B5fd7968e14A71f1de13CFC5"),
		ConsensusKey:  common.Hex2Bytes("a9ecabb98635fd1cd25e485ece6dbf2833939635ee1d981d74b310fa30e6badede8cb8773d04d9f2c3f62facda9ba2b1"),
		Enode:         "enode://b7a3b981bc1c1f636519cc13c28b5930329db28d810e93f79735e38f5d5e8283470c98aa4478440e223d776426325d5d37f26ef8ec1657b37272e0d39870806f@35.217.2.146:30303",
		BondedStake:   Ntn1,
	},
	// Tgcs7W
	{
		NodeAddress:   hexToAddressPtr("0x31B1f4C0875F4940e88f9F4fe3247435cE334165"),
		Treasury:      common.HexToAddress("0x578b27feCbb7741EcdE6Ae8197F37061122dD880"),
		OracleAddress: common.HexToAddress("0x31342aB7A751292468532c97CC62B35e57df1695"),
		ConsensusKey:  common.Hex2Bytes("ae2e81112dc72ff303b29a97e2d8b3ef622e319851014a43c809b91c5c75a63d41e0c0bd33444609a535c56aa73a6b15"),
		Enode:         "enode://0a2a825c8f6bfe0475a9d77f1108755b6828a817bd1dc1c28c70b57793d878deea18f543082bbdb35cfc1ce61cf423755c366011b1fa6de1d348794e07ab68b6@178.205.102.224:30303",
		BondedStake:   Ntn1,
	},
	// TyXZGk
	{
		NodeAddress:   hexToAddressPtr("0x5909673b7690Ef6Bf61Df6D26A11efe5Afbf0D27"),
		Treasury:      common.HexToAddress("0x1103AFfeb7AACB5bc45991cC6BCf6023AD3769ed"),
		OracleAddress: common.HexToAddress("0xB315Fe43994a6D3B92fAA67E9061De12CB7E75fD"),
		ConsensusKey:  common.Hex2Bytes("a4c0ec12e19f35cae2c8fb75e7471fbd280311c20e9a2ce80578908248f8258facdd50e8c22a9b150c2386aebc1ae4cd"),
		Enode:         "enode://0f1d6fef2eb3542e15d8b7c5c6d385dd62537b6d8aae4fc79bb99d35dd1e5a14738984fa7e9b3bf85407a5381a53adda05dd70dcf2097f0101d411978cc03f1c@173.231.41.98:14403",
		BondedStake:   Ntn1,
	},
	// ULaZrL
	{
		NodeAddress:   hexToAddressPtr("0xf63e6a8F976FA5e0815BEF73bce4973f40555A42"),
		Treasury:      common.HexToAddress("0x9949c71cd6B8CBfd8B9B9Ab4bc8b58A16C39fE88"),
		OracleAddress: common.HexToAddress("0x1d5c10cDF36cb2333Bb53BA9E1CeDb07023184C8"),
		ConsensusKey:  common.Hex2Bytes("86be742a455f4d580a0eacfdb802ab212a33c3c3f5d2286a3368f4dde638fc2a5538f693ddf0857b817b302b1ba06517"),
		Enode:         "enode://641c5ca840a654febe1f8ecff797554aa2f2148205a6194e5dd49ea2d95b0b30df0c9a43815334055f320dcc52cd160a38519f86534f63b136899420491c2780@37.27.123.37:30303",
		BondedStake:   Ntn1,
	},
	// vd73Er
	{
		NodeAddress:   hexToAddressPtr("0xd97aEb638B6aa915a0F2bebb1b8BB80596c6e00C"),
		Treasury:      common.HexToAddress("0x3D3A1b7AD09C87d40EA333D26FF22cBED367aD98"),
		OracleAddress: common.HexToAddress("0x0f92ff67948A871d6E9732FBAceC25eCE55E8d2b"),
		ConsensusKey:  common.Hex2Bytes("974eebf0d6d2caea87451a9043310163b82c8eba6f7890323074bb8a701a34cada1efbc057a12a58287a91a93a7692fd"),
		Enode:         "enode://be20daaedcd1a1bbbc2ac033036762f49cc414981227d14b3bfbd86a6a3c6d43a560ab7193c856495c655c0f253308bc9fb4063920f2abdf72caebedd852cccd@54.195.150.136:30136",
		BondedStake:   Ntn1,
	},
	// WHgEhi
	{
		NodeAddress:   hexToAddressPtr("0x515DD5Dc41Cf3Ad2BFdfb03f0F955D22f696aC06"),
		Treasury:      common.HexToAddress("0x7d1373A0b7bEb6615999c2D316e6566732a6a1fc"),
		OracleAddress: common.HexToAddress("0x3525D3E9f270e70a7CFfdBD2CCb8B748F59F129A"),
		ConsensusKey:  common.Hex2Bytes("b6e01a752853a6ab324e318da4da68c5563186c18627b6d20c8e22ec74908e2a557ad83f4b55c99018c027de854d1790"),
		Enode:         "enode://a15fa248277caaa7a4e17435cf6be12255fb76fa82ec2e70d85acebecbd6ea4733f995243b848c71a898a685c563f6dedff5848afa9283635b2dfb2461d5f23a@51.89.40.88:30303",
		BondedStake:   Ntn1,
	},
}

// BakerlooChainConfig contains the chain parameters to run a node on the Bakerloo test network.
var BakerlooChainConfig = &ChainConfig{
	ChainID:                 big.NewInt(65_010_004),
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
		MinBaseFee:               10_000_000_000,
		EpochPeriod:              DefaultEpochPeriod,
		GasLimit:                 30_000_000,
		GasLimitBoundDivisor:     DefaultGasLimitBoundDivisor,
		BaseFeeChangeDenominator: DefaultBaseFeeChangeDenominator,
		ElasticityMultiplier:     DefaultElasticityMultiplier,
		UnbondingPeriod:          6 * 60 * 60,
		BlockPeriod:              1,
		MaxCommitteeSize:         27,
		Operator:                 common.HexToAddress("0x83e5e0eab996Bb894814fa8F0AC96a0D314f06F3"),
		Treasury:                 common.HexToAddress("0xd735174cf1d0D9150cb57750C45B6e8095160f6A"),
		WithheldRewardsPool:      common.HexToAddress("0xd735174cf1d0D9150cb57750C45B6e8095160f6A"),
		TreasuryFee:              50_000_000_000_000_000, // 5%
		InitialInflationReserve:  (*math.HexOrDecimal256)(new(big.Int).Mul(big.NewInt(40_000_000), NtnPrecision)),
		DelegationRate:           1000,
		WithholdingThreshold:     0,   // 0%, no tolerance
		ProposerRewardRate:       500, // 5%
		OracleRewardRate:         500, // 5%
		SkipGenesisVerification:  false,
		// Schedules are set in the DefaultBakerlooGenesisBlock
		MaxScheduleDuration: uint64(4*SecondsInYear + SecondsInDay),
		ClusteringThreshold: 100,
		Validators:          BakerlooValidators,
		TokenBond:           (*math.HexOrDecimal256)(new(big.Int).Mul(big.NewInt(1_620_027), NtnPrecision)),
		TokenMint:           (*math.HexOrDecimal256)(new(big.Int).Mul(big.NewInt(60_000_000), NtnPrecision)),
	},
	OracleContractConfig: &OracleContractGenesis{
		Symbols:                   []string{"AUD-USD", "CAD-USD", "EUR-USD", "GBP-USD", "JPY-USD", "SEK-USD"},
		VotePeriod:                600,
		OutlierDetectionThreshold: 3,
		OutlierSlashingThreshold:  225,
		BaseSlashingRate:          10,
		NonRevealThreshold:        5,
		RevealResetInterval:       10,
		SlashingRateCap:           50,
	},
	ASM: AsmConfig{
		ACUContractConfig: &AcuContractGenesis{
			Symbols:    []string{"AUD-USD", "CAD-USD", "EUR-USD", "GBP-USD", "JPY-USD", "SEK-USD", "USD-USD"},
			Quantities: []uint64{1_826_272, 1_634_362, 1_023_622, 886_976, 177_065_877, 11_513_754, 1_192_007},
			Scale:      uint64(7),
		},
		StabilizationContractConfig: &StabilizationContractGenesis{
			BorrowInterestRate:        (*math.HexOrDecimal256)(math.MustParseBig256("50_000_000_000_000_000")),
			AnnouncementWindow:        (*math.HexOrDecimal256)(math.MustParseBig256("3600")), // 1 hour
			LiquidationRatio:          (*math.HexOrDecimal256)(math.MustParseBig256("1_800_000_000_000_000_000")),
			MinCollateralizationRatio: (*math.HexOrDecimal256)(math.MustParseBig256("2_000_000_000_000_000_000")),
			MinDebtRequirement:        (*math.HexOrDecimal256)(math.MustParseBig256("1_000_000")),
			TargetPrice:               (*math.HexOrDecimal256)(math.MustParseBig256("1_618_034_000_000_000_000")),
			DefaultNTNATNPrice:        (*math.HexOrDecimal256)(math.MustParseBig256("666_619_300_000_000_000")),
			DefaultNTNUSDPrice:        (*math.HexOrDecimal256)(math.MustParseBig256("900_000_000_000_000_000")),
			DefaultACUUSDPrice:        (*math.HexOrDecimal256)(math.MustParseBig256("834_405_200_000_000_000")),
		},
		SupplyControlConfig: &SupplyControlGenesis{
			InitialAllocation: (*math.HexOrDecimal256)(
				new(big.Int).Sub(
					new(big.Int).Sub(
						new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil), // 2^256
						new(big.Int).Mul(big.NewInt(403), big.NewInt(Ether)),
					),
					common.Big1,
				),
			),
		},
		AuctioneerContractConfig: DefaultAuctioneerGenesis,
	},
	AccountabilityConfig: DefaultAccountabilityConfig,
	OmissionAccountabilityConfig: &OmissionAccountabilityGenesis{
		InactivityThreshold:    1500,   // 15%
		LookbackWindow:         60,     // 60 blocks
		PastPerformanceWeight:  1000,   // 10%
		InitialJailingPeriod:   10_000, // 10000 blocks
		InitialProbationPeriod: 8,      // 8 epochs
		InitialSlashingRate:    5,      // 0.05%
		Delta:                  5,      // 5 blocks
	},
	InflationContractConfig: &InflationControllerGenesis{
		InflationRateInitial:      (*math.HexOrDecimal256)(new(big.Int).Div(new(big.Int).Mul(big.NewInt(75), DecimalFactor), big.NewInt(1000*SecondsInYear))),        // 7.5% AR
		InflationRateTransition:   (*math.HexOrDecimal256)(new(big.Int).Div(new(big.Int).Mul(big.NewInt(55), DecimalFactor), big.NewInt(1000*SecondsInYear))),        // 5.5% AR
		InflationReserveDecayRate: (*math.HexOrDecimal256)(new(big.Int).Div(new(big.Int).Mul(big.NewInt(17_317), DecimalFactor), big.NewInt(100_000*SecondsInYear))), // 17.3168186793% AR
		InflationTransitionPeriod: (*math.HexOrDecimal256)(new(big.Int).Mul(big.NewInt(4*SecondsInYear+SecondsInDay), DecimalFactor)),                                // (4+1/365) years
		InflationCurveConvexity:   (*math.HexOrDecimal256)(new(big.Int).Div(new(big.Int).Mul(big.NewInt(2_765), DecimalFactor), big.NewInt(1_000))),                  // 2.765

	},
}
