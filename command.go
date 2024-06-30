package konata

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type CmdType int

const (
	CycleSet CmdType = iota // C= 	<CYCLE>
	Cycle                   // C	<CYCLE>
	Inst                    // I	<ID>	<SIM_ID>	<THREAD_ID>
	Label                   // L	<ID>	<TYPE>	<TEXT>
	Stage                   // S	<ID>	<LANE_ID>	<STAGE_NAME>
	End                     // E	<ID>	<LANE_ID>	<STAGE_NAME>
	Retire                  // R	<ID>	<RETIRE_ID>	<TYPE>
	WakeUp                  // W	<CONSUMER_ID>	<PRODUCER_ID>	<TYPE>
)

type LabelType int

const (
	LeftPane LabelType = iota
	MouseOver
	CurrentStage
	LabelTypeCount
)

type RetireType int

const (
	Successful RetireType = iota
	Flush
)

type DependencyType int

const (
	WakeUpDependency DependencyType = iota
)

type Command struct {
	T              CmdType
	Cycle          int            // CycleSet/Cycle
	Id             int            // Inst/Label/Stage/End/Retire/WakeUp
	SimId          int            // Inst
	ThreadId       int            // Inst
	LabelType      LabelType      // Label
	Text           string         // Label
	LaneId         int            // Stage
	StageName      string         // Stage
	RetireId       int            // Retire
	RetireType     RetireType     // Retire
	Consumer       int            // WakeUp
	Producer       int            // WakeUp
	DependencyType DependencyType // WakeUp
}

func ParseLine(l string) (Command, error) {
	w := strings.Split(l, "\t")
	switch w[0] {
	case "C=":
		cycle, err := strconv.Atoi(w[1])
		if err != nil {
			return Command{}, fmt.Errorf("failed to Atoi(%s) in '%s', %s", w[1], l, err)
		}
		return Command{
			T:     CycleSet,
			Cycle: cycle,
		}, nil
	case "C":
		cycle, err := strconv.Atoi(w[1])
		if err != nil {
			return Command{}, fmt.Errorf("failed to Atoi(%s) in '%s', %s", w[1], l, err)
		}
		return Command{
			T:     Cycle,
			Cycle: cycle,
		}, nil
	case "I":
		id, err := strconv.Atoi(w[1])
		if err != nil {
			return Command{}, fmt.Errorf("failed to Atoi(%s) in '%s', %s", w[1], l, err)
		}
		sid, err := strconv.Atoi(w[2])
		if err != nil {
			return Command{}, fmt.Errorf("failed to Atoi(%s) in '%s', %s", w[1], l, err)
		}
		tid, err := strconv.Atoi(w[3])
		if err != nil {
			return Command{}, fmt.Errorf("failed to Atoi(%s) in '%s', %s", w[1], l, err)
		}
		return Command{
			T:        Inst,
			Id:       id,
			SimId:    sid,
			ThreadId: tid,
		}, nil
	case "L":
		id, err := strconv.Atoi(w[1])
		if err != nil {
			return Command{}, fmt.Errorf("failed to Atoi(%s) in '%s', %s", w[1], l, err)
		}
		var lt LabelType
		switch w[2] {
		case "0":
			lt = LeftPane
		case "1":
			lt = MouseOver
		case "2":
			lt = CurrentStage
		default:
			return Command{}, fmt.Errorf("invalid label type(%s) in '%s'", w[2], l)
		}
		return Command{
			T:         Label,
			Id:        id,
			LabelType: lt,
			Text:      w[3],
		}, nil
	case "S":
		id, err := strconv.Atoi(w[1])
		if err != nil {
			return Command{}, fmt.Errorf("failed to Atoi(%s) in '%s', %s", w[1], l, err)
		}
		lid, err := strconv.Atoi(w[2])
		if err != nil {
			return Command{}, fmt.Errorf("failed to Atoi(%s) in '%s', %s", w[1], l, err)
		}
		return Command{
			T:         Stage,
			Id:        id,
			LaneId:    lid,
			StageName: w[3],
		}, nil
	case "R":
		id, err := strconv.Atoi(w[1])
		if err != nil {
			return Command{}, fmt.Errorf("failed to Atoi(%s) in '%s', %s", w[1], l, err)
		}
		rid, err := strconv.Atoi(w[2])
		if err != nil {
			return Command{}, fmt.Errorf("failed to Atoi(%s) in '%s', %s", w[1], l, err)
		}
		var rt RetireType
		switch w[3] {
		case "0":
			rt = Successful
		case "1":
			rt = Flush
		default:
			return Command{}, fmt.Errorf("unknown Retire Type(%s) in '%s'", w[3], l)
		}
		return Command{
			T:          Retire,
			Id:         id,
			RetireId:   rid,
			RetireType: rt,
		}, nil
	case "E":
		id, err := strconv.Atoi(w[1])
		if err != nil {
			return Command{}, fmt.Errorf("failed to Atoi(%s) in '%s', %s", w[1], l, err)
		}
		lid, err := strconv.Atoi(w[2])
		if err != nil {
			return Command{}, fmt.Errorf("failed to Atoi(%s) in '%s', %s", w[1], l, err)
		}
		return Command{
			T:         End,
			Id:        id,
			LaneId:    lid,
			StageName: w[3],
		}, nil
	case "W":
		cid, err := strconv.Atoi(w[1])
		if err != nil {
			return Command{}, fmt.Errorf("failed to Atoi(%s) in '%s', %s", w[1], l, err)
		}
		pid, err := strconv.Atoi(w[2])
		if err != nil {
			return Command{}, fmt.Errorf("failed to Atoi(%s) in '%s', %s", w[1], l, err)
		}
		var dt DependencyType
		switch w[3] {
		case "0":
			dt = WakeUpDependency
		default:
			return Command{}, fmt.Errorf("failed to parse dependency type(%s) in '%s'", w[3], l)
		}
		return Command{
			T:              WakeUp,
			Consumer:       cid,
			Producer:       pid,
			DependencyType: dt,
		}, nil
	default:
		return Command{}, fmt.Errorf("parse error '%s'", l)
	}
}

