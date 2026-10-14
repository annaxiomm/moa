package main

import (
	"fmt"
	"os"
	"path"

	"github.com/annaxiomm/moa/cli"
	"github.com/annaxiomm/moa/emulator"
)

func main() {

	args, err := cli.ParseCLI(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// STI address, value
	/*program := []uint8{
	0x10, 0x00, 0x30, 0x00,
	0x10, 0x01, 0x04, 0x00,
	0x10, 0x02, 0x05, 0x00,
	0x23, 0x01, 0x02,
	0x23, 0x00, 0x01,
	0x13, 0x00, 0xa6, 0x00,

	0x01,
	}*/

	if args.Name == "run" {
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		bytecode, err := os.ReadFile(path.Join(wd, args.Input))
		if err != nil {
			panic(err)
		}

		emulator.Run(bytecode, args.NoGUI)
	} else if args.Name == "asm" {
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		code, err := os.ReadFile(path.Join(wd, args.Input))
		if err != nil {
			panic(err)
		}

		bytecode, err := assembler.Assemble(code)
		if err != nil {
			panic(err)
		}

		err := os.WriteFile(path.Join(wd, args.Output), bytecode, 0644)
		if err != nil {
			panic(err)
		}
	}
}
