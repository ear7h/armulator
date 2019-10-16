package is

type Cond byte // actually 4 bits

const (
	CondEQ Cond = iota
	CondNE
	CondCS
	CondCC
	CondMI
	CondPL
	CondVS
	CondVC
	CondHI
	CondLS
	CondGE
	CondLT
	CondGT
	CondLE
	CondAL1
	CondAL2
)

func (c Cond) String() string {
	switch c {
	case CondEQ:
		return "EQ"
	case CondNE:
		return "NE"
	case CondCS:
		return "CS"
	case CondCC:
		return "CC"
	case CondMI:
		return "MI"
	case CondPL:
		return "PL"
	case CondVS:
		return "VS"
	case CondVC:
		return "VC"
	case CondHI:
		return "HI"
	case CondLS:
		return "LS"
	case CondGE:
		return "GE"
	case CondLT:
		return "LT"
	case CondGT:
		return "GT"
	case CondLE:
		return "LE"
	case CondAL1, CondAL2:
		return "AL"
	default:
		return "UNKNOWN CONDITION"
	}
}
