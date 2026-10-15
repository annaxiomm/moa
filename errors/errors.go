package errors

import (
	"errors"
	"fmt"
)

type ErrorCode int

const (
	WriteToROM ErrorCode = iota
	InvalidOpcode
	EndOfROM
	InvalidRegister
	DivideByZero

	InvalidCharacter

	UnexpectedToken
	MissingComma
	MissingOperand
)

type ExecutionError struct {
	Type     ErrorCode
	Pc       uint16
	Address  uint16
	Opcode   uint8
	Register uint8
}

type LexError struct {
	Type ErrorCode
	Char string
	Line int
}

type ParseError struct {
	Type     ErrorCode
	Tok      string
	Expected string
}

func (e *LexError) Error() string {
	switch e.Type {
	case InvalidCharacter:
		return fmt.Sprintf("illegal character %s in assembly source", e.Char)
	default:
		return "unknown lexical error"
	}
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
	case DivideByZero:
		return "attempted to divide by zero"
	default:
		return "unknown execution error"
	}
}

func (e *ParseError) Error() string {
	switch e.Type {
	case UnexpectedToken:
		return fmt.Sprintf("unexpected token %s in source\n    expected %s", e.Tok, e.Expected)
	case InvalidRegister:
		return fmt.Sprintf("invalid register %s in source\n    registers should start with $R, and there are 8 registers (0-7)", e.Tok)
	case MissingComma:
		return "missing comma in operand list\n    operands must be separated by commas"
	case MissingOperand:
		return "missing operand in operand list"
	default:
		return "unknown parsing error"
	}
}

func PrintError(err error) {
	var execErr *ExecutionError
	var lexErr *LexError
	var parseErr *ParseError

	switch {
	case errors.As(err, &execErr):
		fmt.Println("moa: execution error")
		fmt.Printf("    %v\n", execErr)
		fmt.Printf("    at instruction 0x%04X\n", execErr.Pc)
	case errors.As(err, &lexErr):
		fmt.Println("moa: lexical error")
		fmt.Printf("    %v\n", lexErr)
		fmt.Printf("    at line %d\n", lexErr.Line)
	case errors.As(err, &parseErr):
		fmt.Println("moa: parsing error")
		fmt.Printf("    %v\n", parseErr)

	}
}
