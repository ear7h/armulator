package vm

type State struct {
	Memory []byte
	Registers [30]uint64
}