func ParseFile(filename string) ([]Command, error) {
	fp, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	buf := bufio.NewScanner(fp)
	var cmds []Command

	if buf.Scan() {
		line := buf.Text()
		if !strings.HasPrefix(line, "Kanata	0004") {
			return cmds, fmt.Errorf("failed to ParseFile(%s): invalid header '%s'", filename, line)
		}
	} else {
		return cmds, fmt.Errorf("failed to ParseFile(%s): no header", filename)
	}

	for buf.Scan() {
		line := buf.Text()
		cmd, err := ParseLine(line)
		if err != nil {
			return nil, err
		}
		cmds = append(cmds, cmd)
	}
	return cmds, nil
}

func (c *Command) String() string {
	switch c.T {
	case CycleSet:
		return fmt.Sprintf("C=\t%d", c.Cycle)
	case Cycle:
		return fmt.Sprintf("C\t%d", c.Cycle)
	case Inst:
		return fmt.Sprintf("I\t%d\t%d\t%d", c.Id, c.SimId, c.ThreadId)
	case Label:
		return fmt.Sprintf("L\t%d\t%d\t%s", c.Id, c.LabelType, c.Text)
	case Stage:
		return fmt.Sprintf("S\t%d\t%d\t%s", c.Id, c.LaneId, c.StageName)
	case End:
		return fmt.Sprintf("E\t%d\t%d\t%s", c.Id, c.LaneId, c.StageName)
	case Retire:
		return fmt.Sprintf("R\t%d\t%d\t%d", c.Id, c.RetireId, c.RetireType)
	case WakeUp:
		return fmt.Sprintf("W\t%d\t%d\t%d", c.Consumer, c.Producer, c.DependencyType)
	default:
		return fmt.Sprintf("UNKNOWN COMMAND TYPE: %#v", c)
	}
}
