# DB Conn

A Go utility to test database connection strings in parallel, with configurable driver, timeout, and input file.

## Features

- Reads a list of database connection strings from a file (one per line)
- Connects to each database in parallel using the specified driver
- Configurable connection timeout
- Reports success/failure for each connection and overall statistics

## Usage

```bash
make run ARGS="[flags]"
```

## Example Usage

```bash
make run ARGS="-f my_targets.txt -d mysql -t 2s"
```

## targets.txt Example

```txt
postgres://user:pass@localhost:5432/db1?sslmode=disable
postgres://user:pass@localhost:5432/db2?sslmode=disable
```
