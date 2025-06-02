BIN_FILE_PATH := ./bin/db-conn

build:
	go build -o $(BIN_FILE_PATH) cmd/main.go

run: build
	$(BIN_FILE_PATH) $(ARGS)

build_lambda:
	env GOOS=linux go build -o bootstrap lambda/main.go
	zip -r db-conn.zip bootstrap targets.txt
	
test:
	go test -v ./...
