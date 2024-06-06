package disasm_test

import (
	"testing"

	"github.com/clkbug/konata-tools/disasm"
)

func TestParseFile(t *testing.T) {
	p, err := disasm.ParseFile("coremark.dump")
	if err != nil {
		t.Fatal(err)
	}

	testTable := []disasm.Instruction{
		{PC: 0x1000, FunName: "_start", Asm: "j	1010 <_init>"},
		{PC: 0x1004, FunName: "_end", Asm: "j	1004 <_end>"},
		{PC: 0x1008, FunName: "_call_main", Asm: "jal	62e0 <main>"},
		{PC: 0x100c, FunName: "_call_main", Asm: "j	1004 <_end>"},
	}

	for i, expected := range testTable {
		if p[i] != expected {
			t.Errorf("%d: expected %v, but got %v", i, expected, p[i])
		}
	}
}
