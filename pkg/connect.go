package pkg

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type ConnectFunc func(
	ctx context.Context,
	driverName string,
	connStr string,
	timeout time.Duration,
) error

func connect(
	ctx context.Context,
	driverName string,
	connStr string,
	timeout time.Duration,
) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	errChan := make(chan error)

	go func() {
		db, err := sql.Open(driverName, connStr)
		if err != nil {
			errChan <- err
			return
		}
		defer db.Close()
		errChan <- db.PingContext(ctx)
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("reached timeout of '%s'", timeout.String())
	case err := <-errChan:
		return err
	}
}
