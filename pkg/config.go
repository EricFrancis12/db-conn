package pkg

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type Config struct {
	Ctx            context.Context
	Logger         Logger
	ConnStrs       []string
	DriverName     string
	ConnTimeout    time.Duration
	ListConnErrors bool
	Connect        ConnectFunc
}

func (c Config) Valid() error {
	problems := make(map[string]string)

	if c.Ctx == nil {
		problems["Context"] = "Context is nil"
	}

	if c.Logger == nil {
		problems["Logger"] = "Logger is nil"
	}

	if len(fmtLines(c.ConnStrs)) == 0 {
		problems["ConnStrs"] = "no valid connection strings"
	}

	if c.DriverName == "" {
		problems["DriverName"] = "DriverName is missing"
	}

	if c.ConnTimeout < MinConnTimeout {
		problems["ConnTimeout"] = fmt.Sprintf(
			"ConnTimeout '%s' is less than the minimum '%s'",
			c.ConnTimeout.String(),
			MinConnTimeout.String(),
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

	if c.DriverName == "" {
		c.DriverName = DefaultDriverName
	}

	if c.ConnTimeout == 0 {
		c.ConnTimeout = DefaultConnTimeout
	}

	if c.Connect == nil {
		c.Connect = connect
	}

	return c.Valid()
}
