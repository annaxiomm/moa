package main

import (
	"fmt"
	"os"

	"github.com/annaxiomm/moa/cpu"
	"github.com/annaxiomm/moa/errors"
	"github.com/annaxiomm/moa/screen"
)

func main() {

	cpu := cpu.InitCpu()
	// STI address, value
	program := []uint8{
		0x11, 0x00, 0x66, 0xFF, // Pixel 0: white
		0x11, 0x01, 0x66, 0xCC, // Pixel 1: light grey
		0x11, 0x02, 0x66, 0x88, // Pixel 2: grey
		0x11, 0x03, 0x66, 0x44, // Pixel 3: dark grey
		0x01, // HLT
	}
	cpu.Load(program)

	screen := screen.New(&cpu)
	if err := screen.Run(); err != nil {
		errors.PrintError(err)
		os.Exit(1)
	}

	fmt.Println(cpu.Registers)
	fmt.Println(cpu.Memory.Read(0xffff))
}
