package pkg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

const (
	DefaultFilePath    string        = "targets.txt"
	DefaultDriverName  string        = "postgres"
	DefaultConnTimeout time.Duration = time.Second * 5
	MinConnTimeout     time.Duration = time.Second / 10
)

func Run(c *Config) error {
	if err := c.Init(); err != nil {
		return err
	}

	connCount := len(c.ConnStrs)
	errChan := make(chan error, connCount)
	var wg sync.WaitGroup

	c.Logger.Printf("Starting (%d) connection(s)\n", connCount)
	now := time.Now()

	for _, connStr := range c.ConnStrs {
		wg.Add(1)
		go func(cs string) {
			err := c.Connect(c.Ctx, c.DriverName, cs, c.ConnTimeout)
			if err != nil {
				errChan <- fmt.Errorf("error connecting to '%s': %v", connStr, err)
			}
			wg.Done()
		}(connStr)
	}

	wg.Wait()
	took := time.Since(now)
	c.Logger.Printf("Finished (%d) connection(s) in %s\n", connCount, took.String())
	close(errChan)

	// len() measures unread elements in a channel,
	// so we must take the len before ranging over
	errorCount := len(errChan)

	if c.ListConnErrors && errorCount > 0 {
		c.Logger.Println("--- START ERRORS ---")
		for err := range errChan {
			c.Logger.Printf("%v\n", err)
		}
		c.Logger.Println("---  END ERRORS  ---")
	}

	successCount := connCount - errorCount
	successRate := float64(successCount) / float64(connCount) * 100

	msg := fmt.Sprintf(
		"%d/%d connection attempt(s) succeeded (%.2f%% success rate)",
		successCount,
		connCount,
		successRate,
	)

	if errorCount > 0 {
		return errors.New(msg)
	}

	c.Logger.Println(msg)
	return nil
}

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

func readFile(filePath string) ([]byte, error) {
	stat, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("file '%s' not found", filePath)
	} else if err != nil {
		return nil, fmt.Errorf("error reading file '%s': %v", filePath, err)
	} else if stat.IsDir() {
		return nil, fmt.Errorf("'%s' needs to be a path to a file (found dir)", filePath)
	}

	b, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file '%s': %v", filePath, err)
	}

	return b, nil
}

func fmtLines(lines []string) []string {
	result := []string{}
	for _, line := range lines {
		s := strings.TrimSpace(line)

		// ignore empty lines & comments
		if s == "" || s[0] == '#' {
			continue
		}

		result = append(result, s)
	}
	return result
}

func ReadToConnStrs(filePath string) ([]string, error) {
	b, err := readFile(filePath)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(b), "\n")
	return fmtLines(lines), nil
}
