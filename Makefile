BIN_FILE_PATH := ./bin/db-conn

build:
	go build -o $(BIN_FILE_PATH) .

run: build
	$(BIN_FILE_PATH) $(ARGS)

build_lambda:
	env GOOS=linux go build -ldflags="-X main.BuildMode=lambda" -o bootstrap .
	zip -r db-conn.zip bootstrap targets.txt
	
test:
	go test -v ./...
