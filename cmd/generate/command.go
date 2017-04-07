package generate

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/zhexuany/leetcode-ctl/config"
	"github.com/zhexuany/leetcode-ctl/html"
)

type Command struct {
}

func NewCommand() *Command {
	return &Command{}
}

func (cmd *Command) Run(args ...string) error {
	opts, err := cmd.parseFlags(args...)
	if err != nil {
		return err
	}

	if opts.Validate() != nil {
		return opts.Validate()
	}

	// parse html and generate file
	html.GetJsonObjectFromLeetcode()
	fileName := html.QueryByID(opts.problemID)
	if fileName == "" {
		panic("problem id is not existed.")
	}

	cfg, err := config.NewConfig(opts.configPath)
	if err != nil {
		return err
	}

	ex := html.Extracter{}
	bs := []byte(ex.Find().Json().GetDefaultCode("golang"))
	fmt.Println("type", cfg.LangeType)
	ioutil.WriteFile(fileName+getFileExtenison(cfg.LangeType), bs, 0644)

	return nil
}

func getFileExtenison(language string) string {
	switch language {
	case "golang":
		return ".go"
	case "java":
		return ".java"
	case "csharp":
		return ".cs"
	case "cpp":
		return ".cc"
	case "c":
		return ".c"
	case "javascript":
		return ".js"
	}
	return ""
}
func (cmd *Command) parseFlags(args ...string) (Options, error) {
	var opt Options
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.StringVar(&opt.configPath, "config", "", "config file")
	fs.IntVar(&opt.problemID, "id", 0, "generate a problem file according to its id")
	fs.Usage = func() {
		fmt.Println(usage)
	}

	if err := fs.Parse(args); err != nil {
		return Options{}, err
	}

	return opt, nil
}

type Options struct {
	problemID  int
	configPath string
}

func (opt *Options) Validate() error {
	if opt.problemID == 0 {
		return errors.New("problem id is not set")
	}

	if opt.configPath == "" {
		return errors.New("config path is not set")
	}

	return nil
}

const usage = `
Generate the problem file according to problem's id
Usage leetcode-ctl generate [flags]
    -id <int>
            The problem id that answer will submit to.
`
