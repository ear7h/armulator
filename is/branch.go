package is

import (
	"fmt"
	"strconv"

	"github.com/ear7h/armulator/vm"
)

type BranchCondImm struct {
	Imm   int64 // 19 bits
	Label string
	Cond  Cond
}

func (i BranchCondImm) String() string {
	cond := i.Cond.String()

	imm := i.Label

	if len(imm) == 0 {
		imm = "#"+strconv.FormatInt(i.Imm, 10)
	}

	return fmt.Sprintf(
		"B%s %s",
		cond, imm)
}


func (i BranchCondImm) Execute (state *vm.State) {

}

func DecodeBranchCondImm(ei Encoded) Instr {

	if ei.bit(4) || ei.bit(24) {
		return UnallocInstr(
			fmt.Sprintf("unallocated conditional branch: %x", uint(ei)))
	}

	return BranchCondImm{
		Imm: ei.int64(5, 19) << 2,
		Cond: Cond(ei.byte(0, 4)),
	}
}

type ExceptionType byte

const (
	ExceptionSVC ExceptionType = iota
	ExceptionHVC
	ExceptionSMC
	ExceptionBRK
	ExceptionHLT
	ExceptionDCPS1
	ExceptionDCPS2
	ExceptionDCPS3
)

func (typ ExceptionType) String() string {
	switch typ {
	case ExceptionSVC:
		return "SVC"
	case ExceptionHVC:
		return "HVC"
	case ExceptionSMC:
		return "SMC"
	case ExceptionBRK:
		return "BRK"
	case ExceptionHLT:
		return "HLT"
	case ExceptionDCPS1:
		return "DCPS1"
	case ExceptionDCPS2:
		return "DCPS2"
	case ExceptionDCPS3:
		return "DCPS3"
	default:
		panic("unknown exception type")
	}
}

type ExceptionInstr struct {
	Type ExceptionType
	Imm uint16
}

func (i ExceptionInstr) Execute (vm *vm.State) {

}

func (i ExceptionInstr) String() string {
	return fmt.Sprintf("%s #%d", i.Type, i.Imm)
}

func DecodeExceptionGeneration(ei Encoded) Instr {
	opc := ei.byte(21, 3)
	op2 := ei.byte(2, 3)
	ll := ei.byte(0, 2)

	var exType ExceptionType

	switch {
	case opc == 0 && op2 == 0:
		switch ll {
		case 1:
			exType = ExceptionSVC
		case 2:
			exType = ExceptionHVC
		case 3:
			exType = ExceptionSMC
		default:
			panic("unreachable")
		}

	case opc == 1 && op2 == 0 && ll == 0:
		exType = ExceptionBRK

	case opc == 2 && op2 == 0 && ll == 0:
		exType = ExceptionHLT

	case opc == 5 && op2 == 0:
		switch ll {
		case 1:
			exType = ExceptionDCPS1
		case 2:
			exType = ExceptionDCPS2
		case 3:
			exType = ExceptionDCPS3
		default:
			panic("unreachable")
		}

	default:
		return UnallocInstr(
			fmt.Sprintf("unallocated instruction: %x", uint(ei)))
	}

	return ExceptionInstr {
		Type: exType,
		Imm: ei.uint16(5, 16),
	}
}

type BranchUncondRegType byte

const (
	BranchUnalloc BranchUncondRegType = iota
	BranchBR
	BranchBRAA
	BranchBRAAZ
	BranchBRAB
	BranchBRABZ

	BranchBLR
	BranchBLRAA
	BranchBLRAAZ
	BranchBLRAB
	BranchBLRABZ

	BranchRET
	BranchRETAA
	BranchRETAAZ
	BranchRETAB
	BranchRETABZ

	BranchERET
	BranchERETAA
	BranchERETAB

	BranchDRPS
)

func (typ BranchUncondRegType) String() string {
	switch typ {
	case BranchBR:
		return "BR"
	case BranchBRAA:
		return "BRAA"
	case BranchBRAAZ:
		return "BRAAZ"
	case BranchBRAB:
		return "BRAB"
	case BranchBRABZ:
		return "BRABZ"
	case BranchBLR:
		return "BLR"
	case BranchBLRAA:
		return "BLRAA"
	case BranchBLRAAZ:
		return "BLRAAZ"
	case BranchBLRAB:
		return "BLRAB"
	case BranchBLRABZ:
		return "BLRABZ"
	case BranchRET:
		return "RET"
	case BranchRETAA:
		return "RETAA"
	case BranchRETAAZ:
		return "RETAAZ"
	case BranchRETAB:
		return "RETAB"
	case BranchRETABZ:
		return "RETABZ"
	case BranchERET:
		return "ERET"
	case BranchERETAA:
		return "ERETAA"
	case BranchERETAB:
		return "ERETAB"
	case BranchDRPS:
		return "DRPS"
	default:
		return "unknown brnach type"
	}
}

