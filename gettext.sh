#!/usr/bin/env bash
API_ENDPOINT=$1
FILE=$2
FILENAME=$(basename $FILE)
EXTENSION="${FILENAME##*.}"
echo "Uploading file $FILENAME to $API_ENDPOINT"
ENCODED=$(base64 -i $FILE)

JSON="{\"imageBase64\": \"$ENCODED\", \"fileName\": \"$FILENAME\", \"extension\": \"$EXTENSION\"}"

curl -XPOST -d "$JSON" -H "Content-type: application/json" -v $API_ENDPOINT 
#curl -XPOST -d "$JSON" -H "Content-type: application/json" -v $API_ENDPOINT
(echo -n '{"imageBase64": "'; base64 $FILE; echo '", "fileName":"$FILENAME", "extension":"png"}') | curl -H "Content-Type: application/json" -d @- $API_ENDPOINT
