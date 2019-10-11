package is

import (
	"fmt"
	"strconv"

	"github.com/ear7h/armulator/vm"
)

type ADR struct {
	IsP bool
	Imm int32 // Actually 21 bits
	Label string
	Rd  byte
}

func (i ADR) String() string {
	mnem := "ADR"
	if i.IsP {
		mnem += "P"
	}

	// a label if known or rel addr
	op2 := i.Label
	if op2 == "" {
		op2 = immStr(int(i.Imm))
	}

	rd := regStr(true, i.Rd)
	return fmt.Sprintf("%s %s, %s",
		mnem, rd, op2)
}

func (i ADR) Execute(state *vm.State) {
	return
}

func DecodeADR(ei Encoded) ADR {
	immHi := ei.int32(5, 19)
	immLo := ei.int32(29, 2)

	return ADR{
		IsP: ei.bit(31),
		Imm:  immHi << 2 | immLo,
		Rd: ei.byte(0, 5),
	}
}

type ADDSUBImm struct {
	Is64Bit bool
	IsSub bool
	IsS bool
	Shift   byte
	Imm12   uint16
	Rn      byte
	Rd      byte
}

func (i ADDSUBImm) String() string {
	mnem := "ADD"
	if i.IsSub {
		mnem = "SUB"
	}

	if i.IsS {
		mnem += "S"
	}

	imm := immStr(int(i.Imm12))
	rn := regStr(i.Is64Bit, i.Rn)
	rd := regStr(i.Is64Bit, i.Rd)

	shift := ""
	switch i.Shift {
	case 0:
		// nothing
	case 1:
		shift = ", LSL #12"
	default:
		shift = ", RESERVED VALUE: " + strconv.Itoa(int(i.Shift))
	}

	return fmt.Sprintf(
		"%s %s, %s, %s%s",
		mnem, rd, rn, imm, shift)
}

func (i ADDSUBImm) Execute(state *vm.State) {
	return
}

func DecodeADDSUBImm(ei Encoded) ADDSUBImm {
	return ADDSUBImm{
		Is64Bit: ei.bit(31),
		IsSub:   ei.bit(30),
		IsS:     ei.bit(29),
		Shift:   ei.byte(22, 2),
		Imm12:   ei.uint16(10, 12),
		Rn:      ei.byte(5, 5),
		Rd:      ei.byte(0, 5),
	}
}

type LogicImm struct {
	Is64 bool
	Opc byte
	IsN bool
	Immr, Imms, Rn, Rd byte
}
