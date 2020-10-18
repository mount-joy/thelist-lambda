#!/usr/bin/env bash

aws dynamodb delete-table \
  --endpoint-url http://localhost:8000 \
  --region eu-west-2 \
  --table-name items

aws dynamodb delete-table \
  --endpoint-url http://localhost:8000 \
  --region eu-west-2 \
  --table-name lists
