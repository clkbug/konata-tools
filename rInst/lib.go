package rInst

import (
	"fmt"
	"strings"

	"github.com/clkbug/konata-tools"
	"github.com/clkbug/konata-tools/kInst"
	"github.com/goccy/go-yaml"
)

type Instruction struct {
	kInst.Instruction
	InstInfo
}

type InstInfo struct {
	PC        uint64    `yaml:"pc"`
	Code      uint64    `yaml:"code"`
	Mnemonic  string    `yaml:"mnemonic"`
	RobId     int       `yaml:"roid"`
	LogRd     int       `yaml:"rd"`
	LogRs1    int       `yaml:"rs1"`
	LogRs2    int       `yaml:"rs2"`
	LogRs3    int       `yaml:"rs3"`
	PhysRd    int       `yaml:"rd_prid"`
	PhysRs1   int       `yaml:"rs1_prid"`
	PhysRs2   int       `yaml:"rs2_prid"`
	PhysRs3   int       `yaml:"rs3_prid"`
	FlushInfo FlushInfo `yaml:"flush_info"`
	TrapInfo  TrapInfo  `yaml:"trap_info"`
}

type FlushInfo struct {
	T      string `yaml:"type"`
	PC     uint64 `yaml:"pc"`
	Target uint64 `yaml:"target"`
}
type TrapInfo struct {
	T     string `yaml:"type"`
	RobId int    `yaml:"roid"`
	Tval  int    `yaml:"tval"`
}

type Program []Instruction

func ToProgram(cmds []konata.Command) (Program, error) {
	var prog Program
	p, err := kInst.ToProgram(cmds)
	if err != nil {
		return prog, err
	}
	for _, inst := range p {
		var info InstInfo
		label := strings.ReplaceAll(inst.Label[1], "\\n", "\n")
		err := yaml.Unmarshal([]byte(label), &info)
		if err != nil {
			return prog, err
		}
		fmt.Printf("%v\n", info)
		prog = append(prog, Instruction{
			inst,
			info,
		})
	}
	return prog, nil
}

func (p Program) ToCommand() ([]konata.Command, error) {
	var insts kInst.Program
	for _, i := range p {
		insts = append(insts, i.Instruction)
	}
	return insts.ToCommand()
}
