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
	Label      [konata.LabelTypeCount]string
	Stage      []Stage
}

type Stage struct {
	Start int
	End   int // opened. E command cycle not displayed
	Name  string
	Lane  int
}

type Program []Instruction

func ToProgram(cmds []konata.Command) (Program, error) {
	var prog Program
	c := 0 // current cycle
cmdloop:
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
			prog[cmd.Id].Label[cmd.LabelType] = cmd.Text
		case konata.Stage:
			for i := range prog[cmd.Id].Stage {
				if prog[cmd.Id].Stage[i].Lane == cmd.LaneId {
					prog[cmd.Id].Stage[i].End = c
					break
				}
			}
			prog[cmd.Id].Stage = append(prog[cmd.Id].Stage, Stage{
				Start: c,
				End:   -1,
				Name:  cmd.StageName,
				Lane:  cmd.LaneId,
			})
		case konata.End:
			for i := range prog[cmd.Id].Stage {
				if prog[cmd.Id].Stage[i].End != -1 {
					continue
				}
				if prog[cmd.Id].Stage[i].Lane == cmd.LaneId && prog[cmd.Id].Stage[i].Name == cmd.StageName {
					prog[cmd.Id].Stage[i].End = c
					continue cmdloop
				}
				if prog[cmd.Id].Stage[i].Lane == cmd.LaneId || prog[cmd.Id].Stage[i].Name == cmd.StageName {
					return prog, fmt.Errorf("ToProgram(%v) (E command '%s' process):\n\tcurrent command stage/lane: %s/%d\n\tcurrent state stage/lane: %s/%d",
						cmd,
						cmd.String(),
						cmd.StageName, cmd.LaneId,
						prog[cmd.Id].Stage[i].Name, prog[cmd.Id].Stage[i].Lane,
					)
				}
			}
			return prog, fmt.Errorf("ToProgram(%v) (E command '%s' process):\n\tnot in stage?\n\tcurrent command stage/lane: %s/%d",
				cmd,
				cmd.String(),
				cmd.StageName, cmd.LaneId,
			)
		case konata.Retire:
			prog[cmd.Id].Retire = c
			prog[cmd.Id].RetireType = cmd.RetireType
		case konata.WakeUp:
		}
	}
	return prog, nil
}
