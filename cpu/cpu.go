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
)

type Registers [8]uint8

type Cpu struct {
	Memory    memory.Memory
	Registers Registers
	Pc        uint16
	Halted    bool
}

func InitCpu() Cpu {
	return Cpu{Memory: memory.Memory{}, Registers: Registers{}, Pc: 0, Halted: false}
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
		byte := c.FetchByte()

		if err := c.LoadRegister(register, byte); err != nil {
			err.(*errors.ExecutionError).Pc = instructionPc
			return err
		}

	case STI:
		// Moa is little endian so the least significant byte comes firts
		lo := uint16(c.FetchByte())
		hi := uint16(c.FetchByte())
		value := c.FetchByte()

		location := (hi << 8) | lo

		if err := c.Memory.Write(location, value); err != nil {
			err.(*errors.ExecutionError).Pc = instructionPc
			return err
		}

	case LD:
		register := c.FetchByte()
		lo := uint16(c.FetchByte())
		hi := uint16(c.FetchByte())

		location := (hi << 8) | lo

		if err := c.LoadRegister(register, c.Memory.Read(location)); err != nil {
			err.(*errors.ExecutionError).Pc = instructionPc
			return err
		}
	case ST:
		lo := uint16(c.FetchByte())
		hi := uint16(c.FetchByte())
		register := c.FetchByte()

		location := (hi << 8) | lo

		value, err := c.ReadRegister(register)
		if err != nil {
			err.(*errors.ExecutionError).Pc = instructionPc
			return err
		}

		if err := c.Memory.Write(location, value); err != nil {
			err.(*errors.ExecutionError).Pc = instructionPc
			return err
		}
	default:
		fmt.Printf("0x%02X\n", c.Memory.Read(c.Pc))
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

func (c *Cpu) LoadRegister(register uint8, value uint8) error {
	if register > uint8(len(c.Registers)-1) {
		return &errors.ExecutionError{Type: errors.InvalidRegister, Register: value}
	}

	c.Registers[register] = value
	return nil
}

func (c *Cpu) ReadRegister(register uint8) (uint8, error) {
	if register > uint8(len(c.Registers)-1) {
		return 0, &errors.ExecutionError{Type: errors.InvalidRegister, Register: register}
	}

	value := c.Registers[register]
	return value, nil
}
