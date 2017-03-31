package main

import (
	"errors"
	"github.com/BurntSushi/toml"
	"strings"
)

// Config defines your identity info and your prefered programming language
type Config struct {
	LeetcodeSession string `toml:"leetcode-session"`
	CsrfToken       string `toml:"csrf-token"`
	LangeType       string `toml:"lang-type"`
}

// NewConfig will decode the passed file and if such file is legal toml file
// a config will be created.
func NewConfig(path string) (*Config, error) {
	cfg := Config{}
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, err
	}

	return &cfg, cfg.Validate()
}

func (c *Config) Validate() error {
	if c.CsrfToken == "" || c.LeetcodeSession == "" {
		return errors.New("invalid config")
	}
	return nil
}
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
