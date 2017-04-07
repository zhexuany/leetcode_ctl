package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/zhexuany/leetcode-ctl/cmd/generate"
	"github.com/zhexuany/leetcode-ctl/cmd/submit"
	"github.com/zhexuany/leetcode-ctl/config"
)

// These variables are populated via the Go linker.
var (
	version string
	commit  string
	branch  string
)

func init() {
	// If commit, branch, or build time are not set, make that clear.
	if version == "" {
		version = "unknown"
	}
	if commit == "" {
		commit = "unknown"
	}
	if branch == "" {
		branch = "unknown"
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	m := NewMain()
	if err := m.Run(os.Args[1:]...); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type Main struct {
	Logger *log.Logger
	Stdin  io.Writer
	Stdout io.Writer
	Stderr io.Writer
}

func NewMain() *Main {
	return &Main{
		Logger: log.New(os.Stderr, "[run] ", log.LstdFlags),
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}

func (m *Main) Run(args ...string) error {
	name, args := ParseCommandName(args)

	switch name {
	case "submit":
		cmd := submit.NewCommand()
		return cmd.Run(args...)
	case "generate":
		cmd := generate.NewCommand()
		return cmd.Run(args...)
	case "config":
		config.NewPrintConfigCommand().Run(args...)
	case "version":
		if err := NewVersionCommand().Run(args...); err != nil {
			return fmt.Errorf("version: %s", err)
		}

	case "":
	case "help":
		fmt.Println(usage)
	default:
		return fmt.Errorf(`unknown command %s"`+`Run 'leetcode-ctl help' for usage`+"\n\n", name)
	}

	return nil
}

const usage = `leetcode-ctl is a command line controller. Via this awesome tool, you can submit your answer to leetcode inside a terminal.

Usage: leetcode-ctl [[command]] [[arguments]]

The commands are:
    config               display the default configuration
    help                 display this help message
    generate             uses a snapshot of a data node to rebuild a cluster
    submit               submit the solution to leetcode judge
    version              displays the version
`

// VersionCommand represents the command executed by "influxd version".
type VersionCommand struct {
	Stdout io.Writer
	Stderr io.Writer
}

// NewVersionCommand return a new instance of VersionCommand.
func NewVersionCommand() *VersionCommand {
	return &VersionCommand{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}

// Run prints the current version and commit info.
func (cmd *VersionCommand) Run(args ...string) error {
	// Parse flags in case -h is specified.
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.Usage = func() { fmt.Fprintln(cmd.Stderr, versionUsage) }
	if err := fs.Parse(args); err != nil {
		return err
	}

	// Print version info.
	fmt.Fprintf(cmd.Stdout, "leetcode-ctl v%s (git: %s %s)\n", version, branch, commit)

	return nil
}

var versionUsage = `Displays the leetcode-ctl version, build branch and git commit hash.

Usage: leetcode-ctl version
`

func ParseCommandName(args []string) (string, []string) {
	// Retrieve command name as first argument
	var name string
	if len(args) > 0 {
		if !strings.HasPrefix(args[0], "-") {
			name = args[0]
		} else if args[0] == "-h" || args[0] == "-help" || args[0] == "--help" {
			name = "help"
		}
	}

	// If command is "help" and has an argument then rewrite args to use "-h"
	if name == "help" && len(args) > 2 && !strings.HasPrefix(args[1], "-") {
		return args[1], []string{"-h"}
	}

	// If a named command is specified then return it with its arguments.
	if name != "" {
		return name, args[1:]
	}
	return "", args
}
