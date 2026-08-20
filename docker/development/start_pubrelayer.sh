#!/bin/bash

set -euo pipefail
# Read environment variables from $ENV_FILE
if [ -f "$ENV_FILE" ]; then
    source "$ENV_FILE"
else
    echo "Error: $ENV_FILE not found"
    exit 1
fi

# Start public relayer with debugging enabled
/go/bin/dlv exec --listen=:$GO_DEBUG_PORT --headless=true --log=true --accept-multiclient --api-version=2 --continue ./tmp/app/pubrelayer-app -- run --env $ENV_FILE