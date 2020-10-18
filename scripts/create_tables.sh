#!/usr/bin/env bash

aws dynamodb create-table \
  --endpoint-url http://localhost:8000 \
  --region eu-west-2 \
  --table-name items \
  --attribute-definitions "AttributeName=ListId,AttributeType=S" "AttributeName=Id,AttributeType=S" \
  --key-schema "AttributeName=ListId,KeyType=HASH" "AttributeName=Id,KeyType=SORT" \
  --billing-mode PAY_PER_REQUEST

aws dynamodb create-table \
  --endpoint-url http://localhost:8000 \
  --region eu-west-2 \
  --table-name lists \
  --attribute-definitions "AttributeName=Id,AttributeType=S" \
  --key-schema "AttributeName=Id,KeyType=HASH" \
  --billing-mode PAY_PER_REQUEST
