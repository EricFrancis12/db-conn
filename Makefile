BIN_FILE_PATH := ./bin/db-conn

build:
	go build -o $(BIN_FILE_PATH) cmd/main.go

run: build
	$(BIN_FILE_PATH) $(ARGS)

codegen:
	go generate ./...

LAMBDA_BIN_FILE_PATH := bootstrap
ZIP_FILE_PATH := db-conn.zip

build_lambda:
	make codegen
	env GOOS=linux go build -o $(LAMBDA_BIN_FILE_PATH) ./lambda
	zip -r $(ZIP_FILE_PATH) $(LAMBDA_BIN_FILE_PATH)
	
test:
	go test -v ./...
