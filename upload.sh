#!/usr/bin/env bash

API=$1
FILE=$2
FILENAME=$(basename $FILE)

echo "upload $FILENAME to $API"
(echo -n '{"imageBase64": "'; base64 $FILE; echo '"}') | curl -H "Content-Type: application/json" -d @- $API
