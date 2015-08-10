package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
)

var (
	commands     = make(map[string]Command)
	commandFlags = make(map[string]*flag.FlagSet)
)

var ErrUsage = UsageError("invalid command")

type UsageError string

func (e UsageError) Error() string {
	return fmt.Sprintf("Usage error: %s", string(e))
}

func Exit(code int) {
	os.Exit(code)
}

type Command interface {
	Usage()
	Run(args []string) error
}

type describer interface {
	Describe() string
}

type exampler interface {
	Examples() []string
}

// RegisterCommand registers a command to run for a given command name.
func RegisterCommand(name string, makeCmd func(Flags *flag.FlagSet) Command) {
	if _, dup := commands[name]; dup {
		log.Fatalf("duplicate command %q registered", name)
	}
	flags := flag.NewFlagSet(name+"options", flag.ContinueOnError)
	flags.Usage = func() {}

	commandFlags[name] = flags
	commands[name] = makeCmd(flags)
}

func Errorf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
}

func hasFlags(flags *flag.FlagSet) bool {
	any := false
	flags.VisitAll(func(*flag.Flag) {
		any = true
	})
	return any
}

func usage(msg string) {
	executable := filepath.Base(os.Args[0])
	if msg != "" {
		Errorf("Error: %v\n", msg)
	}
	Errorf(`
Usage: ` + executable + ` <command> [commandopts] [commandargs]

Commands:

`)

	var names []string
	for name, _ := range commands {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		cmd, _ := commands[name]
		if des, ok := cmd.(describer); ok {
			Errorf("  %-15s %s\n", name, des.Describe())
		}
	}

	Errorf("\nExamples:\n")
	for _, name := range names {
		cmd, _ := commands[name]
		if ex, ok := cmd.(exampler); ok {
			exs := ex.Examples()
			if len(exs) > 0 {
				Errorf("\n")
			}
			for _, example := range exs {
				Errorf("  %s %s %s\n", executable, name, example)
			}
		}
	}

	Errorf(`
For command-specific help:

  ` + executable + ` help <command>
`)
	//flag.PrintDefaults()
	Exit(1)
}

func help(name string) {
	executable := os.Args[0]
	cmd := commands[name]
	cmdFlags := commandFlags[name]
	cmdFlags.SetOutput(os.Stderr)
	if des, ok := cmd.(describer); ok {
		Errorf("%s\n", des.Describe())
	}
	Errorf("\n")
	cmd.Usage()
	if hasFlags(cmdFlags) {
		cmdFlags.PrintDefaults()
	}
	if ex, ok := cmd.(exampler); ok {
		Errorf("\nExamples:\n")
		for _, example := range ex.Examples() {
			Errorf("  %s %s %s\n", executable, name, example)
		}
	}
}

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		usage("No command given.")
	}

	name := args[0]
	cmd, ok := commands[name]
	if !ok {
		usage(fmt.Sprintf("Unknown command %q", name))
	}

	cmdFlags := commandFlags[name]
	cmdFlags.SetOutput(os.Stderr)
	err := cmdFlags.Parse(args[1:])
	if err != nil {
		err = ErrUsage
	} else {
		err = cmd.Run(cmdFlags.Args())
	}
	if e, isUsageErr := err.(UsageError); isUsageErr {
		Errorf("%s\n", e)
		cmd.Usage()
		Errorf("\nGlobal options:\n")
		flag.PrintDefaults()

		if hasFlags(cmdFlags) {
			Errorf("\nSpecific options for command %q:\n", name)
			cmdFlags.PrintDefaults()
		}
		Exit(1)
	}

	if err != nil {
		Errorf("Error: %v\n", err)
		Exit(2)
	}
}
