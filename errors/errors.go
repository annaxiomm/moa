package errors

import "fmt"

type ErrorCode int

const (
	WriteToROM ErrorCode = iota
	InvalidOpcode
	EndOfROM
	InvalidRegister
)

type ExecutionError struct {
	Type     ErrorCode
	Pc       uint16
	Address  uint16
	Opcode   uint8
	Register uint8
}

func (e *ExecutionError) Error() string {
	switch e.Type {
	case WriteToROM:
		return fmt.Sprintf("attempted to write to program ROM at address 0x%04X", e.Address)
	case InvalidOpcode:
		return fmt.Sprintf("invalid opcode 0x%02X", e.Opcode)
	case EndOfROM:
		return "program counter reached end of ROM"
	case InvalidRegister:
		return fmt.Sprintf("attempted to access invalid register %d", e.Register)
	default:
		return "unknown execution error"
	}
}

func PrintError(err error) {
	execErr, ok := err.(*ExecutionError)
	if !ok {
		fmt.Printf("moa: %v\n", err)
		return
	}

	fmt.Println("moa: execution error")
	fmt.Printf("    %v\n", execErr)
	fmt.Printf("    instruction: 0x%04X\n", execErr.Pc)
}
