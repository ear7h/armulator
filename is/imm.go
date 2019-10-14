package is

import (
	"fmt"
	"math/bits"
	"strconv"

	"github.com/ear7h/armulator/vm"
)

type ADR struct {
	IsP   bool
	Imm   int32 // Actually 21 bits
	Label string
	Rd    byte
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
		Imm: immHi<<2 | immLo,
		Rd:  ei.byte(0, 5),
	}
}

type ADDSUBImm struct {
	Is64Bit bool
	IsSub   bool
	IsS     bool
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
	Is64Bit bool
	Opc     byte
	Imm     uint64
	Rn      byte
	Rd      byte
}

const (
	LogicImmOpcAnd = iota
	LogicImmOpcOrr
	LogicImmOpcEor
	LogicImmOpcAnds
)

func (i LogicImm) String() string {

	var opc string
	switch i.Opc {
	case LogicImmOpcAnd:
		opc = "AND"
	case LogicImmOpcOrr:
		opc = "ORR"
	case LogicImmOpcEor:
		opc = "EOR"
	case LogicImmOpcAnds:
		opc = "ANDS"
	default:
		panic("unknown opc")
	}

	rd := regStr(i.Is64Bit, i.Rd)
	rn := regStr(i.Is64Bit, i.Rn)

	return fmt.Sprintf(
		"%s %s, %s, #0x%x",
		opc, rd, rn, i.Imm)
}

func (i LogicImm) Execute(state *vm.State) {
	return
}

func DecodeLogicImm(ei Encoded) LogicImm {
	// no clue what's happening here but
	// llvm knows
	// https://llvm.org/doxygen/AArch64AddressingModes_8h_source.html#l00293
	// for reference when encoding
	// https://llvm.org/doxygen/AArch64AddressingModes_8h_source.html#l00213
	// https://dinfuehr.github.io/blog/encoding-of-immediate-values-on-aarch64/
	is64 := ei.bit(31)
	regSize := uint32(32)
	if is64 {
		regSize = uint32(64)
	}

	N := ei.uint32(22, 1)
	immr := ei.uint32(16, 6)
	imms := ei.uint32(10, 6)

	// assert
	if !(is64 || N == 0) {
		panic("undefined logical immediate encoding")
	}

	l := 31 - bits.LeadingZeros32((N<<6)|(^imms&0x3f))

	// assert
	if l < 0 {
		panic("undefined logical immediate encoding")
	}

	size := uint32(1 << l)
	R := immr & (size - 1)
	S := imms & (size - 1)

	// assert
	if S == size-1 {
		panic("undefined logical immediate encoding")
	}

	pattern := (uint64(1) << (S + 1)) - 1
	for i := uint32(0); i < R; i++ {
		// variable size ror
		pattern =  ((pattern & 1) << (size-1)) | (pattern >> 1)
	}

	for size != regSize {
		pattern |= (pattern << size)
		size *= 2
	}

	return LogicImm{
		Is64Bit: ei.bit(31),
		Opc:     ei.byte(29, 2),
		Imm:     pattern,
		Rn:      ei.byte(5, 5),
		Rd:      ei.byte(0, 5),
	}
}
