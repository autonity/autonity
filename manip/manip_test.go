package main

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/clearmatics/autonity/ethclient"
)

func TestManip(t *testing.T) {
	t.SkipNow()
	client, err := ethclient.Dial("http://localhost:6000")
	if err != nil {
		t.Fatal(err)
	}
	b, err := client.BlockByNumber(context.Background(), big.NewInt(2))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", b)
}
