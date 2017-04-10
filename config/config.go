package config

import (
	"errors"
	"github.com/BurntSushi/toml"
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

func NewDemoConfig() *Config {
	cfg := Config{}
	cfg.LangeType = "golang"
	cfg.LeetcodeSession = "cookie"
	cfg.CsrfToken = "csrftoken"

	return &cfg
}

func (c *Config) Validate() error {
	if c.CsrfToken == "" || c.LeetcodeSession == "" {
		return errors.New("invalid config")
	}
	if !isValidLanguageType(c.LangeType) {
		return errors.New("invliad  language type")
	}
	return nil
}

func isValidLanguageType(language string) bool {
	switch language {
	case "golang", "java", "csharp", "cpp", "c", "javascript":
		return true
	}
	return false
}
