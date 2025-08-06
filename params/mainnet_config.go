package params

import (
	"encoding/json"
	"log"
	"math/big"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/common/math"
)

// Ideally should be in the GO format directly
func MainnetSchedules() []Schedule {
	const schedulesJSON = `
[
    {
      "startTime": 1760486400,
      "totalDuration": 63072000,
      "amount": 1000000000000000000000000,
      "vaultAddress": "0x6D78AF53357626a2ccB0B555e0DA6846d1489659"
    },
    {
      "startTime": 1760486400,
      "totalDuration": 63072000,
      "amount": 1000000000000000000000000,
      "vaultAddress": "0x59E1F17410E854BF1b3B21639C254d34AAD2148f"
    },
    {
      "startTime": 1760486400,
      "totalDuration": 63072000,
      "amount": 1000000000000000000000000,
      "vaultAddress": "0x0d4F741f28f92A048B904E87A885C26FD0DCb146"
    },
    {
      "startTime": 1760486400,
      "totalDuration": 63072000,
      "amount": 250000000000000000000000,
      "vaultAddress": "0x358B57b61CacB2B94E51d71683Aabc6cdee85c17"
    },
    {
      "startTime": 1767225600,
      "totalDuration": 63072000,
      "amount": 1000000000000000000000000,
      "vaultAddress": "0x2Cc4b2ECF1A7f4dDC85d4067dd85a9f4FeAcb2a0"
    },
    {
      "startTime": 1767225600,
      "totalDuration": 63072000,
      "amount": 1000000000000000000000000,
      "vaultAddress": "0x500d93485FEBCB4B9217E54fe781523E0f7C5cbD"
    },
    {
      "startTime": 1767225600,
      "totalDuration": 63072000,
      "amount": 1000000000000000000000000,
      "vaultAddress": "0x118205bdA68e16839B25ce6aa9E050238d2e581e"
    },
    {
      "startTime": 1767225600,
      "totalDuration": 63072000,
      "amount": 1000000000000000000000000,
      "vaultAddress": "0xF43dcE089810C4284D7FaE75f9181CeC5E2888d9"
    },
    {
      "startTime": 1767225600,
      "totalDuration": 63072000,
      "amount": 1000000000000000000000000,
      "vaultAddress": "0x76cE9c136277D5e0F94BAA1d5c0F4D8163887Add"
    },
    {
      "startTime": 1767225600,
      "totalDuration": 63072000,
      "amount": 250000000000000000000000,
      "vaultAddress": "0x961755887191D3dbf81c6ba4c121dd7f77828B57"
    },
    {
      "startTime": 1756080000,
      "totalDuration": 7884000,
      "amount": 1000000000000000000000000,
      "vaultAddress": "0x7774218cC31CA0A6B022cC34354a06F4c2D31399"
    },
    {
      "startTime": 1756080000,
      "totalDuration": 7884000,
      "amount": 500000000000000000000000,
      "vaultAddress": "0x19De27800D5b28D139e7B54C6FE01cA68C2B9Fb7"
    },
    {
      "startTime": 1756080000,
      "totalDuration": 7884000,
      "amount": 250000000000000000000000,
      "vaultAddress": "0x746Bd217f838D1020245a80bA8365320d3c7df71"
    },
    {
      "startTime": 1764547200,
      "totalDuration": 63072000,
      "amount": 1000000000000000000000000,
      "vaultAddress": "0xa08D9B10A24747A21D5dB0B862a98E4Ee1693e7b"
    },
    {
      "startTime": 1764547200,
      "totalDuration": 63072000,
      "amount": 1000000000000000000000000,
      "vaultAddress": "0x4F91BE0Befa72fe03Fc163654fd258d231285B91"
    },
    {
      "startTime": 1764547200,
      "totalDuration": 63072000,
      "amount": 1000000000000000000000000,
      "vaultAddress": "0x92d0C8A3F8F2fa0829D7F6318Bda06C8A103124c"
    },
    {
      "startTime": 1764547200,
      "totalDuration": 63072000,
      "amount": 500000000000000000000000,
      "vaultAddress": "0x05BE870d6683fA71dC0cD34483BdC6f79D7e6bfA"
    },
    {
      "startTime": 1764547200,
      "totalDuration": 63072000,
      "amount": 250000000000000000000000,
      "vaultAddress": "0xE5507365A62D042850d501dBbdA9BA66b1E2Ade9"
    },
    {
      "startTime": 1764547200,
      "totalDuration": 63072000,
      "amount": 1000000000000000000000000,
      "vaultAddress": "0x2789cbc4987BdD7020929cB7C0B1313026D892D7"
    },
    {
      "startTime": 1764547200,
      "totalDuration": 63072000,
      "amount": 1000000000000000000000000,
      "vaultAddress": "0xe447c8C6E756Ce8983C6C468F5B9114bDc046e5b"
    },
    {
      "startTime": 1755003600,
      "totalDuration": 60444000,
      "amount": 12985416666666680000000000,
      "vaultAddress": "0x946b4f0e4593F630Dac3F6100A154B7D417709cC"
    },
    {
      "startTime": 1755003600,
      "totalDuration": 94608000,
      "amount": 2420000000000000000000000,
      "vaultAddress": "0x946b4f0e4593F630Dac3F6100A154B7D417709cC"
    }
]
`
	var schedules []Schedule
	err := json.Unmarshal([]byte(schedulesJSON), &schedules)
	if err != nil {
		log.Fatal(err)
	}
	return schedules
}

