# DB Conn

A Go utility to test database connection strings in parallel, with configurable driver and input file.

## Features

- Reads a list of database connection strings from a file, one per line (see `targets.example.txt`).
- Connects to each database in parallel using the specified driver
- Reports success/failure for each connection and overall success rate

## Usage

```bash
make run ARGS="[flags]"
```

## Example Usage

```bash
make run ARGS="-f my_targets.txt -d mysql -t 2s"
```

## Supported Drivers (-d flag)

- postgres (default)
- mysql
- sqlite3

## AWS Lambda Deployment

You can deploy this utility as an AWS Lambda function. Follow these steps:

### 1. Create a targets.txt file

Create a file called `targets.txt` at the project root, and list your connection strings on it (one per line).

### 2. Build the Lambda Binary

Build a Linux executable named `bootstrap` (required name for Lambda custom runtime):

```bash
make build_lambda
```

This will:

- Build the binary for Linux (`GOOS=linux`)
- Output it as `bootstrap`
- Zip it as `db-conn.zip`

### 3. Upload to AWS Lambda

1. Go to the AWS Lambda console.
2. Create a new Lambda function (see image below).
3. Upload the `db-conn.zip` file as the function code.

<img src="https://github.com/user-attachments/assets/9bf6ef5e-7193-4e97-a591-04537aeaba19"/>
