package cpu

import (
	"fmt"

	"github.com/annaxiomm/moa/errors"
	"github.com/annaxiomm/moa/memory"
)

const (
	NOP = 0x00
	HLT = 0x01 // halt execution

	LDI = 0x10 // load instantaneous value into a register
	STI = 0x11 // store instantaneous value in memory
	LD  = 0x12 // load 1 byte from memory into a register
	ST  = 0x13 // store the value of a register into memory
	MOV = 0x14 // copy one register to another
	STR = 0x15 // stores the value of a register in the memory address stored in another register

	INC = 0x20 // increment a register
	DEC = 0x21 // decrement a register
	CMP = 0x22 // compare two registers, set zero flag to true if they are equal

	JMP = 0x30 // unconditional jump
	JNZ = 0x31 // jump if zero flag is not set
)

type Registers [8]uint16

type Cpu struct {
	Memory    memory.Memory
	Registers Registers
	Pc        uint16
	Halted    bool
	Zero      bool
}

func InitCpu() Cpu {
	return Cpu{Memory: memory.Memory{}, Registers: Registers{}, Pc: 0, Halted: false, Zero: false}
}

func (c *Cpu) Load(program []uint8) {
	for i, byte := range program {
		c.Memory.WriteUnsafe(uint16(i), byte)
	}
}

func (c *Cpu) Run() error {
	for !c.Halted {
		if err := c.Step(); err != nil {
			execErr := err.(*errors.ExecutionError)
			return execErr
		}

		if c.Pc >= 0x3FFF {
			c.Halted = true
			return &errors.ExecutionError{
				Type: errors.EndOfROM,
			}
		}
	}

	return nil
}

func (c *Cpu) Step() error {
	instructionPc := c.Pc
	opcode := c.FetchByte()

	switch opcode {
	case NOP:
		// do nothing

	case HLT:
		c.Halted = true

	case LDI:
		register := c.FetchByte()
		word := c.FetchWord()

		if err := c.LoadRegister(register, word); err != nil {
			err.(*errors.ExecutionError).Pc = instructionPc
			return err
		}

	case STI:
		location := c.FetchWord()
		value := c.FetchByte()

		if err := c.Memory.Write(location, value); err != nil {
			err.(*errors.ExecutionError).Pc = instructionPc
			return err
		}

	case LD:
		register := c.FetchByte()
		location := c.FetchWord()

		if err := c.LoadRegister(register, uint16(c.Memory.Read(location))); err != nil {
			err.(*errors.ExecutionError).Pc = instructionPc
			return err
		}
	case ST:
		location := c.FetchWord()
		register := c.FetchByte()

		value, err := c.ReadRegister(register)
		if err != nil {
			err.(*errors.ExecutionError).Pc = instructionPc
			return err
		}

		if err := c.Memory.Write(location, uint8(value)); err != nil {
			err.(*errors.ExecutionError).Pc = instructionPc
			return err
		}

		if location == 0xA600 {
			fmt.Print(string(c.Memory.Read(location)))
		}

	case MOV:
		register1 := c.FetchByte()
		register2 := c.FetchByte()

		value, err := c.ReadRegister(register2)
		if err != nil {
			err.(*errors.ExecutionError).Pc = instructionPc
			return err
		}

		if err := c.LoadRegister(register1, value); err != nil {
			err.(*errors.ExecutionError).Pc = instructionPc
			return err
		}

	case STR:
		register1 := c.FetchByte()
		register2 := c.FetchByte()

		location, err := c.ReadRegister(register1)
		if err != nil {
			err.(*errors.ExecutionError).Pc = instructionPc
			return err
		}

		value, err := c.ReadRegister(register2)
		if err != nil {
			err.(*errors.ExecutionError).Pc = instructionPc
			return err
		}

		if err := c.Memory.Write(location, uint8(value)); err != nil {
			err.(*errors.ExecutionError).Pc = instructionPc
			return err
		}

	case INC:
		register := c.FetchByte()
		current, err := c.ReadRegister(register)
		if err != nil {
			err.(*errors.ExecutionError).Pc = instructionPc
			return err
		}

		if err := c.LoadRegister(register, current+1); err != nil {
			err.(*errors.ExecutionError).Pc = instructionPc
			return err
		}

	case DEC:
		register := c.FetchByte()
		current, err := c.ReadRegister(register)
		if err != nil {
			err.(*errors.ExecutionError).Pc = instructionPc
			return err
		}

		if err := c.LoadRegister(register, current-1); err != nil {
			err.(*errors.ExecutionError).Pc = instructionPc
			return err
		}

	case CMP:
		register1, err := c.ReadRegister(c.FetchByte())
		if err != nil {
			return err
		}

		register2, err := c.ReadRegister(c.FetchByte())
		if err != nil {
			return err
		}

		if register1 == register2 {
			c.Zero = true
		} else {
			c.Zero = false
		}

	case JMP:
		lo := uint16(c.FetchByte())
		hi := uint16(c.FetchByte())
		c.Pc = uint16((hi << 8) | lo)

	case JNZ:
		lo := uint16(c.FetchByte())
		hi := uint16(c.FetchByte())

		if !c.Zero {
			c.Pc = uint16((hi << 8) | lo)
		}

	default:
		return &errors.ExecutionError{
			Type:   errors.InvalidOpcode,
			Opcode: opcode,
			Pc:     instructionPc,
		}
	}

	return nil
}

func (c *Cpu) FetchByte() uint8 {
	value := c.Memory.Read(c.Pc)
	c.Pc++
	return value
}

func (c *Cpu) FetchWord() uint16 {
	lo := uint16(c.FetchByte())
	hi := uint16(c.FetchByte())
	return (hi << 8) | lo
}

func (c *Cpu) LoadRegister(register uint8, value uint16) error {
	if register > uint8(len(c.Registers)-1) {
		return &errors.ExecutionError{Type: errors.InvalidRegister, Register: register}
	}

	c.Registers[register] = value
	return nil
}

func (c *Cpu) ReadRegister(register uint8) (uint16, error) {
	if register > uint8(len(c.Registers)-1) {
		return 0, &errors.ExecutionError{Type: errors.InvalidRegister, Register: register}
	}

	value := c.Registers[register]
	return value, nil
}
