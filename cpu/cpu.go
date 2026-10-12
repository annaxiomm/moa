package cpu

import (
	"fmt"

	"github.com/annaxiomm/moa/memory"
)

const (
	NOP = 0x00
	HLT = 0x01

	LDI = 0x10
	STI = 0x11
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
		c.Memory.Write(uint16(i), byte)
	}
}

func (c *Cpu) Run() {
	for !c.Halted {
		c.Step()

		if c.Pc >= 0x3FFF {
			c.Halted = true
			fmt.Println("[error] EOF reached (did you forget a HLT call?)")
		}
	}
}

func (c *Cpu) Step() {
	opcode := c.FetchByte()

	switch opcode {
	case NOP:
		// do nothing

	case HLT:
		c.Halted = true

	case LDI:
		register := c.FetchByte()
		byte := c.FetchByte()

		c.LoadRegister(register, byte)

	case STI:
		// Moa is little endian so the least significant byte comes firts
		lo := uint16(c.FetchByte())
		hi := uint16(c.FetchByte())
		value := c.FetchByte()

		location := (hi << 8) | lo

		fmt.Println(location, value)
	}
}

func (c *Cpu) FetchByte() uint8 {
	value := c.Memory.Read(c.Pc)
	c.Pc++
	return value
}

func (c *Cpu) LoadRegister(register uint8, value uint8) {
	if register > uint8(len(c.Registers)-1) {
		return
	}

	c.Registers[register] = value
}
