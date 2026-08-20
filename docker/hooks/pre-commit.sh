#!/usr/bin/env bash

set -euo pipefail

# Save all output (stdout and stderr) to a file (pre-commit.log)
#exec &> pre-commit.log
echo "Starting pre-commit checks..."
echo "df -h: $(df -h)"
echo "free -m: $(free -m)"
echo "nproc: $(nproc)"

# Start development environment
./start_dev.sh --no-governance --no-otel --clean 6 & #&> pre-commit.log &
#./start_dev.sh --clean 6 &> output.log &
ENV_PID=$!
echo "ENV_PID=${ENV_PID}"

# Wait for relayer to be healthy
while true; do
    if RELAYER_HEALTHCHECK_URL=http://0.0.0.0:9000/healthcheck ./docker/development/relayer_healthcheck.sh; then
        break
    else
        # check if $ENV_PID process returned an error, or if it exited already
        if ! kill -0 $ENV_PID; then
            echo "Environment setup failed. Exiting..."
            exit 1
        fi
        sleep 5
    fi
done

# Temporarily disable exit-on-error for the test check
set +e
TESTS_FAILED=0
if ! ./compose.sh exec contracts bash -c "npm run test:e2e -- --bail"; then
    TESTS_FAILED=1
fi
set -e  # Re-enable exit-on-error

# Cleanup environment
./compose.sh down #&> shutdown.log

wait $ENV_PID

if [ "$TESTS_FAILED" -eq 0 ]; then
    echo "✅ All tests passed."
else
    echo "❌ Some tests failed."
fi

# Exit with appropriate status
exit $TESTS_FAILED
