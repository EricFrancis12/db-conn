package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	defaultDriverName  string        = "postgres"
	defaultConnTimeout time.Duration = time.Second * 5
	minConnTimeout     time.Duration = time.Second / 10
)

type Config struct {
	Ctx         context.Context
	Logger      Logger
	ConnStrs    []string
	DriverName  string
	ConnTimeout time.Duration
	Connect     ConnectFunc
}

func (c Config) Valid() error {
	problems := make(map[string]string)

	if c.Ctx == nil {
		problems["Context"] = "Context is nil"
	}

	if c.Logger == nil {
		problems["Logger"] = "Logger is nil"
	}

	if c.ConnStrs == nil {
		problems["ConnStrs"] = "ConnStrs is nil"
	} else if len(c.ConnStrs) == 0 {
		problems["ConnStrs"] = "ConnStrs has length of 0"
	}

	if c.DriverName == "" {
		problems["DriverName"] = "DriverName is missing"
	}

	if c.ConnTimeout < minConnTimeout {
		problems["ConnTimeout"] = fmt.Sprintf(
			"ConnTimeout '%s' is less than the minimum '%s'",
			c.ConnTimeout.String(),
			minConnTimeout.String(),
		)
	}

	if c.Connect == nil {
		problems["Connect"] = "Connect is nil"
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

	connStrs := []string{}
	for _, s := range c.ConnStrs {
		s = strings.TrimSpace(s)
		if s != "" {
			connStrs = append(connStrs, s)
		}
	}
	c.ConnStrs = connStrs

	if c.DriverName == "" {
		c.DriverName = defaultDriverName
	}

	if c.ConnTimeout < minConnTimeout {
		c.ConnTimeout = minConnTimeout
	}

	if c.Connect == nil {
		c.Connect = connect
	}

	return c.Valid()
}