var MainnetValidators = []*Validator{
	// 3u89kP
	{
		NodeAddress:   hexToAddressPtr("0xAd68B8994bc1433c15Dd8cCa3eC3AeFd7c5cb85a"),
		Treasury:      common.HexToAddress("0x843eBA81dAFDb7D151099e88386B8D86aa806Bef"),
		OracleAddress: common.HexToAddress("0x89551a34Dd9b9a99E1e826992E0225b81f668FA8"),
		ConsensusKey:  common.Hex2Bytes("8089d4bee0a890cc48ba5dba03be95debf011e23170a08fa4536b7f49aaff68b07b8e6aa304382403f002fadc21712be"),
		Enode:         "enode://8a0b5999e97aae4bfa34dbd9e5603db8f49ed0b56bf687d50c8baf9249e58017df02a02cf01b20fdf1c732261f4f7b9ff2433b3af7ddd8253373b777c311b6eb@173.231.40.170:14403",
		BondedStake:   Ntn1,
	},
	// 5u9RnB
	{
		NodeAddress:   hexToAddressPtr("0xF26b0Ca600eF22Ab0A7d9fe854D4095421A98dDB"),
		Treasury:      common.HexToAddress("0x723E55368213EB6b11A7AE931F7645C7d3A9323d"),
		OracleAddress: common.HexToAddress("0x2C1feB9858E6fb71457cfd86b0250d0FaD3AD256"),
		ConsensusKey:  common.Hex2Bytes("8b651f78fd958781198f9ea29601a09d8d335ac8c6a5fb9599146686d12e73239eafe07c5fc3e85276464e99ce6d3825"),
		Enode:         "enode://e8e720e95166d78eb73927fcd338b393d7a2bc16f7df0560ee2593b89c899b0fda278c857fe03cc56d1a1cc39fd8bdfe2279c231113ddc98f8010bb5c8896d7e@57.128.187.32:30303",
		BondedStake:   Ntn1,
	},
	// 6xNQr7
	{
		NodeAddress:   hexToAddressPtr("0x185332Cd8c8bCA5F902d0e329f9374937F8B3379"),
		Treasury:      common.HexToAddress("0x2728800AB45F050F8106919022c35b1E500Cbb2F"),
		OracleAddress: common.HexToAddress("0x6047B2E26Ecb5EBE6FfB3dc45C89a05c9f68C542"),
		ConsensusKey:  common.Hex2Bytes("b2d43210b3ed0b04e32667e30a1560f6f608894424dbdfc23a30489ebdbf90323c9f5eb446f8934f1b75abeef044015e"),
		Enode:         "enode://75d6eb3a35ad10c2d32575783456a0c2c98637281c3d7e37acbdc3c00829997bdbe5f9f1b919256fa1bbf11524e15e8434b9674b5e95b56a8983ec9f26bd38e2@162.55.199.235:30303",
		BondedStake:   Ntn1,
	},
	// 7kX9Qd
	{
		NodeAddress:   hexToAddressPtr("0x30A3Bf62Eb474C0A49fa8B4151EB02d6493d2f9f"),
		Treasury:      common.HexToAddress("0xb3b3EF4D2Bc7d4DA487aab157092Ed8442F67101"),
		OracleAddress: common.HexToAddress("0xab5D28Db79C124c136225763d48E289D7d711D4b"),
		ConsensusKey:  common.Hex2Bytes("8fccb7fca127f77443d1c0a84788d498c34d26b5229b137e7b53004199d2840f69d4390a8364c40ffa30680fde6b36bf"),
		Enode:         "enode://057a4338aa7e77227c9d8c8cac69d03c0cf3c131b8813b05221a3f0b531b13813cafbc06a477ac2f6304cc599155a9a51eec28d388c2fc00c7d1604c0d2a7d45@64.176.54.103:30303",
		BondedStake:   Ntn1,
	},
	// A5zXc9
	{
		NodeAddress:   hexToAddressPtr("0x0cBa879804EF58Fd77E34abBB46f391979721B25"),
		Treasury:      common.HexToAddress("0x99A898F17B34711B912830186AC066B0499E88B0"),
		OracleAddress: common.HexToAddress("0x2246Fceb09e59079f6F75F8365Ee6E6dBe4BbA70"),
		ConsensusKey:  common.Hex2Bytes("b3159fad138c73ec789aa31879e53d90a73c0b256059a8ec52c0868da145393b1e317156230692112fe56b3a209a39a6"),
		Enode:         "enode://f1c0231b93492a88f62d09168545a41d0ee62ca6592ca7260bc4f62bab4ebba32f35d89a8527604c988066f331fbb1954d07aeb85878e2d083ee136b466acb7c@207.121.63.124:30303",
		BondedStake:   Ntn1,
	},
	// aZ7K2M
	{
		NodeAddress:   hexToAddressPtr("0xCf0E242bFBA2F57051FD5Cc1d738Dd6f25f514F5"),
		Treasury:      common.HexToAddress("0x0881B2CC40745C94274052f33B3E9Aa1F53b7F19"),
		OracleAddress: common.HexToAddress("0xda133740E48bcaC7407EB753F8faD1a6EF0CB567"),
		ConsensusKey:  common.Hex2Bytes("b222c408d967b9d4b63bbf12e3f7845ef2aa16a0162636f02332ecad49e47c0b818fe139af3355206399fcdc12f4079b"),
		Enode:         "enode://f3349b12760b3225fc1342845ec475d6967dab1e830cc053ea9a0dd482bf6f011a6b986c2fd1a2f35d750f87d38ed836e10949c941ccc87b903caff8d4a4662c@160.250.106.108:30303",
		BondedStake:   Ntn1,
	},
	// b2D8Zh
	{
		NodeAddress:   hexToAddressPtr("0xA07D93291FE3dC2193F87C4142d71457cfdd81ac"),
		Treasury:      common.HexToAddress("0x20Be8aD8c175986B3f0342d9a220859F7A070aAb"),
		OracleAddress: common.HexToAddress("0x32a9213cdFf5840D3bF82282dcE69C8147C7A877"),
		ConsensusKey:  common.Hex2Bytes("a73f1f792f9ec4e6af18cc012cee0f511e51f000cd0f23c07d645e18c8b86b9b4e74ab663a5c8fc6bd0b80aad33d2af7"),
		Enode:         "enode://ff7ed4fa73bee5ed7aad746855c917fc9179e5f5845ae21bd83d9c25208fe41a7eb4d00286b895d12d14e8ccb3dc819832822f5866300f2481e1a86ef5a7ed66@35.197.235.224:30303",
		BondedStake:   Ntn1,
	},
	// B80yXe
	{
		NodeAddress:   hexToAddressPtr("0xD302AE484906AD70FC714F9348E738aaB6c87DF7"),
		Treasury:      common.HexToAddress("0x565AB60042F22b555d69A28437F1021a48D6567D"),
		OracleAddress: common.HexToAddress("0xA10bDFCD09C2EfB8F8392d40ea710EFbc7b33068"),
		ConsensusKey:  common.Hex2Bytes("b1f7c86e669925cfd6e08ab4f5468c5a18711b21310238223b46e4c6951f83a48de931f755ce134466eceae14b3e0770"),
		Enode:         "enode://709d3bf926ba0fef64476b4673396344aaecc3f91c4b7e4c238016b0f38eac05b746ae37e99c0e2fd11d52f1ceae7109c821493ab864570b633832fbc75c94ed@185.26.11.57:30303",
		BondedStake:   Ntn1,
	},
	// c9F7Dz
	{
		NodeAddress:   hexToAddressPtr("0x56FE65a264a192e55F374858604CA57191b6f71f"),
		Treasury:      common.HexToAddress("0x1610F16c7562dE770D4C47876A6169789F7Fd224"),
		OracleAddress: common.HexToAddress("0xf07CDD9f6F9fFFBd174611D996c107E7e3335638"),
		ConsensusKey:  common.Hex2Bytes("98f2cff3d042cc5ea9346caf3d38116ab104829469e29d961cfc4c8f7514cdd37f0d992aba217a4c1f3c47418d9a8882"),
		Enode:         "enode://4f0acb75485d6c93cccea9c69b5a91882b6225542066a1592530792efdb735b5af7a83ccadc7b279100ec1d04982025e0284af3998e7945ce6959ae90de51c61@74.50.95.190:30303",
		BondedStake:   Ntn1,
	},
	// d5Y9Tr
	{
		NodeAddress:   hexToAddressPtr("0x6E595d9cE1bD44dc6f997B18D2699619AaBE9c55"),
		Treasury:      common.HexToAddress("0xD234c076D59dE386a534825deEC5421584fFeFD8"),
		OracleAddress: common.HexToAddress("0x09A2506a2B5039DbD34a84cd26724244d4b2d854"),
		ConsensusKey:  common.Hex2Bytes("b4c818fa98544d3d3a92c908c2def4334be645f2983c961c8b827d7d6a4526152f358f263339c6a272ef402b25ae669f"),
		Enode:         "enode://54aee31c8181eb6f88cac43401f1f456fba8f2ec918069a33f7beb1e46eedc062275380c975b618d0540126bbe90c6a63fd2d9c34a4400cff9d8b08a0e3a9ac0@217.66.20.45:30303",
		BondedStake:   Ntn1,
	},
	// E0y37p
	{
		NodeAddress:   hexToAddressPtr("0x82e4ed31a437a8e2EEA6012531A2c754C8c3aB82"),
		Treasury:      common.HexToAddress("0x6876123246777ABb3F840e4489D1810D678822fA"),
		OracleAddress: common.HexToAddress("0x2664d4480Df26C628F3D048Bc1B3CeaD6241Ee2B"),
		ConsensusKey:  common.Hex2Bytes("97c987b18b1ef1522357c256bc70e801c35d84a21d636c33012ae51383049ed8483ccea691e1a3799463d9278ad9f7d9"),
		Enode:         "enode://c1384f53a7ae66adb6c4b24ef1f31c7ecac0790aee85b0630efc321a24286856b99348f9221cf96d4e598d6a8a6893f4d17730284c3329f11276459664d3b5fe@35.208.175.249:30303",
		BondedStake:   Ntn1,
	},
	// GxA9v3
	{
		NodeAddress:   hexToAddressPtr("0x8264506336d4cDa914Bedb3F3e8D6577CcD1f2ae"),
		Treasury:      common.HexToAddress("0x8dEA3eFA35bB783F1F7561FC6077AF3b2162D0cd"),
		OracleAddress: common.HexToAddress("0x1D7ef3B0da315bF57a79986e23b98A60740387f1"),
		ConsensusKey:  common.Hex2Bytes("8819eb34b7505be4c8bfc56e2f2677de9e8cb2050f52d25e211f3bd6951fe287f8c5d7c8bc0664da6b68f9be542cd462"),
		Enode:         "enode://6556d646609fdc3ca323797c5edaad4a47844d6f9a9c149c800788eb7deec558c9874292d8b1e31e99bc5f032b07f0186543b2c44fafb83c4531a398333ba9b4@45.153.35.148:30250",
		BondedStake:   Ntn1,
	},
	// h3Y8fq
	{
		NodeAddress:   hexToAddressPtr("0x37c7DC11812E11Fab9127C8746a4d576F5AB5280"),
		Treasury:      common.HexToAddress("0x63abA1d8735947C5241D1b1FAF3D3c989B69Ca9c"),
		OracleAddress: common.HexToAddress("0x6Ebc8917FDE618F59310B175af1e35aDEDbFC47a"),
		ConsensusKey:  common.Hex2Bytes("8a79e8ddb666b0e59ad98c8d3c941c7b210db8712a0d2a55fb4e807a5b7b501ec161e3ed438f201a138be7af5e1c119e"),
		Enode:         "enode://3b4c22cb86a8cfa0e235c31e5d674b6a5ed3149621196b03c7467e6de6cbfb2f7f634a0c07ffbd6a2fdb18fcb7a922911c47ff7378fd86fbd17ec03540299d31@65.109.117.219:30303",
		BondedStake:   Ntn1,
	},
	// j90EwR
	{
		NodeAddress:   hexToAddressPtr("0x398601D19FDCeA40b1C5eF386e272A8b29E3A62D"),
		Treasury:      common.HexToAddress("0x57D801f412af651b9c2245852EFD78F72a1109E4"),
		OracleAddress: common.HexToAddress("0xd312E3e99ba1cAb19e5592C94B2Fe9D4dF35ea59"),
		ConsensusKey:  common.Hex2Bytes("b3a0551bfc03ce5bdc2579c123852b2220aa7b6a3dce47396ffb01db2b7e821d62ceaa7822e11813b6c4a62286600884"),
		Enode:         "enode://103bc560fb54ee649d095443a7b2934f2898e95781cfe52729cd68064d24c85886b01740a1ceb54268461606df4339fbbed1ae71407e6ff44690aeb6d1235a87@142.132.156.99:30303",
		BondedStake:   Ntn1,
	},
	// m3Yv20
	{
		NodeAddress:   hexToAddressPtr("0xae172a1b52331885Da7C3fb48F4c5fC943105e4B"),
		Treasury:      common.HexToAddress("0x5180798028A4B9742a9346d09b48C953c7750373"),
		OracleAddress: common.HexToAddress("0x65171526EEC9804124752F3f4491bb7Bd8a9d7eB"),
		ConsensusKey:  common.Hex2Bytes("96992d91ef611ac6b1c2b0ffa31de47c9f472190d5b663c07ac590073f08cbebdfa5dfa2bbb3ac77d70c282e14a5f873"),
		Enode:         "enode://69ceab52c1a30a53a67b9748232b951f47d003f03d0661dac85395713b42e0c49cab040347f6594f47adeb69602fd3a104d62dba9b2647287036138a846a8f81@50.115.46.42:30303",
		BondedStake:   Ntn1,
	},
	// n7R0cG
	{
		NodeAddress:   hexToAddressPtr("0x0f87a3c5F45dAc811A3Ba894e8e11BBE28103eE5"),
		Treasury:      common.HexToAddress("0x2F7955Ce3d9EA7f24e0635d5645Bd8760537a80d"),
		OracleAddress: common.HexToAddress("0x4df785EBf687B14d33F3f3a06d40AF38a32453E5"),
		ConsensusKey:  common.Hex2Bytes("a8ca98db2c925a2ec5da4ab6096c3c655641508aade67072c77de0c3265714995886258a6a56b6868663a613eca49366"),
		Enode:         "enode://9ee6b1a88f67c74d0c5d678a1ba91fcd05c2c0e54101d93fa6fd84862a4fb6adb7c873758b160a49f95bb81cf10135c86877b886f08cc768f5b97d3b936c45d4@135.181.215.62:30303",
		BondedStake:   Ntn1,
	},
	// P8xJ4q
	{
		NodeAddress:   hexToAddressPtr("0x769e0fa6F16f61B3F03B261279a91b882781d6bA"),
		Treasury:      common.HexToAddress("0xE6A288D9AEa9fB4A6F496b9205adfA43A281EDeD"),
		OracleAddress: common.HexToAddress("0xe8F811b17343629cB5BA4A1491F5E1305d8Bd4Bf"),
		ConsensusKey:  common.Hex2Bytes("8ed20c8beb8b9ad3578e9cd18ce2d16bd6e209c80d57daf72fe84be7f2801b6dd894fdf9470105d3be6bece84c7d1e1b"),
		Enode:         "enode://042bb47fe1a11f9c0383f53fe8907d2d7d5f404eb777dfc18d402c251f3924fb23be0e41529c8baa1b093ec6841c2259b1027f83548765f16da615a0707cd7de@190.2.141.78:30303",
		BondedStake:   Ntn1,
	},
	// Qn52Rw
	{
		NodeAddress:   hexToAddressPtr("0x111ce88feb53CE74E6468C9c31FAc820EE96BCd3"),
		Treasury:      common.HexToAddress("0x48f488a0258cd45250bf3c54c20de013c60d47e0"),
		OracleAddress: common.HexToAddress("0xc3a7015eABda3902214176f25cB2ff3025Cf43F1"),
		ConsensusKey:  common.Hex2Bytes("8406d4fa7d5e9dfd3939672852a54be6db8dd28b44a9878cf2db64941867003c8d3d99225004fd233d9c26d90b390d1b"),
		Enode:         "enode://60677fbff1dd8d072949703ac13333a3175c2b7d312584b8f6167bb4bbad4c6cefde6099e8a14a47dab3c6494b78c6bf3cd53d6d027aecc0f79b3a0e542af55d@141.94.141.80:30303",
		BondedStake:   Ntn1,
	},
	// Qw5DzV
	{
		NodeAddress:   hexToAddressPtr("0x4c5E7643026D6d131F9b5cc3e9BdF3A42271aD82"),
		Treasury:      common.HexToAddress("0x9BA01C2b5F3E128Ba27DDC6C6daA261fDa804C83"),
		OracleAddress: common.HexToAddress("0xe617D515a8C61Fad3bc7B55E468FEc43c07d693f"),
		ConsensusKey:  common.Hex2Bytes("976dc7c53e0e8fa58b75aa9c15aae2e296a08e62530209aa7acd2a291c6b3a60577cf658d4e5c24889d645748e037e03"),
		Enode:         "enode://9013cce57ebd49641ec2286d331d41c0198f012db0ca279ca402b6b45e6b2c52dc323f1f9728fbb19b3b89e029acceca25773069a0a15b864f36c03475caed0c@52.208.224.137:30281",
		BondedStake:   Ntn1,
	},
	// r7U2Nk
	{
		NodeAddress:   hexToAddressPtr("0x45D3787551499f8c23De7bde106ff05509147B74"),
		Treasury:      common.HexToAddress("0x643fE2808E3c6a91eF0a360517f9DBf39b9aA8fB"),
		OracleAddress: common.HexToAddress("0x9A6e3A7B45B084fD8E140a081001217f76985cfB"),
		ConsensusKey:  common.Hex2Bytes("872f91173d097e86e29be4c4f8a4d212edd3e852925262b7172dd8288b5733090f798f0033adcb787f5382805e59246f"),
		Enode:         "enode://243a3255a63f48cdd00ec35b9a5432abf02d3c62f6fd3255928e0a67b63213399fe5c4c26c262e3c6b57edd443b06603b381fcb66bda850b0eb0a334cb8e3368@176.9.125.13:30303",
		BondedStake:   Ntn1,
	},
	// u7PnX4
	{
		NodeAddress:   hexToAddressPtr("0x612E73E6c3232a92e1c5171f659c41400fEA6Bcc"),
		Treasury:      common.HexToAddress("0xE4101Be42F880f2536f6F6b888b358a2b86f0927"),
		OracleAddress: common.HexToAddress("0xd90256637799CC2F6Cd6a746bBACBbAAdc6c9A53"),
		ConsensusKey:  common.Hex2Bytes("83bbcb98d1dcb549865ef10cc45e9c8480a284c840be67070def85ae6c11d36d2b3e4e117147b873b595cad602e94b1f"),
		Enode:         "enode://e91757e7fb5eda843e6e04228c71632dfa3cb899f17fce2d164250beaa0fb73aa8af497356db6a8987b1fd3d05a01e716ca20d362fd0c85fee308e3a17d6b152@185.100.10.182:30310",
		BondedStake:   Ntn1,
	},
	// vJ7qX5
	{
		NodeAddress:   hexToAddressPtr("0x6e9ea56Ab25ecaB9Bcfe1a128c54cf07d0EF1400"),
		Treasury:      common.HexToAddress("0x6E129707dF5Eb03dCb92427c60B3176f8e099101"),
		OracleAddress: common.HexToAddress("0xE310dAe3Da82FDB2c92FEf88E32F7A1c533C6Afe"),
		ConsensusKey:  common.Hex2Bytes("8351826288273f290873d3086ad3584cbb9d0a8ff57a5295ba3a3287565c4e4cc6a64da12dca347936820bdeb7a11bb0"),
		Enode:         "enode://becc3c02019c3b20ba17a1f576a6563fa1e33f19edd7eea4e6a9b6d88e31a8f2d45145af84b9d241399d94e3a414488628e6f8942e9b39d47016d6c1e4f97eab@37.27.64.210:33333",
		BondedStake:   Ntn1,
	},
	// w2G8Ve
	{
		NodeAddress:   hexToAddressPtr("0x8dE949f5e6aB1108bCc7EC301F83BF55a2214441"),
		Treasury:      common.HexToAddress("0xFc226fD756e45894B88944881643C92Ff26CCd3a"),
		OracleAddress: common.HexToAddress("0x230115d5D02F92BDA7c75BAEb28566Ad6c0D3CeA"),
		ConsensusKey:  common.Hex2Bytes("b43793524261eed73c027b5068827e5fa4ee5006c87730f2495a6d50ccb1daee63dc1c49a50fb0065f76b1f75e84a2b5"),
		Enode:         "enode://716b005017770a9e885740e6ffb403b52721cea9d6e8ae3eaa14e34d1e82f2e6337bd137c83ca17e9ef52d2919773f7fee124ba04703572ea549aca13070d383@37.252.186.236:30303",
		BondedStake:   Ntn1,
	},
	// X3g9Vy
	{
		NodeAddress:   hexToAddressPtr("0x8042dD69D91099F1a07d86043CeF50E5666D7cFA"),
		Treasury:      common.HexToAddress("0x98F5DEb9C437bD7B99B61E719b4F7673aFDB51af"),
		OracleAddress: common.HexToAddress("0x3003Be771cD6B8Dffb428a1f6fdC3B55cd06112f"),
		ConsensusKey:  common.Hex2Bytes("84d9c46dc8f751d68ac59f257867657b32a651be54098b36fccd7772aa7a6da94ee3854af9daf67b91ad304a3462e6be"),
		Enode:         "enode://f867013d367c445a4735b04a7b9d381679583b0123fa3624364eadbe522ee7292f58e1274e7ee86586625837e8c12e803f9f78d728320404230f8d4623d10a68@65.21.47.31:30303",
		BondedStake:   Ntn1,
	},
	// xGq8y2
	{
		NodeAddress:   hexToAddressPtr("0x2e2E10E811b3dD5691ab2744E694072CD749BDe7"),
		Treasury:      common.HexToAddress("0xfF2da72bA8e4759b24cA700886545caD7e76d6f5"),
		OracleAddress: common.HexToAddress("0xDEf2025e4d249E631f1fFD93DfDA0bCDecf0e6B9"),
		ConsensusKey:  common.Hex2Bytes("91af6e5b868f8288ea1425893867f8aac55ee31fe2bdaebfe1e27e5130151faaec9d8f6cf6c10277cefb1e88ff25376e"),
		Enode:         "enode://dcf007c459160b3e53c7cced0eec41fce9ccbcec97984811a95536a14b760614f28f8e9a7aae3814c824bc6632cf3e22aa6f453f2eba8cb82635396034e6cb7f@35.217.1.95:30303",
		BondedStake:   Ntn1,
	},
	// yZ8w3C
	{
		NodeAddress:   hexToAddressPtr("0x50cf23B31C78e65552EC338C5B6a32D65447A2e3"),
		Treasury:      common.HexToAddress("0x37F2803d217cCD63E270eB73e8B1F840df8EB97A"),
		OracleAddress: common.HexToAddress("0xd1F8d57C73C99A57054e5F850b79fa2B5c71d0Fd"),
		ConsensusKey:  common.Hex2Bytes("addf8fac1fcf3727b0ed59437fd43ed1bdf5a6f60ee73d1bdaaefe8dfed1e44dbd21022b516dd6c8fb0440ae013a5c96"),
		Enode:         "enode://d9f49af41de6d6db4171943d18077798e059d56f56611636d375009d98e69e3995b93ea4592e4f72e2ef813cd0aa7f979448a967607e86b34ef410a35c303bab@190.2.149.83:30303",
		BondedStake:   Ntn1,
	},
	// z90WxC
	{
		NodeAddress:   hexToAddressPtr("0x05e5417c65f5F81BF6dCcD3bDEd30FA779A5ec78"),
		Treasury:      common.HexToAddress("0x19797Fac59E59fd2daD9aEc74deD7d104904faC3"),
		OracleAddress: common.HexToAddress("0xDA543790CF7af44b90c875d46337C518C1A0B415"),
		ConsensusKey:  common.Hex2Bytes("84203e7a29b9c1c53bbecbbf6012cf97350d42b4b9eccc9800dcdba1244014d2ed8d879d8ca52f921fc937c5a268b7fb"),
		Enode:         "enode://69d236c308349a1894719dadf3ca97475cdfa68d9c60aa899f1b45ff2dfa09607cbf38dde6fc750633836ade690f2a721eef1d83c25abb833c06d721070a05f4@149.50.116.70:30303",
		BondedStake:   Ntn1,
	},
}

const AutMainnetNetworkID = 65_000_000

// MainnetChainConfig contains the chain parameters to run a node on the Autonity main network.
var AutMainnetChainConfig = &ChainConfig{
	ChainID:                 big.NewInt(AutMainnetNetworkID),
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
		Operator:                 common.HexToAddress("0x5fB82096CdFc95755b7b766E94F80265C9Fc4bFC"),
		Treasury:                 common.HexToAddress("0x37563754deD96B4dC645fe10Ebf6d3A8be56555F"),
		WithheldRewardsPool:      common.HexToAddress("0x37563754deD96B4dC645fe10Ebf6d3A8be56555F"),
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
		Validators:          MainnetValidators,
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
						new(big.Int).Mul(big.NewInt(551), big.NewInt(Ether)),
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
