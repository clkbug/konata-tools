package instruction

import (
	"fmt"

	"github.com/clkbug/konata-tools"
)

type Instruction struct {
	Id         int
	Start      int
	Retire     int
	RetireType konata.RetireType
}

type Program []Instruction

func ToProgram(cmds []konata.Command) (Program, error) {
	var prog Program
	c := 0 // current cycle
	for _, cmd := range cmds {
		switch cmd.T {
		case konata.CycleSet:
			c = cmd.Cycle
		case konata.Cycle:
			c += cmd.Cycle
		case konata.Inst:
			if len(prog) != cmd.Id {
				return prog, fmt.Errorf("id? len(prog) = %d, but cmd.Id = %d", len(prog), cmd.Id)
			}
			prog = append(prog, Instruction{Id: cmd.Id, Start: c})
		case konata.Label:
		case konata.Stage:
		case konata.End:
		case konata.Retire:
			prog[cmd.Id].Retire = c
			prog[cmd.Id].RetireType = cmd.RetireType
		case konata.WakeUp:
		}
	}
	return prog, nil
}
