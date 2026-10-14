package emulator

import (
	"github.com/annaxiomm/moa/cpu"
	"github.com/annaxiomm/moa/screen"
)

func Run(program []uint8, nogui bool) error {
	cpu := cpu.InitCpu()
	cpu.Load(program)

	if !nogui {
		screen := screen.New(&cpu)
		if err := screen.Run(); err != nil {
			return err
		}
	} else {
		if err := cpu.Run(); err != nil {
			return err
		}
	}

	return nil
}
