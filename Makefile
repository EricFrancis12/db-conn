BIN_FILE_PATH := ./bin/db-conn

build:
	go build -o $(BIN_FILE_PATH) cmd/main.go

run: build
	$(BIN_FILE_PATH) $(ARGS)

LAMBDA_BIN_FILE_PATH := bootstrap
ZIP_FILE_PATH := db-conn.zip

codegen_lambda:
	go generate ./gen/lambda/...

build_lambda: codegen_lambda
	env GOOS=linux go build -o $(LAMBDA_BIN_FILE_PATH) ./lambda
	zip -r $(ZIP_FILE_PATH) $(LAMBDA_BIN_FILE_PATH)
	
test:
	go test -v ./pkg/...
