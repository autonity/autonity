package keygenerator

import (
	"crypto/ecdsa"
	"github.com/clearmatics/autonity/crypto"
)

func init() {
	root = New()
}

var keyList = []string{
	"98a97816178ad8fffef02446eb0aef59cdca43bfc34f3d25cb04dc320ba14778",
	"4d531ca75258a85a150e799826c86a3d815df095dbd4256019a8ab737deba029",
	"7ddf4be8b3f671bbedf5190e9b411f389ff05eec3bd64c2b7baefba20e0e2676",
	"fbe928da20f970153374ea4b90f49c98c9ccdbb372580e260f403f0bf29a89e1",
	"3738f2b598f13ca77a48f8dd82a8e0d3f9fa71116b75ec3e03ba7e6d042791ff",
	"a8c0c66ff46544cd94be5d743b618129e3f166f6be95f73d25ddc91db452be9d",
	"6230477f89c93cd4bfe3b7d13bc9aa9743446a2cbf7ef1fc3db2ffcf4fe783f8",
	"c8c6d4121c6c1e46d8a97d2bc773cf728d0326173409941f9dafff1f5602539e",
	"962abbd3f22855b4559bb99598de1f19e43776e0aec6c756778423c359de7d7b",
	"d43f88073cb4fe44f7fd4208827c1d431d66833e9f0cc10c53a31ce0346256a9",
	"5fc5ff9e47b73180a654ac3e20c344147fb57cc968d609ef96ea8e387fd6d564",
	"ce44a2546699d2a2d0593ef6b382eb6f5119700e43018bab92079a2c1c36c04c",
	"76e96c1dd79975a0105fba36e1bbdef1398faafbae296c5cb72feb6738d76672",
	"abfa66e8d6bac749402191ee8dd260a02d16123296f451bc507845d5859489bf",
	"fc1685d961b8fa110b499ab2a99b97332bda3115c48ec46b39725678a6000efe",
	"8cb558c74f6247eba6645f2e59660cfed33ea58ffb18c1b0c70c4119436a749d",
	"91874ade765d95cbda8d74e06925b22dc4d0dd66e3ac74599e9960aefd437f98",
	"35a35b0316dcfca554126e62309c3cfbc7a2e14fd6da02c9989e60bcc216ec1d",
	"3c29cff57380ec82d561ce7a9387f853381cf5fcc1689c2582a3c169ba1b80eb",
	"58cc1590d11860f9a393a183252e56f4c76064680ea2b488964938d4c4c0df96",
	"5b5fe176c8c51a79d009b36ae4653dcf377d54ca4ef53c0264463a91841a3f82",
	"cfc10285eef0bda61755d7311157f575e4304c2898816599cc56aeb2e86347eb",
	"58fac33e75efc6f968ef27b8089a19d6cb76de6a1a72e1a2057b07aacf01f200",
	"dc5ffbcec41b8a88ad20e16480c83640967c868c1de2b7abcf2a0459bb98a7e0",
	"2a5f9854de0b7d7339996661b50a591b14e788fef1cdd28b126a823340a38f93",
	"c017585878518bbbcf55ec0a45bc77f5f8044915367762d953f91c255f4b43ca",
	"4d7fdd124528d55f39a0d138faa449b55490b1b90a7ab3dc1247dcdfb32d5fb9",
	"dca81a8bd059b9bdb8039c3362c449c0574fbe952f5c9ce09d0d597ceb26a8ff",
	"fdf8b01e768bbca2c3e1e8944216b2797c6ff39fe0e9204817a9589af08cdc51",
	"f59ea4d883b047bcfdc24b95ba361235aa6bef8a42961e22679a89b205c669f5",
	"ed5839529b3e34de8124afb2927cfb438a63e84b5840d00cecd304fe34d1d5a0",
	"1e63d8a27e2f2904578e4b71d0886015b863a4334be37ae46d0aed14a6eab26f",
	"fd7c0f28791c25ff4cb0be844d0cf09ec891d37d7aa0bc2e0c908a5b0e6317ec",
	"01147a52ab73bb22ef1827dc8b2f135102f0099a29ff488cbe631fe3b274f55b",
	"4a82b510ab33d8ad75dec8a5881e63f01331de898c696ca3584a192e278ea814",
	"8b42c38c95c40422b43fc52ecb0ea2d63db659453f60c6509b404baa7bbe1ead",
	"f9d056979a78c28a54a729e696787e3e3ba66e8d405d2f7450d5b149424cfa5e",
	"8d9514d349f73308f033f4f06e8b319aa04e663b37a10bfdddaef0af53f83922",
	"9ae193cb0a9356876b87ea5fe8edb8500fbc9331056d0167f2e91312f3cd8519",
	"1ed71a733c5a542ab767734603753fa8bbe730f63ae70b24937436265cfd5b67",
	"6913a48610911c766a70e2f3331762c798d7d00c2345ae60648f58a35efac2f6",
	"f266cd15787b3db960eb23c7d86e0d9313b3978beec053c36a2134ebb15fa7e3",
	"b853edf646842deccad8cb217429aee70427aed52835347c64fb3718e96bc66d",
	"1ce4cd30252dace83bc5a19b96b4b80ed3e9aea133551a511118c066050415f7",
	"81cb06bc238bd7d6f108e8f538b6f8be16e92c5ffcf00981b593e1cea180d50e",
	"f818c9c3779227e953906a0dfde8922cdb3c754c476011e88d73ea475c858cce",
	"eaed9f32da51e3f1b5f32d096ec4501c0027db5efc89f79acae0390794b82c59",
	"e09e7515e339f3b542946c2546c3ea56b47b30cda0182e12c14e66ca56d97afd",
	"740252189c63047f853516a2ba711e80d55575e1526302ad17cc71963c7e84e8",
	"4be96c6c612b40479b2841138dba771764d16db3f0955dd9cdff91ab7bfccc2a",
}
var root *Generator

func New() *Generator {
	return &Generator{
		keys:  keyList,
		index: 0,
	}
}

type Generator struct {
	keys  []string
	index int
}

func (gen *Generator) Next() (*ecdsa.PrivateKey, error) {
	key, err := crypto.HexToECDSA(gen.keys[gen.index])
	gen.index++
	if gen.index >= len(gen.keys) {
		gen.index -= len(gen.keys)
	}
	return key, err
}

//Next returns new private key
func Next() (*ecdsa.PrivateKey, error) {
	return root.Next()
}
