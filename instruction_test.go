package konata_test

import (
	"bufio"
	"os"
	"strings"
	"testing"

	"github.com/clkbug/konata-tools"
)

var inputs = []string{
	"kanata-sample-1.log",
	"kanata-sample-2.log",
}

func TestParser(t *testing.T) {
	for _, input := range inputs {
		_, err := konata.ParseFile(input)
		if err != nil {
			t.Errorf("failed to ParseFile(%s): %s", input, err)
		}
	}
}
func TestStringer(t *testing.T) {
	for _, input := range inputs {
		insts, err := konata.ParseFile(input)
		if err != nil {
			t.Errorf("failed to ParseFile(%s): %s", input, err)
		}

		fp, err := os.Open(input)
		if err != nil {
			t.Error(err)
		}
		buf := bufio.NewScanner(fp)
		i := 0
		if !buf.Scan() || !strings.HasPrefix(buf.Text(), "Kanata\t0004") {
			t.Errorf("invalid file header %s", input)
		}
		for buf.Scan() {
			l := buf.Text()
			lt := strings.TrimSpace(l)

			r := insts[i]
			if l != r.String() && lt != r.String() {
				t.Errorf("'%s' != '%s'[%v.String()]", l, r.String(), r)
			}
			i++
		}
	}
}
