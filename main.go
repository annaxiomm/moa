package main

import (
	"fmt"

	"github.com/annaxiomm/moa/cpu"
)

func main() {
	cpu := cpu.InitCpu()
	cpu.Load([]uint8{0x11, 0x10, 0x10, 0x01, 0x01})

	cpu.Run()
	fmt.Println(cpu.Registers)
}
