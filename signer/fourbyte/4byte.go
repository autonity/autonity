// Code generated by go-bindata. DO NOT EDIT.
// sources:
// 4byte.json (6.048kB)

package fourbyte

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type asset struct {
	bytes  []byte
	info   os.FileInfo
	digest [sha256.Size]byte
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

//nolint:misspell
var __4byteJson = []byte(`{
  "2f2c3f2e": "COMMISSION_RATE_PRECISION()",
  "f56b5615": "COMMITTEE_FRACTION_PRECISION()",
  "5115840b": "PROPOSER_REWARD_RATE_PRECISION()",
  "b46e5520": "activateValidator(address)",
  "dd62ed3e": "allowance(address,address)",
  "095ea7b3": "approve(address,uint256)",
  "e98712b2": "atnTotalRedistributed()",
  "70a08231": "balanceOf(address)",
  "a515366a": "bond(address,uint256)",
  "9dc29fac": "burn(address,uint256)",
  "852c4849": "changeCommissionRate(address,uint256)",
  "872cf059": "completeContractUpgrade()",
  "ae1f5fa0": "computeCommittee()",
  "79502c55": "config()",
  "313ce567": "decimals()",
  "d5f39488": "deployer()",
  "c9d97af4": "epochID()",
  "1604e416": "epochReward()",
  "9c98e471": "epochTotalBondedStake()",
  "3b2f2fac": "finalize(address[],address,uint256,bool)",
  "d861b0e8": "finalizeInitialization()",
  "43645969": "getBlockPeriod()",
  "ab8f6ffe": "getCommittee()",
  "a8b2216e": "getCommitteeEnodes()",
  "b45e2675": "getConsensusViewOfHeight(uint256)",
  "96b477cb": "getEpochFromBlock(uint256)",
  "dfb1a4d2": "getEpochPeriod()",
  "731b3a03": "getLastEpochBlock()",
  "819b6463": "getMaxCommitteeSize()",
  "11220633": "getMinimumBaseFee()",
  "b66b3e79": "getNewContract()",
  "e7f43c68": "getOperator()",
  "833b1fce": "getOracle()",
  "5f7d3949": "getProposer(uint256,uint256)",
  "b6ef8c67": "getRevertingAmount(uint256)",
  "f7866ee3": "getTreasuryAccount()",
  "29070c6d": "getTreasuryFee()",
  "6fd2c80b": "getUnbondingPeriod()",
  "a9af1059": "getUnbondingReleaseState(uint256)",
  "1904bb2e": "getValidator(address)",
  "b7ab4db5": "getValidators()",
  "0d8e6e2c": "getVersion()",
  "9ac31074": "inflationReserve()",
  "c2362dd5": "lastEpochBlock()",
  "89c614b8": "lastEpochTime()",
  "cef98450": "maxBondAppliedGas()",
  "5bdedd16": "maxRewardsDistributionGas()",
  "386a827b": "maxUnbondAppliedGas()",
  "71d1bc59": "maxUnbondReleasedGas()",
  "40c10f19": "mint(address,uint256)",
  "06fdde03": "name()",
  "0ae65e7a": "pauseValidator(address)",
  "161605e3": "receiveATN()",
  "84467fdb": "registerValidator(string,address,bytes,bytes)",
  "cf9c5719": "resetContractUpgrade()",
  "1250a28d": "setAccountabilityContract(address)",
  "d372c07e": "setAcuContract(address)",
  "8bac7dad": "setCommitteeSize(uint256)",
  "6b5f444c": "setEpochPeriod(uint256)",
  "f03b959d": "setInflationControllerContract(address)",
  "0fe50109": "setMaxBondAppliedGas(uint256)",
  "7e660ac9": "setMaxRewardsDistributionGas(uint256)",
  "0b21fb1d": "setMaxUnbondAppliedGas(uint256)",
  "ee7ff28d": "setMaxUnbondReleasedGas(uint256)",
  "cb696f54": "setMinimumBaseFee(uint256)",
  "1a0cf2e5": "setNonStakableVestingContract(address)",
  "b47c3eaa": "setOmissionAccountabilityContract(address)",
  "520fdbbc": "setOperatorAccount(address)",
  "496ccd9b": "setOracleContract(address)",
  "cfd19fb9": "setStabilizationContract(address)",
  "6a929cef": "setStakingGasPrice(uint256)",
  "b3ecbadd": "setSupplyControlContract(address)",
  "d886f8a2": "setTreasuryAccount(address)",
  "77e741c7": "setTreasuryFee(uint256)",
  "114eaf55": "setUnbondingPeriod(uint256)",
  "ceaad455": "setUpgradeManagerContract(address)",
  "6bb9cb0e": "stakingGasPrice()",
  "95d89b41": "symbol()",
  "18160ddd": "totalSupply()",
  "a9059cbb": "transfer(address,uint256)",
  "23b872dd": "transferFrom(address,address,uint256)",
  "a5d059ca": "unbond(address,uint256)",
  "784304b5": "updateEnode(address,string)",
  "35be16e0": "updateValidatorAndTransferSlashedFunds((address,address,address,string,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,address,uint256,uint256,uint256,uint256,uint256,bytes,uint8))",
  "b2ea9adb": "upgradeContract(bytes,string)",
  "187cf4d7": "FEE_FACTOR_UNIT_RECIP()",
  "372500ab": "claimRewards()",
  "5ea1d6f8": "commissionRate()",
  "282d3fdf": "lock(address,uint256)",
  "59355736": "lockedBalanceOf(address)",
  "a0ce552d": "redistribute(uint256)",
  "19fac8fd": "setCommissionRate(uint256)",
  "61d027b3": "treasury()",
  "949813b8": "unclaimedRewards(address)",
  "7eee288d": "unlock(address,uint256)",
  "84955c88": "unlockedBalanceOf(address)",
  "3a5381b5": "validator()",
  "4dc925d3": "ACCUSATION_CONTRACT()",
  "2090a442": "COMPUTE_COMMITTEE_CONTRACT()",
  "c13974e1": "ENODE_VERIFIER_CONTRACT()",
  "8e153dc3": "INNOCENCE_CONTRACT()",
  "925c5492": "MISBEHAVIOUR_CONTRACT()",
  "50d93720": "POP_VERIFIER_CONTRACT()",
  "d0a6d1a6": "SUCCESS()",
  "a4ad5d91": "UPGRADER_CONTRACT()",
  "55463ceb": "autonity()",
  "570ca735": "operator()",
  "b3ab15fb": "setOperator(address)",
  "6e3d9ff0": "upgrade(address,string)",
  "7adbf973": "setOracle(address)",
  "a2e62045": "update()",
  "7ecc2b56": "availableSupply()",
  "44df8e70": "burn()",
  "db7f521a": "setStabilizer(address)",
  "7e47961c": "stabilizer()",
  "a8031a1d": "distributeRewards(address,uint256)",
  "6c9789b0": "finalize(bool)",
  "92eff3cd": "calculateSupplyDelta(uint256,uint256,uint256,uint256)",
  "f968f493": "unlockTokens()",
  "eeb92233": "distributeProposerRewards(uint256)",
  "c1a48245": "finalize(address[],address,uint256,bool,bool)",
  "9a11e0e6": "getInactivityScore(address)",
  "5ca1809c": "getLookbackWindow()",
  "7f5e2f11": "getScaleFactor()",
  "d5baf908": "setCommittee(address[],address[])",
  "e1e8cac6": "setLastEpochBlock(uint256)",
  "4bb278f3": "finalize()",
  "9670c0bc": "getPrecision()",
  "9f8743f7": "getRound()",
  "3c8510fd": "getRoundData(uint256,string)",
  "df7f710e": "getSymbols()",
  "b78dec52": "getVotePeriod()",
  "cdd72253": "getVoters()",
  "33f98c77": "latestRoundData(string)",
  "8d4f75d2": "setSymbols(string[])",
  "845023f2": "setVoters(address[])",
  "307de9b6": "vote(uint256,int256[],uint256)",
  "9dfd1b8e": "bondingApplied(uint256,address,uint256,bool,bool)",
  "d18736ab": "rewardsDistributed(address[])",
  "a8920241": "unbondingApplied(uint256,address,bool)",
  "3c54c290": "unbondingReleased(uint256,uint256,bool)",
  "416f86a2": "setProposerRewardRate(uint256)",
  "39f9eb1f": "setWithheldRewardsPool(address)",
  "53b1821b": "getTotalEffort()"
}
`)

func _4byteJsonBytes() ([]byte, error) {
	return __4byteJson, nil
}

func _4byteJson() (*asset, error) {
	bytes, err := _4byteJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "4byte.json", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xa9, 0x2a, 0x78, 0x82, 0xb3, 0xaa, 0x5f, 0xa3, 0xd7, 0x20, 0xb3, 0x89, 0x59, 0xa9, 0x3b, 0x41, 0xf7, 0xc3, 0x7d, 0xcc, 0xa5, 0xb0, 0x71, 0x20, 0x32, 0x88, 0xc5, 0x86, 0xe5, 0x9a, 0x83, 0xaf}}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// AssetString returns the asset contents as a string (instead of a []byte).
