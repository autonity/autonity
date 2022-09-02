package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/autonity/autonity/common/hexutil"
	"github.com/autonity/autonity/crypto"
)

// generates entries selector -> function signature from the ABI of a contract
func generateSelectors(ABIpath string) map[string]string {
	// read contract ABI
	content, err := os.ReadFile(ABIpath)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	var payload []map[string]interface{}
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	// extract function signatures from ABI
	var signatures []string
	var signature string
	var argType string
	var inputsArray []interface{}
	var input map[string]interface{}
	for i := 0; i < len(payload); i++ {
		if payload[i]["type"] == "function" {
			signature = fmt.Sprintf("%v(", payload[i]["name"])
			inputsArray = payload[i]["inputs"].([]interface{})
			// get function arguments types
			for j := 0; j < len(inputsArray); j++ {
				input = inputsArray[j].(map[string]interface{})
				argType = fmt.Sprintf("%v", input["type"])
				if j != len(inputsArray)-1 {
					argType += ","
				}
				signature += argType
			}
			signature += ")"
			signatures = append(signatures, signature)
		}
	}

	// compute 4byte selector for each function signature
	selectors := make(map[string]string)
	for i := 0; i < len(signatures); i++ {
		hash := crypto.Keccak256Hash([]byte(signatures[i]))
		hashBytes := hash.Bytes()
		selectorBytes := hashBytes[0:4]
		selectorHex := hexutil.Encode(selectorBytes)
		selectors[selectorHex[2:]] = signatures[i]
	}

	return selectors
}

func main() {

	// extract Autonity and Liquid contract selectors
	autonitySelectors := generateSelectors("./common/acdefault/generated/Autonity.abi")
	liquidSelectors := generateSelectors("./common/acdefault/generated/Liquid.abi")

	// load existing clef 4byte selectors
	content, err := os.ReadFile("./signer/fourbyte/4byte.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	var clefSelectors map[string]string
	err = json.Unmarshal(content, &clefSelectors)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	// merge maps
	for k, v := range liquidSelectors {
		clefSelectors[k] = v
	}
	for k, v := range autonitySelectors {
		clefSelectors[k] = v
	}

	// write custom 4byte db to file
	clefSelectorsJSON, err := json.Marshal(clefSelectors)
	if err != nil {
		log.Fatal("Error during Marshal: ", err)
	}
	err = os.WriteFile("./signer/fourbyte/4byte.json", clefSelectorsJSON, 0600)
	if err != nil {
		log.Fatal("Error while writing file: ", err)
	}
}
