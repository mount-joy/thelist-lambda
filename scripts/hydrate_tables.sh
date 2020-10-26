#!/usr/bin/env bash
# Creates 5 lists and randomly adds up to 10 items to each one

ITEMS=("Apples" "Pears" "Oranges" "Limes" "Lemons" "Butter" "Salt" "Cereal" "Bread" "Milk" "Cheese")

for list in {1..5}; do
  ID=$(uuidgen | cut -c1-8)
  aws dynamodb put-item \
    --endpoint-url http://localhost:8000 \
    --region eu-west-2 \
    --table-name lists \
    --item "{ \"Id\": { \"S\": \"$ID\" }, \"Name\": { \"S\": \"List $list\" } }" \
    --condition-expression "attribute_not_exists(Id)"

  N=$(( $RANDOM % 10 ))
  echo "Adding $N items to list $list ($ID)"

  for item in $(seq 0 $N); do
    ITEM_ID=$(uuidgen | cut -c1-8)
    ITEM_NAME=${ITEMS[RANDOM%${#ITEMS[@]}]}

    aws dynamodb put-item \
      --endpoint-url http://localhost:8000 \
      --region eu-west-2 \
      --table-name items \
      --item "{ \"ListId\": { \"S\": \"$ID\" }, \"Id\": { \"S\": \"$ITEM_ID\" }, \"Item\": { \"S\": \"$ITEM_NAME\" } }"
  done
done
