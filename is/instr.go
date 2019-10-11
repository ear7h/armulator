package is

import (
	"github.com/ear7h/armulator/vm"
)


type Instr interface {
	Execute(*vm.State)
	String() string
}

type UnallocInstr string

func (instr UnallocInstr) Execute(state *vm.State) {
	return
}

func (instr UnallocInstr) String() string {
	return string(instr)
}
