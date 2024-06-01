package symtab

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type SymType int

const (
	SymNoType SymType = iota
	SymFunc
	SymObj
	SymSection
	SymFile
)

type Symbol struct {
	Name string
	Addr uint64
	Size uint64
	Type SymType
}

type SymbolTable []Symbol

func ParseSymbolTableFile(filename string) (SymbolTable, error) {
	fp, err := os.Open(filename)
	if err != nil {
		return SymbolTable{}, err
	}
	defer fp.Close()

	st := SymbolTable{}
	buf := bufio.NewScanner(fp)
	i := 0
	for buf.Scan() {
		i++
		l := buf.Text()
		if i <= 3 {
			continue
		}
		l = strings.TrimSpace(l)
		for range 10 {
			l = strings.Replace(l, "  ", " ", -1)
		}
		w := strings.Split(l, " ")

		addr, err := strconv.ParseUint(w[1], 16, 64)
		if err != nil {
			return SymbolTable{}, fmt.Errorf("failed to ParseUint(%s) in %s:%d '%s'", w[1], filename, i, l)
		}
		size, err := strconv.ParseUint(w[2], 16, 64)
		if err != nil {
			return SymbolTable{}, fmt.Errorf("failed to ParseUint(%s) in %s:%d '%s'", w[1], filename, i, l)
		}
		var t SymType
		switch w[3] {
		case "NOTYPE":
			t = SymNoType
		case "FILE":
			t = SymFile
		case "SECTION":
			t = SymSection
		case "FUNC":
			t = SymFunc
		case "OBJECT":
			t = SymObj
		default:
			return SymbolTable{}, fmt.Errorf("unknown SymbolType '%s'", w[3])
		}
		if len(w) == 7 {
			w = append(w, "")
		}
		s := Symbol{
			Addr: addr,
			Size: size,
			Type: t,
			Name: w[7],
		}
		st = append(st, s)
	}
	return st, nil
}

func (st SymbolTable) Sort() {
	sort.Slice(st, func(i, j int) bool { return st[i].Addr < st[j].Addr })
}

func (st SymbolTable) Search(addr uint64) (Symbol, bool) {
	low := 0
	high := len(st) - 1
	if addr < st[low].Addr || st[high].Addr+st[high].Size <= addr {
		return Symbol{}, false
	}
	if st[high].Addr <= addr {
		return st[high], true
	}
	for 1 < high-low {
		mid := (low + high) / 2
		if st[mid].Addr <= addr {
			low = mid
		} else {
			high = mid
		}
	}
	if st[low].Addr <= addr && addr < st[low].Addr+st[low].Size {
		return st[low], true
	}
	return Symbol{}, false
}
