package konata_test

import (
	"fmt"
	"testing"

	"github.com/clkbug/konata-tools"
)

func TestParserAndStringer(t *testing.T) {
	inputs := []string{
		"kanata-sample-1.log",
		"kanata-sample-2.log",
	}
	for _, input := range inputs {
		insts, err := konata.ParseFile(input)
		if err != nil {
			t.Errorf("failed to ParseFile(%s): %s", input, err)
		}
		fmt.Printf("%s: %d\n", input, len(insts))
	}
}
