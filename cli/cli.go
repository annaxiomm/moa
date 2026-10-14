package cli

import (
	"fmt"
	"path/filepath"
	"strings"
)

type Command struct {
	Name   string
	Input  string
	Output string
	NoGUI  bool
	Help   bool
}

func ParseCLI(args []string) (Command, error) {
	if len(args) == 0 {
		return Command{}, fmt.Errorf("moa: expected a command")
	}

	cmd := Command{Name: args[0]}
	args = args[1:]

	switch cmd.Name {
	case "asm", "run":
		// Parse command arguments below.
	case "help", "version":
		return cmd, nil
	default:
		return Command{}, fmt.Errorf("moa: unknown command %q", cmd.Name)
	}

	var positional []string

	for i := 0; i < len(args); i++ {
		arg := args[i]

		switch {
		case arg == "-h" || arg == "--help":
			cmd.Help = true

		case cmd.Name == "asm" &&
			(arg == "-o" || arg == "--output"):
			if i+1 >= len(args) {
				return Command{}, fmt.Errorf(
					"%s requires a filename", arg,
				)
			}
			i++
			cmd.Output = args[i]

		case cmd.Name == "run" &&
			(arg == "-nogui" || arg == "--nogui"):
			cmd.NoGUI = true

		case strings.HasPrefix(arg, "-"):
			return Command{}, fmt.Errorf(
				"unknown option %q for %s", arg, cmd.Name,
			)

		default:
			positional = append(positional, arg)
		}
	}

	if cmd.Help {
		return cmd, nil
	}

	if len(positional) != 1 {
		return Command{}, fmt.Errorf(
			"moa: %s expects exactly one input file", cmd.Name,
		)
	}

	cmd.Input = positional[0]

	if cmd.Name == "asm" {
		if filepath.Ext(cmd.Input) != ".moa" {
			return Command{}, fmt.Errorf(
				"moa: expected a .moa source file, got %q", cmd.Input,
			)
		}

		if cmd.Output == "" {
			cmd.Output = strings.TrimSuffix(
				cmd.Input, filepath.Ext(cmd.Input),
			) + ".mxe"
		}
	}

	if cmd.Name == "run" &&
		filepath.Ext(cmd.Input) != ".mxe" {
		return Command{}, fmt.Errorf(
			"moa: expected a .mxe executable, got %q", cmd.Input,
		)
	}

	return cmd, nil
}
