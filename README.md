# thelist-lambda

## Requirements

* AWS CLI already configured with Administrator permission
* [Docker installed](https://www.docker.com/community-edition)
* [Golang](https://golang.org)
* SAM CLI - [Install the SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html)

## Running

`make build` - build the lambda ready for running locally.
`make start` - run locally on port 3000.
`make test` - runs all unit tests.
`make lambda.zip` - creates the lambda.zip file ready for deployment.
