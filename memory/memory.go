package memory

import (
	"github.com/annaxiomm/moa/errors"
)

const MemorySize = 1024 * 64 // 64 KiB of memory

type Memory [MemorySize]uint8

/* MEMORY REGIONS
 * 0x0000 - 0x3FFF: Program ROM (16KiB)
 * 0x4000 - 0x5FFF: General RAM (8KiB)
 * 0x6000 - 0x65FF: Palette RAM (1.5KiB - 256 colours * 3 bytes)
 * 0x6600 - 0xBFFF: Framebuffer (23,040B)
 * 0xC000 - 0xFFFF: Reserved for future use
 */

func (m *Memory) Read(addr uint16) uint8 {
	return m[addr]
}

// Primary write function for use in programs
func (m *Memory) Write(addr uint16, value uint8) error {
	// Programs shouldn't be able to write to their own ROM
	if addr <= 0x3FFF {
		return &errors.ExecutionError{
			Type:    errors.WriteToROM,
			Address: addr,
		}
	}

	m[addr] = value

	return nil
}

// Unsafe write function that allows you to write anywhere in memory - should only be used for program loading!
func (m *Memory) WriteUnsafe(addr uint16, value uint8) {
	m[addr] = value
}
