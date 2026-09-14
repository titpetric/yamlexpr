package main

// Command is one subcommand of the tool: it runs with the arguments after
// its name and returns the exit code, and it describes itself for help.
type Command interface {
	Run(args []string) int
	Help() string
}
