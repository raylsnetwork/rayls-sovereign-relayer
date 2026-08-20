#!/bin/bash

set -euo pipefail
# Read environment variables from $ENV_FILE, to retrieve CTS URL
if [ -f "$ENV_FILE" ]; then
    source "$ENV_FILE"
else
    echo "Error: $ENV_FILE not found"
    exit 1
fi

# Extract hostname from gRPC address (e.g. "cts-a:8080" -> "cts-a")
CTS_HOST=$(echo "${CTS_GRPC_URL}" | sed 's|https\?://||' | cut -d: -f1)
CTS_HEALTH_URL="http://${CTS_HOST}:${CTS_HTTP_PORT}/health"

while true; do
    code=$(curl -s -o /dev/null -w "%{http_code}" "${CTS_HEALTH_URL}")
    if [ "$code" = "200" ] || [ "$code" = "401" ]; then
        echo "CTS is healthy (${CTS_HEALTH_URL})"
        /go/bin/dlv exec --listen=:$GO_DEBUG_PORT --headless=true --log=true --accept-multiclient --api-version=2 --continue ./tmp/app/relayer-app -- run --env $ENV_FILE
        break
    else
        echo "Waiting for CTS..."
        sleep 5
    fi
done
