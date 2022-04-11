package main

import (
	"fmt"
	"os"

	"github.com/autonity/autonity/cmd/gengen/gengen"
)

func main() {
	err := gengen.NewCmd().Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

}
