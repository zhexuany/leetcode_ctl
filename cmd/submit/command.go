package submit

import (
	"flag"

	"fmt"
	"github.com/zhexuany/leetcode-ctl/client"
	"github.com/zhexuany/leetcode-ctl/config"
)

type Command struct {
}

func NewCommand() *Command {
	return &Command{}
}

func (cmd *Command) Run(args ...string) error {
	opts, err := cmd.ParseFlags(args...)
	if err != nil {
		return err
	}
	cfg, err := config.NewConfig("./default.toml")
	if err != nil {
		return err
	}
	c, err := client.NewClient(cfg)
	if err != nil {
	}
	if err := c.Submit(opts.filePath); err != nil {
		return err
	}

	return nil
}

// ParseFlags parses the command line flags from args and returns an options set.
func (cmd *Command) ParseFlags(args ...string) (Options, error) {
	var opt Options
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.StringVar(&opt.filePath, "file", "", "specify the path of the file")
	fs.IntVar(&opt.problemID, "id", 0, "specify the problem id")
	fs.Usage = func() { fmt.Println(usage) }
	if err := fs.Parse(args); err != nil {
		return Options{}, err
	}
	return opt, nil
}

const usage = `Submit the answer to leetcode
Usage leetcode-ctl submit [flags]
    -config <path>
            Set the path to the configuration file.
    -file <path>
            The file path which contains the answer
    -id <int>
            The problem id that answer will submit to.
`

type Options struct {
	ConfigPath string
	filePath   string
	problemID  int
}
