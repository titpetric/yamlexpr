package main

import (
	"fmt"
	"os"
)

func main() {
	if err := start(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

// start handles the main command dispatch and execution.
func start(args []string) error {
	if len(args) < 1 {
		printUsage()
		return nil
	}

	// Parse command from first argument
	cmd := args[0]

	// Get remaining arguments
	cmdArgs := args[1:]

	var command Command

	switch cmd {
	case "process":
		command = NewProcessCommand()
	case "test":
		command = NewTestCommand()
	case "gen":
		command = NewGenCommand()
	case "help":
		if len(cmdArgs) > 0 {
			// Help for specific command
			subCmd := cmdArgs[0]
			switch subCmd {
			case "process":
				fmt.Println(NewProcessCommand().Help())
			case "test":
				fmt.Println(NewTestCommand().Help())
			case "gen":
				fmt.Println(NewGenCommand().Help())
			default:
				return fmt.Errorf("unknown command: %s", subCmd)
			}
		} else {
			printUsage()
		}
		return nil
	default:
		return fmt.Errorf("unknown command: %s", cmd)
	}

	// Run the command
	exitCode := command.Run(cmdArgs)
	if exitCode != 0 {
		return fmt.Errorf("command failed with exit code %d", exitCode)
	}
	return nil
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `yamlexpr - YAML expression evaluation and composition

Usage: yamlexpr <command> [options] [arguments]

Commands:
  process   Process and evaluate YAML files (default if no command given)
  test      Run fixture tests
  gen       Generate documentation from fixtures
  help      Show help for a command

Examples:
  yamlexpr process config.yaml
  yamlexpr test -dir testdata/fixtures-by-feature
  yamlexpr gen -feature for-loops
  yamlexpr help process

Use 'yamlexpr help <command>' for detailed help on a command.
`)
}
