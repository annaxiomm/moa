package main

import (
	"fmt"
	"os"

	"github.com/annaxiomm/moa/cpu"
	"github.com/annaxiomm/moa/errors"
)

func main() {
	cpu := cpu.InitCpu()
	cpu.Load([]uint8{0x10, 0x00, 0xff, 0x13, 0xff, 0xff, 0x00, 0x01})

	/* ldi $R0, 255
	 * st [0xffff], $R0
	 */

	if err := cpu.Run(); err != nil {
		errors.PrintError(err)
		os.Exit(1)
	}
	fmt.Println(cpu.Registers)
	fmt.Println(cpu.Memory.Read(0xffff))
}
