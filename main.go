package main

import (
	"os"

	"github.com/annaxiomm/moa/cpu"
	"github.com/annaxiomm/moa/errors"
	"github.com/annaxiomm/moa/screen"
)

func main() {

	cpu := cpu.InitCpu()
	// STI address, value
	program := []uint8{
		0x10, 0x00, 0x30, 0x00,
		0x10, 0x01, 0x04, 0x00,
		0x10, 0x02, 0x05, 0x00,
		0x23, 0x01, 0x02,
		0x23, 0x00, 0x01,
		0x13, 0x00, 0xa6, 0x00,

		0x01,
	}

	cpu.Load(program)

	screen := screen.New(&cpu)
	if err := screen.Run(); err != nil {
		errors.PrintError(err)
		os.Exit(1)
	}
}
