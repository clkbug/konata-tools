package instruction

import (
	"fmt"

	"github.com/clkbug/konata-tools"
)

type Instruction struct {
	Id, SimId, ThreadId int
	Start               int // absolute
	Retire              int // absolute, so latency = Retire - Start
	RetireId            int
	RetireType          konata.RetireType
	Label               [konata.LabelTypeCount]string
	Stage               []Stage
}

type Stage struct {
	Start int // relative to Instruction.Start
	End   int // opened. E command cycle not displayed
	Name  string
	Lane  int
}

type Program []Instruction

const NotRetired int = -1
const NotStageFinished int = -1

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
			prog = append(prog,
				Instruction{
					Id:       cmd.Id,
					SimId:    cmd.SimId,
					ThreadId: cmd.ThreadId,
					Retire:   NotRetired,
					Start:    c,
				})
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
				Start: c - prog[cmd.Id].Start,
				End:   NotStageFinished,
				Name:  cmd.StageName,
				Lane:  cmd.LaneId,
			})
		case konata.End:
			for i := range prog[cmd.Id].Stage {
				if prog[cmd.Id].Stage[i].End != NotStageFinished {
					continue
				}
				if prog[cmd.Id].Stage[i].Lane == cmd.LaneId && prog[cmd.Id].Stage[i].Name == cmd.StageName {
					prog[cmd.Id].Stage[i].End = c - prog[cmd.Id].Start
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
			prog[cmd.Id].RetireId = cmd.RetireId
			prog[cmd.Id].RetireType = cmd.RetireType
		case konata.WakeUp:
		}
	}

	// 正式にリタイアしていないものをリタイアさせる
	for i := range prog {
		if prog[i].Retire == NotRetired {
			prog[i].Retire = c
			prog[i].RetireId = prog[i].Id // ToDo: set apropriate retire id
			prog[i].RetireType = konata.Successful
		}
	}
	return prog, nil
}

func (p Program) ToCommand() ([]konata.Command, error) {
	var cmds []konata.Command
	start, end := 0, 0 // p[i] ~ p[j-1] in flight
	var retired []bool
	for c := 0; start < len(p); c++ {
		// I command
		for i := end; i < len(p) && p[i].Start == c; i++ {
			cmd, err := konata.MakeCommandInstruction(p[i].Id, p[i].SimId, p[i].ThreadId)
			if err != nil {
				return cmds, err
			}
			cmds = append(cmds, cmd)
			retired = append(retired, false)
			end++

			// L command
			for j, l := range p[i].Label {
				if l == "" {
					continue // Konata can't parse empty labels
				}
				if konata.LabelType(j) == konata.CurrentStage {
					continue // ToDo: support currentStage's label
				}
				cmd, err := konata.MakeCommandLabel(p[i].Id, konata.LabelType(j), l)
				if err != nil {
					return cmds, err
				}
				cmds = append(cmds, cmd)
			}
		}
		for i := start; i < end; i++ {
			// End Stage
			for _, s := range p[i].Stage {
				if p[i].Start+s.End == c {
					e, err := konata.MakeCommandEnd(p[i].Id, s.Lane, s.Name)
					if err != nil {
						return cmds, err
					}
					cmds = append(cmds, e)
				}
			}
			// Start Stage
			for _, s := range p[i].Stage {
				if p[i].Start+s.Start == c {
					e, err := konata.MakeCommandStage(p[i].Id, s.Lane, s.Name)
					if err != nil {
						return cmds, err
					}
					cmds = append(cmds, e)
				}
			}
			// Retire
			if p[i].Retire == c {
				r, err := konata.MakeCommandRetire(p[i].Id, p[i].RetireId, p[i].RetireType)
				if err != nil {
					return cmds, err
				}
				cmds = append(cmds, r)
				retired[i-start] = true
			}
		}
		// increment start
		for i := start; i-start < len(retired) && retired[i-start]; i++ {
			retired = retired[1:]
			start++
		}
		cmd, err := konata.MakeCommandCycle(1)
		if err != nil {
			return cmds, err
		}
		cmds = append(cmds, cmd)
	}
	return cmds, nil
}
