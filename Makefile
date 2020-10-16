.PHONY: build start test build-lambda.zip

build:
	sam build

start: build
	sam local start-api

test:
	go test -v ./...

build-lambda.zip:
	GOARCH=amd64 GOOS=linux go build -o main

lambda.zip: build-lambda.zip
	zip lambda.zip main