type BranchUncondReg struct {
	Type BranchUncondRegType
	Rn byte
	Rm byte // only used in BRAA and BRAB
}

func (i BranchUncondReg) String() string {
	end := regStr(true, i.Rn)
	if i.Type == BranchBRAA || i.Type == BranchBRAB {
		end += ", " + regStr(true, i.Rm)
	}

	return fmt.Sprintf("%s %s", i.Type, end)
}

func (i BranchUncondReg) Execute(state *vm.State) {

}

func DecodeBranchUncondReg(ei Encoded) Instr {
	opc := ei.byte(21, 4)
	op2 := ei.byte(16, 5)
	op3 := ei.byte(10, 6)
	rn := ei.int8(5, 5)
	op4 := ei.int8(0, 5)


	decodeFields := []struct {
		opc, op2, op3 byte
		rn, op4 int8
		typ BranchUncondRegType
	}{
		{
			opc: 0,
			op2: 31,
			op3: 0,
			rn: -1,
			op4: 0,
			typ: BranchBR,
		},
		{
			opc: 0,
			op2: 31,
			op3: 2,
			rn: -1,
			op4: 31,
			typ: BranchBRAAZ,
		},
		{
			opc: 0,
			op2: 31,
			op3: 3,
			rn: -1,
			op4: 31,
			typ: BranchBRABZ,
		},
		{
			opc: 1,
			op2: 31,
			op3: 0,
			rn: -1,
			op4: 0,
			typ: BranchBLR,
		},
		{
			opc: 1,
			op2: 31,
			op3: 2,
			rn: -1,
			op4: 31,
			typ: BranchBLRAAZ,
		},
		{
			opc: 1,
			op2: 31,
			op3: 3,
			rn: -1,
			op4: 31,
			typ: BranchBLRABZ,
		},
		{
			opc: 2,
			op2: 31,
			op3: 0,
			rn: -1,
			op4: 0,
			typ: BranchRET,
		},
		{
			opc: 2,
			op2: 31,
			op3: 2,
			rn: 31,
			op4: 31,
			typ: BranchRETAA,
		},
		{
			opc: 2,
			op2: 31,
			op3: 3,
			rn: 31,
			op4: 31,
			typ: BranchRETAB,
		},
		{
			opc: 4,
			op2: 31,
			op3: 0,
			rn: 31,
			op4: 0,
			typ: BranchERET,
		},
		{
			opc: 4,
			op2: 31,
			op3: 2,
			rn: 31,
			op4: 31,
			typ: BranchERETAA,
		},
		{
			opc: 4,
			op2: 31,
			op3: 3,
			rn: 31,
			op4: 31,
			typ: BranchERETAA,
		},
		{
			opc: 5,
			op2: 31,
			op3: 0,
			rn: 31,
			op4: 0,
			typ: BranchDRPS,
		},
		{
			opc: 8,
			op2: 31,
			op3: 2,
			rn: -1,
			op4: -1,
			typ: BranchBRAA,
		},
		{
			opc: 8,
			op2: 31,
			op3: 3,
			rn: -1,
			op4: -1,
			typ: BranchBRAB,
		},
		{
			opc: 9,
			op2: 31,
			op3: 2,
			rn: -1,
			op4: -1,
			typ: BranchBLRAA,
		},
		{
			opc: 9,
			op2: 31,
			op3: 3,
			rn: -1,
			op4: -1,
			typ: BranchBLRAB,
		},
	}

	typ := BranchUnalloc

	for _, v := range decodeFields {
		if v.opc != opc {
			continue
		}

		if v.op2 != op2 {
			continue
		}

		if v.op3 != op3 {
			continue
		}

		if v.rn != -1 && v.rn != rn {
			continue
		}

		if v.op4 != -1 && v.op4 != op4 {
			continue
		}

		typ = v.typ
		break
	}

	if typ == BranchUnalloc {
		return UnallocInstr(
			fmt.Sprintf("unallocated instruction: %x", uint(ei)))
	}

	return BranchUncondReg{
		Type: typ,
		Rn: byte(rn),
		Rm: byte(op4),
	}

}

