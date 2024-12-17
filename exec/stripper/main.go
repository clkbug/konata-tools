package main

import (
	"bufio"
	"os"
	"strconv"

	"github.com/clkbug/konata-tools"
	"github.com/clkbug/konata-tools/kInst"
)

func main() {
	filename := os.Args[1]
	cmds, err := konata.ParseFile(filename)
	if err != nil {
		panic(err)
	}
	progs, err := kInst.ToProgram(cmds)
	if err != nil {
		panic(err)
	}

	start, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(err)
	}
	end, err := strconv.Atoi(os.Args[3])
	if err != nil {
		panic(err)
	}

	stripped := progs[start:end]

	cmds2, err := stripped.ToCommand()
	if err != nil {
		panic(err)
	}

	fp, err := os.Create(os.Args[4])
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	w := bufio.NewWriter(fp)
	defer w.Flush()
	w.WriteString("Kanata\t004\n")
	w.WriteString("C=\t0\n")
	for _, c := range cmds2 {
		w.WriteString(c.String() + "\n")
	}

}
