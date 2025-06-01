package main

import (
	"context"
	"errors"
	"fmt"
)

type Config struct {
	Ctx      context.Context
	Logger   Logger
	Args     []string
	Connect  ConnectFunc
	ReadFile func(string) ([]byte, error)
}

func (c Config) Valid() error {
	problems := make(map[string]string)

	if c.Ctx == nil {
		problems["Context"] = "Context is nil"
	}

	if c.Logger == nil {
		problems["Logger"] = "Logger is nil"
	}

	if len(c.Args) == 0 {
		problems["Args"] = "Args should be at least length 1 (got length 0)"
	}

	if c.Connect == nil {
		problems["Connect"] = "Connect is nil"
	}

	if c.ReadFile == nil {
		problems["ReadFile"] = "ReadFile is nil"
	}

	if len(problems) > 0 {
		msg := "invalid config:"
		for k, v := range problems {
			msg += fmt.Sprintf(" [%s] %s;", k, v)
		}
		return errors.New(msg)
	}

	return nil
}

func (c *Config) Init() error {
	if c.Ctx == nil {
		c.Ctx = context.Background()
	}

	if c.Logger == nil {
		c.Logger = mockLogger{}
	}

	if len(c.Args) == 0 {
		c.Args = []string{""}
	}

	if c.Connect == nil {
		c.Connect = connect
	}

	if c.ReadFile == nil {
		c.ReadFile = readFile
	}

	return c.Valid()
}
