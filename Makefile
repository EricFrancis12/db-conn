BIN_FILE_PATH := ./bin/db-conn

build:
	go build -o $(BIN_FILE_PATH) cmd/main.go

run: build
	$(BIN_FILE_PATH) $(ARGS)

codegen:
	go generate ./...

build_lambda:
	make codegen
	env GOOS=linux go build -o bootstrap ./lambda
	zip -r db-conn.zip bootstrap targets.txt
	
test:
	go test -v ./...
