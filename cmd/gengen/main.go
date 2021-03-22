package main

import (
	"fmt"
	"os"

	"github.com/clearmatics/autonity/cmd/gengen/gengen"
)

func main() {
	err := gengen.NewCmd().Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

}
