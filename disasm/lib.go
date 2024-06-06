package disasm

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Instruction struct {
	PC      uint64
	FunName string
	Asm     string
}

type Program []Instruction

func ParseFile(filename string) (Program, error) {
	fp, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	buf := bufio.NewScanner(fp)
	f := ""
	var prog Program

	for buf.Scan() {
		l := buf.Text()
		if !strings.Contains(l, ":") {
			continue
		}

		if strings.HasPrefix(l, " ") {
			// instruction
			l = strings.TrimSpace(l)
			w := strings.Split(l, ":\t")

			pc, err := strconv.ParseUint(w[0], 16, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to strconvParseUint(\"%s\") %s, in '%s'", w[0], err.Error(), l)
			}
			w2 := strings.SplitN(w[1], "\t", 2)
			var asm string
			if len(w2) == 1 {
				asm = w2[0]
			} else {
				asm = w2[1]
			}
			prog = append(prog, Instruction{PC: pc, FunName: f, Asm: asm})
		} else if strings.Contains(l, "<") {
			// function
			w := strings.Split(l, "<")
			f = strings.TrimSuffix(w[1], ">:")
		}
	}
	return prog, nil
}
