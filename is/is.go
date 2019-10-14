// Pacakge is contains data relating to
// the aarch64 instruction set
package is

import "fmt"

// Encoded is an encoded Op and operands
type Encoded uint32

func (ei Encoded) bit(n int) bool {
	if n > 31 {
		panic("too many bits to shift")
	}

	return ei>>n&1 == 1
}

func (ei Encoded) int32(n, l int) int32 {
	// sign extended
	ret := int32(ei) << int32(32 - (n + l))
	ret = ret >> int32(32 - (n + l) + n)
	return ret
}

func (ei Encoded) int16(n, l int) int16 {
	return int16(ei.int32(n, l))
}

func (ei Encoded) int8(n, l int) int8 {
	return int8(ei.int32(n, l))
}


func (ei Encoded) uint32(n, l int) uint32 {
	ret := uint32(ei) << uint32(32 - (n + l))
	ret = ret >> uint32(32 - (n + l) + n)
	return ret
}

func (ei Encoded) uint16(n, l int) uint16 {
	return uint16(ei.uint32(n, l))
}

func (ei Encoded) byte(n, l int) byte {
	return byte(ei.uint32(n, l))
}


func (ei Encoded) Instr() Instr {

	op := ei.byte(25, 4)

	switch {
	case op&12 == 0: // 00xx
		return UnallocInstr(
			fmt.Sprintf("unallocated instruction: %x", uint(ei)))

	case op&14 == 8: // 100x
		return ei.dataProcImmInstr()

	case op&14 == 10: // 101x
		return ei.branchInstr()

	case op&5 == 4: // x1x0
		return ei.memInstr()

	case op&7 == 5: // x101
		return ei.dataProcRegInstr()

	case op&7 == 7: // x111
		return ei.dataProcFpInstr()

	default:
		panic("unknown op")
	}
}

//
//		op0 := bits(i, 25, 22)
//		op0 | decode group
//		====|===============
//		00x | pc-rel addr
//		01x | add/sub
//		100 | logical
//		101 | move wide
//		110 | bitfield
//		111 | extract
func (ei Encoded) dataProcImmInstr() Instr {
	op0 := ei.byte(23, 3)

	switch op0 {
	case 0, 1: // PC rel addressing
		return DecodeADR(ei)

	case 2, 3: // add sub
		return DecodeADDSUBImm(ei)

	case 4: // logical
		return DecodeLogicImm(ei)

	case 5: //move wide
		return nil
	case 7:
		return nil
	default:
		panic("unrechable")
	}
}

func (i Encoded) branchInstr() Instr {
	return nil
}

func (i Encoded) memInstr() Instr {
	return nil

}

func (i Encoded) dataProcRegInstr() Instr {
	return nil
}

func (i Encoded) dataProcFpInstr() Instr {
	return nil
}
