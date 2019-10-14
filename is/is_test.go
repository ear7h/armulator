package is

import (
	"strconv"
	"testing"
)

func TestString(t *testing.T) {
	type tcase struct {
		i   Encoded
		str string
	}

	fn := func(tc tcase) func(t *testing.T) {
		return func(t *testing.T) {
			if str := tc.i.Instr().String(); str != tc.str {
				t.Errorf("incorrect output\n%v\n\texpected\n%v",
					str,
					tc.str)
			}
		}
	}

	tcases := map[string][]tcase{
		"unalloc": []tcase{
			{
				i: 0,
				str: "unallocated instruction: 0",
			},
		},
		"adr":[]tcase {
			{ // 0
				i: 0x30000000,
				str: "ADR X0, #1",
			},
			{ // 1
				i: 0xb0000000,
				str: "ADRP X0, #1",
			},
			{ // 2
				i: 0xb000000f,
				str: "ADRP X15, #1",
			},
			{ // 3
				i: 0xb000001f,
				str: "ADRP SP, #1",
			},
			{ // 3
				i: 0x70ffffff,
				str: "ADR SP, #-1",
			},
			{ // 4
				i: 0xf0ffffff,
				str: "ADRP SP, #-1",
			},
		},
		"add sub": []tcase{
			// ADD
			{ // 0
				i: 0x11000000,
				str: "ADD W0, W0, #0",
			},
			{ // 1
				i: 0x11000021,
				str: "ADD W1, W1, #0",
			},
			{ // 2
				// 64 bits
				i: 0x91000021,
				str: "ADD X1, X1, #0",
			},
			{ // 3
				// stack pointer
				i: 0x1100003f,
				str: "ADD WSP, W1, #0",
			},
			{ // 4
				// stack pointer 64
				i: 0x9100003f,
				str: "ADD SP, X1, #0",
			},
			{ // 5
				// stack pointer 64 + 1
				i: 0x9100043f,
				str: "ADD SP, X1, #1",
			},
			{ // 6
				// stack pointer 64 + 6
				i: 0x9100183f,
				str: "ADD SP, X1, #6",
			},
			{ // 7
				// stack pointer 64 + max imm
				i: 0x913ffc3f,
				str: "ADD SP, X1, #4095",
			},
			{ // 8
				// stack pointer 64 + max imm
				i: 0x913ffc21,
				str: "ADD X1, X1, #4095",
			},
			// SUB
			{ // 9
				i: 0x51000000,
				str: "SUB W0, W0, #0",
			},
			{ // 10
				i: 0x51000021,
				str: "SUB W1, W1, #0",
			},
			{ // 11
				// 64 bits
				i: 0xd1000021,
				str: "SUB X1, X1, #0",
			},
			{ // 12
				// stack pointer
				i: 0x5100003f,
				str: "SUB WSP, W1, #0",
			},
			{
				// stack pointer 64
				i: 0xd100003f,
				str: "SUB SP, X1, #0",
			},
			{ // 13
				// stack pointer 64 + 1
				i: 0xd100043f,
				str: "SUB SP, X1, #1",
			},
			{ // 14
				// stack pointer 64 + 6
				i: 0xd100183f,
				str: "SUB SP, X1, #6",
			},
			{ // 15
				// stack pointer 64 + max imm
				i: 0xd13ffc3f,
				str: "SUB SP, X1, #4095",
			},
			{ // 16
				// stack pointer 64 + max imm
				i: 0xd13ffc21,
				str: "SUB X1, X1, #4095",
			},
			// [ADD|SUB]S
			{ // 17
				// stack pointer 64 + max imm
				i: 0xb13ffc21,
				str: "ADDS X1, X1, #4095",
			},
			{ // 18
				// stack pointer 64 + max imm
				i: 0xf13ffc21,
				str: "SUBS X1, X1, #4095",
			},
		},
		//
		// NOTE: the encoding of immediate values in logic
		// class operations is f-ing hard.
		//
		"logic" : []tcase{
			{
				//https://gist.github.com/dinfuehr/51a01ac58c0b23e4de9aac313ed6a06a#file-aarch64-logical-immediates-txt-L1303
				i: 0x12000000,
				str: "AND W0, W0, #0x1",
			},
			{
				// https://gist.github.com/dinfuehr/51a01ac58c0b23e4de9aac313ed6a06a#file-aarch64-logical-immediates-txt-L529
				i: 0x121A1800,
				str: "AND W0, W0, #0x1fc0",
			},
			{
				i: 0x321A1800,
				str: "ORR W0, W0, #0x1fc0",
			},
			{
				i: 0x521A1800,
				str: "EOR W0, W0, #0x1fc0",
			},
			{
				i: 0x721A1800,
				str: "ANDS W0, W0, #0x1fc0",
			},
			{
				i: 0x721A1800,
				str: "ANDS W0, W0, #0x1fc0",
			},
			{
				i: 0x721A1800,
				str: "ANDS W0, W0, #0x1fc0",
			},
			{
				i: 0x721A181F,
				str: "ANDS WSP, W0, #0x1fc0",
			},
			{
				i: 0x92400000,
				str: "AND X0, X0, #0x1",
			},
			{
				i: 0x924003e0,
				str: "AND X0, SP, #0x1",
			},
			{
				i: 0x927757E0,
				str: "AND X0, SP, #0x7ffffe00",
			},
		},
	}

	for k, v := range tcases {
		t.Run(k, func(t *testing.T) {
			for kk, vv := range v {
				t.Run(strconv.Itoa(kk), fn(vv))
			}
		})
	}

}

func TestBits(t *testing.T) {
	i := Encoded(1)
	if !i.bit(0) {
		panic(0)
	}

	i = 0
	if i.bit(0) {
		panic(1)
	}

	i = 0xf0
	if b := i.byte(4, 4); b != 0xf{
		t.Fatalf("%b", b)
	}

	i = 0xf0
	if b := i.int8(4, 4); b != -1 /* 0xff */ {
		t.Fatalf("%b", b)
	}
}
