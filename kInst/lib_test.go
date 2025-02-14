package kInst_test

import (
	"testing"

	"github.com/clkbug/konata-tools"
	"github.com/clkbug/konata-tools/kInst"
)

func TestToInstructionSample1(t *testing.T) {
	cmds, err := konata.ParseFile("../kanata-sample-1.log")
	if err != nil {
		t.Fatal(err)
	}
	prog, err := kInst.ToProgram(cmds)
	if err != nil {
		t.Fatal(err)
	}

	if len(prog) != 2 {
		t.Errorf("len(prog) expected 2, but got %d", len(prog))
	}

	// ToDo: stage test
	expected := []kInst.Instruction{
		{Id: 0, Start: 216, Retire: 216 + 2, RetireType: konata.Successful},
		{Id: 1, Start: 217, Retire: 217 + 2, RetireType: konata.Flush},
	}

	for i, c := range prog {
		if c.Id != expected[i].Id ||
			c.Start != expected[i].Start ||
			c.Retire != expected[i].Retire ||
			c.RetireType != expected[i].RetireType {
			t.Errorf("prog[%d]: expected %v, but got %v", i, expected[i], c)
		}
	}
}
func TestToInstructionSample2(t *testing.T) {
	cmds, err := konata.ParseFile("../kanata-sample-2.log")
	if err != nil {
		t.Fatal(err)
	}
	prog, err := kInst.ToProgram(cmds)
	if err != nil {
		t.Fatal(err)
	}

	if len(prog) != 4041 {
		t.Errorf("len(prog) expected 4041, but got %d", len(prog))
	}

	// ToDo: stage test
	expected := []kInst.Instruction{
		{Id: 0, Start: 0, Retire: 24, RetireType: konata.Successful},
		{Id: 1, Start: 0, Retire: 15, RetireType: konata.Flush},
		{Id: 2, Start: 13, Retire: 16, RetireType: konata.Flush},
		{Id: 98, Start: 681, Retire: 705, RetireType: konata.Successful},
		{Id: 99, Start: 681, Retire: 706, RetireType: konata.Successful},
		{Id: 100, Start: 694, Retire: 709, RetireType: konata.Successful},
	}

	j := 0
	for i, c := range prog {
		if c.Id < expected[j].Id {
			continue
		}

		if c.Id != expected[j].Id ||
			c.Start != expected[j].Start ||
			c.Retire != expected[j].Retire ||
			c.RetireType != expected[j].RetireType {
			t.Errorf("prog[%d]: expected %v, but got %v", i, expected[j], c)
		}
		j++
		if len(expected) <= j {
			break
		}
	}
}