func AssetString(name string) (string, error) {
	data, err := Asset(name)
	return string(data), err
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// MustAssetString is like AssetString but panics when Asset would return an
// error. It simplifies safe initialization of global variables.
func MustAssetString(name string) string {
	return string(MustAsset(name))
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetDigest returns the digest of the file with the given name. It returns an
// error if the asset could not be found or the digest could not be loaded.
func AssetDigest(name string) ([sha256.Size]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s can't read by error: %v", name, err)
		}
		return a.digest, nil
	}
	return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s not found", name)
}

// Digests returns a map of all known files and their checksums.
func Digests() (map[string][sha256.Size]byte, error) {
	mp := make(map[string][sha256.Size]byte, len(_bindata))
	for name := range _bindata {
		a, err := _bindata[name]()
		if err != nil {
			return nil, err
		}
		mp[name] = a.digest
	}
	return mp, nil
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"4byte.json": _4byteJson,
}

// AssetDebug is true if the assets were built with the debug flag enabled.
const AssetDebug = false

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//
//	data/
//	  foo.txt
//	  img/
//	    a.png
//	    b.png
//
// then AssetDir("data") would return []string{"foo.txt", "img"},
// AssetDir("data/img") would return []string{"a.png", "b.png"},
// AssetDir("foo.txt") and AssetDir("notexist") would return an error, and
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		canonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(canonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"4byte.json": {_4byteJson, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory.
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	return os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
}

// RestoreAssets restores an asset under the given directory recursively.
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(canonicalName, "/")...)...)
}
