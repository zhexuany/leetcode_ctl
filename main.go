package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"time"
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
	//temporarily put client in main. TODO(zhexuany) move client to service
	client Client
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
	case "problems":
		ps := PostgresDB{}
		m.Logger.Println("Open database")
		ps.Open()
		m.Logger.Println("write all problems into database")
		ps.write()
	case "submit":
		cfg, err := NewConfig("./default.toml")
		if err != nil {
			return err
		}
		c, err := NewClient(cfg)
		if err != nil {
		}
		if err := c.Submit(""); err != nil {
			return err
		}
	case "generate":
	case "config":
	case "version":
		if err := NewVersionCommand().Run(args...); err != nil {
			return fmt.Errorf("version: %s", err)
		}

	case "":
	case "help":
		fmt.Println(usage)
		return nil
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
