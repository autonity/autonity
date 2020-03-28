package test

import (
	"testing"
)

func TestManip(t *testing.T) {

	// cases := []*testCase{
	// 	{
	// 		name:          "no malicious",
	// 		numValidators: 2,
	// 		numBlocks:     1,
	// 		txPerPeer:     1,
	// 		finalAssert: func(t *testing.T, validators map[string]*testNode) {
	// 			for _, v := range validators {
	// 				v.listener[0].Addr()
	// 				client, err := ethclient.Dial("http://" + v.listener[1].Addr().String())
	// 				if err != nil {
	// 					t.Fatal(err)
	// 				}
	// 				addr, err := client.AutonityContractAddresss(context.Background())
	// 				if err != nil {
	// 					t.Fatal(err)
	// 				}
	// 				fmt.Printf("%v", addr)
	// 				b, err := client.BlockByNumber(context.Background(), big.NewInt(0))
	// 				if err != nil {
	// 					t.Fatal(err)
	// 				}

	// 				if b.Transactions().Len() == 0 {
	// 					t.Fatalf("Expecting block to have transactions")
	// 				}
	// 			}
	// 		},
	// 	},
	// }

	// for _, testCase := range cases {
	// 	testCase := testCase
	// 	t.Run(fmt.Sprintf("test case %s", testCase.name), func(t *testing.T) {
	// 		runTest(t, testCase)
	// 	})
	// }
}
