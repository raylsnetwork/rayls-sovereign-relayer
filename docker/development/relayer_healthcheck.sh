#!/bin/bash

code=$(curl -o /dev/null -s -w "%{http_code}\n" $RELAYER_HEALTHCHECK_URL)

#echo "response code:$code"
if [ "$code" == "200" ]
then
  echo "relayer is healthy"
  exit 0;
else
  echo "waiting for relayer..."
  exit 1;
fi

