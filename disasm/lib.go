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
	inText := false

	for buf.Scan() {
		l := buf.Text()
		if !strings.Contains(l, ":") {
			continue
		}

		if strings.Contains(l, "Disassembly of section") || strings.Contains(l, "セクション") {
			inText = strings.Contains(l, ".text")
			continue
		}

		if !inText {
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

func (p Program) Search(addr uint64) (Instruction, bool) {
	low := 0
	high := len(p) - 1
	if addr < p[low].PC || p[high].PC < addr {
		return Instruction{}, false
	}
	for 1 < high-low {
		mid := (low + high) / 2
		if p[mid].PC <= addr {
			low = mid
		} else {
			high = mid
		}
	}
	if p[low].PC == addr {
		return p[low], true
	}
	return Instruction{}, false
}
