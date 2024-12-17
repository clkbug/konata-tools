package rInst_test

import (
	"testing"

	"github.com/clkbug/konata-tools/rInst"
	"github.com/goccy/go-yaml"
)

func TestYamlUnMarshal(t *testing.T) {
	y := `pc: 0x1000
code: 0x617
roid: 10
flush_info: { type: "branch miss predict", pc: 0x2000, target: 0xFFFFF }
	`

	var info rInst.InstInfo
	err := yaml.Unmarshal([]byte(y), &info)
	if err != nil {
		t.Fatal(err)
	}
	if info.PC != 0x1000 {
		t.Fatalf("info: %v", info)
	}
	if info.Code != 0x617 || info.RobId != 10 || info.FlushInfo.T != "branch miss predict" {
		t.Fatal(info)
	}
}
