.PHONY: build start test build-lambda.zip dynamodb-local dynamodb-create_tables dynamodb-hydrate_tables dynamodb-delete_tables

build:
	sam build

start: build
	sam local start-api --docker-network host

test:
	go test -v ./...

build-lambda.zip:
	GOARCH=amd64 GOOS=linux go build -o main

lambda.zip: build-lambda.zip
	zip lambda.zip main

deps:
	go get -v -t -d

dynamodb-local:
	docker run -d -p 8000:8000 amazon/dynamodb-local:latest

dynamodb-create_tables:
	./scripts/create_tables.sh

dynamodb-hydrate_tables:
	./scripts/hydrate_tables.sh

dynamodb-delete_tables:
	./scripts/delete_tables.sh
