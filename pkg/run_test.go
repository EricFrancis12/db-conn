package pkg

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
				Logger:   logger{},
				ConnStrs: []string{},
				Connect:  mockConnect,
			},
		},
		{
			name:          "Mock connect with connection strings",
			shouldSucceed: true,
			config: Config{
				Logger:   logger{},
				ConnStrs: []string{"my connection string"},
				Connect:  mockConnect,
			},
		},
	}

	for _, test := range tests {
		err := Run(&test.config)

		if test.shouldSucceed && err != nil {
			t.Logf("test '%s': expected success, but got: %v", test.name, err)
			t.Fail()
		} else if !test.shouldSucceed && err == nil {
			t.Logf("test '%s': expected failure, but got nil", test.name)
			t.Fail()
		}
	}
}
