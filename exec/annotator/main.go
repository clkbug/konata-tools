package main

import (
	"fmt"
	"os"

	"github.com/clkbug/konata-tools"
	"github.com/clkbug/konata-tools/rInst"
)

func run() error {
	if len(os.Args) == 2 {
		cmds, err := konata.ParseFile(os.Args[1])
		if err != nil {
			return err
		}
		prog, err := rInst.ToProgram(cmds)
		if err != nil {
			return err
		}
		for i, inst := range prog {
			fmt.Printf("%d\t%08x\t%v\n", i, inst.PC, inst)
		}
		return nil
	} else {
		return fmt.Errorf("! os.Args != 2")
	}
}

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}
