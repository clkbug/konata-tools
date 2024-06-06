package symtab_test

import (
	"testing"

	"github.com/clkbug/konata-tools/symtab"
)

func testSetup() (symtab.SymbolTable, error) {
	const filename = "coremark.symtab.txt"
	return symtab.ParseFile(filename)
}
func TestParseSymbolTableFile(t *testing.T) {
	st, err := testSetup()
	if err != nil {
		t.Error(err)
	}
	if len(st) != 107 {
		t.Errorf("len(st), expected 107, but got %d", len(st))
	}
	f := 0
	for _, x := range st {
		if x.Type == symtab.SymFunc {
			f++
		}
	}
	if f != 43 {
		t.Errorf("len(filter(fun x -> x.Type == SymFunc, st)) == 43, but got %d", f)
	}
}
func TestSymbolTableSort(t *testing.T) {
	st, err := testSetup()
	if err != nil {
		t.Error(err)
	}
	st.Sort()
	var cur uint64
	for _, x := range st {
		if x.Addr < cur {
			t.Errorf("Addr(%d) < %d (not sorted!)", x.Addr, cur)
		}
		cur = x.Addr
	}
}
func TestSymbolTableSearch(t *testing.T) {
	st, err := testSetup()
	if err != nil {
		t.Error(err)
	}
	st.Sort()
	testTable := []struct {
		addr     uint64
		expected string
	}{
		{0x0000476c, "number"},
		{0x0000476c + 1564/2, "number"},
		{0x00001dac + 1508/3, "core_bench_list"},
		{0x00004758, "portable_init"},
		{0x00004758 + 1, "portable_init"},
		{0x00004758 + 11, "portable_init"},
	}
	for _, x := range testTable {
		s, ok := st.Search(x.addr)
		if !ok {
			t.Errorf("Search(%d) expected %s, but not found", x.addr, x.expected)
		}
		if s.Name != x.expected {
			t.Errorf("Search(%d).Name expected %s, but got %s", x.addr, x.expected, s.Name)
		}
	}
	notFounds := []uint64{
		100,
		200,
		0xFFFFFFFF,
	}
	for _, x := range notFounds {
		s, ok := st.Search(x)
		if ok {
			t.Errorf("Search(%d) should be not found, but got %v", x, s)
		}
	}
}
