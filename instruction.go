package konata

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type InstType int

const (
	CycleSet InstType = iota // C= 	<CYCLE>
	Cycle                    // C	<CYCLE>
	Inst                     // I	<ID>	<SIM_ID>	<THREAD_ID>
	Label                    // L	<ID>	<TYPE>	<TEXT>
	Stage                    // S	<ID>	<LANE_ID>	<STAGE_NAME>
	End                      // E	<ID>	<LANE_ID>	<STAGE_NAME>
	Retire                   // R	<ID>	<RETIRE_ID>	<TYPE>
	WakeUp                   // W	<CONSUMER_ID>	<PRODUCER_ID>	<TYPE>
)

type LabelType int

const (
	LeftPane LabelType = iota
	MouseOver
	CurrentStage
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

type Instruction struct {
	T              InstType
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

func ParseLine(l string) (Instruction, error) {
	w := strings.Split(l, "\t")
	switch w[0] {
	case "C=":
		cycle, err := strconv.Atoi(w[1])
		if err != nil {
			return Instruction{}, fmt.Errorf("failed to Atoi(%s) in '%s', %s", w[1], l, err)
		}
		return Instruction{
			T:     CycleSet,
			Cycle: cycle,
		}, nil
	case "C":
		cycle, err := strconv.Atoi(w[1])
		if err != nil {
			return Instruction{}, fmt.Errorf("failed to Atoi(%s) in '%s', %s", w[1], l, err)
		}
		return Instruction{
			T:     Cycle,
			Cycle: cycle,
		}, nil
	case "I":
		id, err := strconv.Atoi(w[1])
		if err != nil {
			return Instruction{}, fmt.Errorf("failed to Atoi(%s) in '%s', %s", w[1], l, err)
		}
		sid, err := strconv.Atoi(w[2])
		if err != nil {
			return Instruction{}, fmt.Errorf("failed to Atoi(%s) in '%s', %s", w[1], l, err)
		}
		tid, err := strconv.Atoi(w[3])
		if err != nil {
			return Instruction{}, fmt.Errorf("failed to Atoi(%s) in '%s', %s", w[1], l, err)
		}
		return Instruction{
			T:        Inst,
			Id:       id,
			SimId:    sid,
			ThreadId: tid,
		}, nil
	case "L":
		id, err := strconv.Atoi(w[1])
		if err != nil {
			return Instruction{}, fmt.Errorf("failed to Atoi(%s) in '%s', %s", w[1], l, err)
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
			return Instruction{}, fmt.Errorf("invalid label type(%s) in '%s'", w[2], l)
		}
		return Instruction{
			T:         Label,
			Id:        id,
			LabelType: lt,
			Text:      w[3],
		}, nil
	case "S":
		id, err := strconv.Atoi(w[1])
		if err != nil {
			return Instruction{}, fmt.Errorf("failed to Atoi(%s) in '%s', %s", w[1], l, err)
		}
		lid, err := strconv.Atoi(w[2])
		if err != nil {
			return Instruction{}, fmt.Errorf("failed to Atoi(%s) in '%s', %s", w[1], l, err)
		}
		return Instruction{
			T:         Stage,
			Id:        id,
			LaneId:    lid,
			StageName: w[3],
		}, nil
	case "R":
		id, err := strconv.Atoi(w[1])
		if err != nil {
			return Instruction{}, fmt.Errorf("failed to Atoi(%s) in '%s', %s", w[1], l, err)
		}
		rid, err := strconv.Atoi(w[2])
		if err != nil {
			return Instruction{}, fmt.Errorf("failed to Atoi(%s) in '%s', %s", w[1], l, err)
		}
		var rt RetireType
		switch w[3] {
		case "0":
			rt = Successful
		case "1":
			rt = Flush
		default:
			return Instruction{}, fmt.Errorf("unknown Retire Type(%s) in '%s'", w[3], l)
		}
		return Instruction{
			T:          Retire,
			Id:         id,
			RetireId:   rid,
			RetireType: rt,
		}, nil
	case "E":
		id, err := strconv.Atoi(w[1])
		if err != nil {
			return Instruction{}, fmt.Errorf("failed to Atoi(%s) in '%s', %s", w[1], l, err)
		}
		lid, err := strconv.Atoi(w[2])
		if err != nil {
			return Instruction{}, fmt.Errorf("failed to Atoi(%s) in '%s', %s", w[1], l, err)
		}
		return Instruction{
			T:         End,
			Id:        id,
			LaneId:    lid,
			StageName: w[3],
		}, nil
	case "W":
		cid, err := strconv.Atoi(w[1])
		if err != nil {
			return Instruction{}, fmt.Errorf("failed to Atoi(%s) in '%s', %s", w[1], l, err)
		}
		pid, err := strconv.Atoi(w[2])
		if err != nil {
			return Instruction{}, fmt.Errorf("failed to Atoi(%s) in '%s', %s", w[1], l, err)
		}
		var dt DependencyType
		switch w[3] {
		case "0":
			dt = WakeUpDependency
		default:
			return Instruction{}, fmt.Errorf("failed to parse dependency type(%s) in '%s'", w[3], l)
		}
		return Instruction{
			T:              WakeUp,
			Consumer:       cid,
			Producer:       pid,
			DependencyType: dt,
		}, nil
	default:
		return Instruction{}, fmt.Errorf("parse error '%s'", l)
	}
}

func ParseFile(filename string) ([]Instruction, error) {
	fp, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	buf := bufio.NewScanner(fp)
	var insts []Instruction

	if buf.Scan() {
		line := buf.Text()
		if !strings.HasPrefix(line, "Kanata	0004") {
			return insts, fmt.Errorf("failed to ParseFile(%s): invalid header '%s'", filename, line)
		}
	} else {
		return insts, fmt.Errorf("failed to ParseFile(%s): no header", filename)
	}

	for buf.Scan() {
		line := buf.Text()
		inst, err := ParseLine(line)
		if err != nil {
			return nil, err
		}
		insts = append(insts, inst)
	}
	return insts, nil
}
