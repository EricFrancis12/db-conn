package main

import (
	"context"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	type Test struct {
		name          string
		shouldSucceed bool
		config        Config
	}

	var mockConnect ConnectFunc = func(
		ctx context.Context,
		driverName string,
		connStr string,
		timeout time.Duration,
	) error {
		return nil
	}

	tests := []Test{
		{
			name:          "Mock connect with no connection strings",
			shouldSucceed: false,
			config: Config{
				ConnStrs: []string{},
				Connect:  mockConnect,
			},
		},
		{
			name:          "Mock connect with connection strings",
			shouldSucceed: true,
			config: Config{
				ConnStrs: []string{"my connection string"},
				Connect:  mockConnect,
			},
		},
	}

	for _, test := range tests {
		success := run(&test.config)

		if test.shouldSucceed && !success {
			t.Logf("test '%s': expected success, but got failure", test.name)
			t.Fail()
		} else if !test.shouldSucceed && success {
			t.Logf("test '%s': expected failure, but got success", test.name)
			t.Fail()
		}
	}
}
