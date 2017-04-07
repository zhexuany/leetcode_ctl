package generate

import (
	"flag"
	"fmt"
	"github.com/zhexuany/leetcode-ctl/html"
	"io/ioutil"
)

type Command struct {
}

func NewCommand() *Command {
	return &Command{}
}

func (cmd *Command) Run(args ...string) error {
	// opts, err := cmd.parseFlags(args...)
	// if err != nil {
	// 	return err
	// }

	// TODO form url first and then send get request to leetcode
	// Given a problem id, we query the problem name and form a url
	// ps := PostgresDB{}
	// m.Logger.Println("Open database")
	// ps.Open()
	// m.Logger.Println("write all problems into database")
	// ps.write()

	html.GetJsonObjectFromLeetcode()
	// parse html and generate file
	fileName := "./two-sum.go"
	// fileName += string(opts.problemID)
	ex := html.Extracter{}
	bs := []byte(ex.Find().Json().GetDefaultCode("golang"))
	ioutil.WriteFile(fileName, bs, 0644)
	return nil
}

func (cmd *Command) parseFlags(args ...string) (Options, error) {
	var opt Options
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.IntVar(&opt.problemID, "id", 0, "generate a problem file according to its id")
	fs.Usage = func() {
		fmt.Println(usage)
	}

	return opt, nil
}

type Options struct {
	problemID int
}

const usage = `
Generate the problem file according to problem's id
Usage leetcode-ctl generate [flags]
    -id <int>
            The problem id that answer will submit to.
`
