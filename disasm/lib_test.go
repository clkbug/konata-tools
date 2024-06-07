package disasm_test

import (
	"testing"

	"github.com/clkbug/konata-tools/disasm"
)

func TestParseFile(t *testing.T) {
	for _, f := range []string{"coremark.dump", "coremarkE.dump"} {
		p, err := disasm.ParseFile(f)
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
}

func TestSearch(t *testing.T) {
	for _, f := range []string{"coremark.dump", "coremarkE.dump"} {
		p, err := disasm.ParseFile(f)
		if err != nil {
			t.Fatal(err)
		}

		testTable := []disasm.Instruction{
			{PC: 0x1000, FunName: "_start", Asm: "j	1010 <_init>"},
			{PC: 0x1004, FunName: "_end", Asm: "j	1004 <_end>"},
			{PC: 0x1008, FunName: "_call_main", Asm: "jal	62e0 <main>"},
			{PC: 0x100c, FunName: "_call_main", Asm: "j	1004 <_end>"},
			{PC: 0x4e18, FunName: "ee_printf", Asm: "sw	s1,356(sp)"},
			{PC: 0x62bc, FunName: "ee_printf", Asm: "j	5500 <ee_printf+0x6ec>"},
			{PC: 0x6300, FunName: "main", Asm: "sw	s6,2000(sp)"},
		}

		for i, expected := range testTable {
			got, found := p.Search(expected.PC)
			if !found {
				t.Errorf("%d: expected %v, but not found", i, expected)
			}
			if got != expected {
				t.Errorf("%d: expected %v, but got %v", i, expected, got)
			}
		}
	}
}
