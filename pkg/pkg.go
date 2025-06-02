package pkg

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

const (
	defaultFilePath    string        = "targets.txt"
	defaultDriverName  string        = "postgres"
	defaultConnTimeout time.Duration = time.Second * 5
	minConnTimeout     time.Duration = time.Second / 10
)

func Run(c *Config) error {
	if err := c.Init(); err != nil {
		return err
	}

	commandLine := flag.NewFlagSet(c.Args[0], flag.ExitOnError)
	var (
		filePath       = commandLine.String("f", defaultFilePath, "Path to targets file. Each target should be on a new line")
		driverName     = commandLine.String("d", defaultDriverName, "SQL driver name")
		connTimeout    = commandLine.Duration("t", defaultConnTimeout, "DB connection timeout")
		listConnErrors = commandLine.Bool("errors", true, "Prints connection errors to terminal")
	)
	commandLine.Parse(c.Args[1:])

	b, err := c.ReadFile(*filePath)
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(b), "\n")

	connStrs := []string{}
	for _, line := range lines {
		s := strings.TrimSpace(line)

		// ignore empty lines & comments
		if s == "" || s[0] == '#' {
			continue
		}

		connStrs = append(connStrs, s)
	}

	connCount := len(connStrs)
	if connCount == 0 {
		return fmt.Errorf("file '%s' has no valid connection strings", *filePath)
	}

	if *connTimeout < minConnTimeout {
		return fmt.Errorf(
			"ConnTimeout '%s' is less than the minimum '%s'",
			connTimeout.String(),
			minConnTimeout.String(),
		)
	}

	errChan := make(chan error, connCount)
	var wg sync.WaitGroup

	c.Logger.Printf("Starting (%d) connection(s)\n", connCount)
	now := time.Now()

	for _, connStr := range connStrs {
		wg.Add(1)
		go func(cs string) {
			err := c.Connect(c.Ctx, *driverName, cs, *connTimeout)
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

	if *listConnErrors && errorCount > 0 {
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
