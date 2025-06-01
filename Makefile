BIN_FILE_PATH := ./bin/db-conn

build:
	go build -o $(BIN_FILE_PATH) .

run: build
	$(BIN_FILE_PATH) $(ARGS)
	
test:
	go test -v ./...
