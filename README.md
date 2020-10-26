# thelist-lambda

## Requirements

* AWS CLI already configured with Administrator permission
* [Docker installed](https://www.docker.com/community-edition)
* [Golang](https://golang.org)
* SAM CLI - [Install the SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html)

## Running

* `make build` - build the lambda ready for running locally.
* `make start` - run locally on port 3000.
* `make test` - runs all unit tests.
* `make lambda.zip` - creates the lambda.zip file ready for deployment.

### Running the database locally
To do so you will need [docker](https://www.docker.com/products/docker-desktop) and the [aws cli](https://docs.aws.amazon.com/cli/latest/userguide/install-cliv2.html ) installed.

* `make dynamodb-local` - run a local version of dynamodb on `http://localhost:8000`.
* `make dynamodb-create_tables` - create a local version of the tables used by the lambda.
* `make dynamodb-hydrate_tables` - creates a few lists and adds up to 10 items to each of them.
* `make dynamodb-delete_tables` - deletes the local tables.
