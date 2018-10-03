package soma

import (
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/params"
)

func deployContract(config *params.SomaConfig, db ethdb.Database) {
	conf := *config
	contractBytecode := common.Hex2Bytes(conf.Bytecode[2:]) // [2:] removes 0x
	log.Printf("\nconf:\n\t%#v\n", contractBytecode)

	sdb := state.NewDatabase(db)
	statedb, _ := state.New(common.Hash{}, sdb)
	log.Println(statedb)
}
