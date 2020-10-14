.PHONY: build start test build-lambda.zip

build:
	sam build

start: build
	sam local start-api

test:
	go test -v ./src/

build-lambda.zip:
	GOARCH=amd64 GOOS=linux go build -o main ./src

lambda.zip: build-lambda.zip
	zip lambda.zip main
