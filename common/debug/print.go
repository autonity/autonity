package debug

import (
	"bytes"
	"fmt"

	"github.com/JekaMas/pretty"
)

func Diff(obtained, expected interface{}) {
	var failMessage bytes.Buffer
	diffs := pretty.Diff(obtained, expected)

	if len(diffs) > 0 {
		failMessage.WriteString("Obtained:\t\tExpected:")
		for _, singleDiff := range diffs {
			failMessage.WriteString(fmt.Sprintf("\n%v", singleDiff))
		}
	}

	res := failMessage.String()
	if len(res) == 0 {
		fmt.Println("Objects are identical")
	} else {
		fmt.Printf("Diff %s\n", failMessage.String())
	}
}
