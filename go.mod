module db-conn

go 1.22.6

require github.com/aws/aws-lambda-go v1.48.0

require (
	github.com/go-sql-driver/mysql v1.9.2
	github.com/lib/pq v1.10.9
	github.com/mattn/go-sqlite3 v1.14.28
)

require filippo.io/edwards25519 v1.1.0 // indirect
