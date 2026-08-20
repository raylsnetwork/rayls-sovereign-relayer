#!/usr/bin/env bash
# Trigger Air hot-reload in all running dev containers.
# Air watches this file inside the container (synced via docker compose watch).
# Bumping its mtime is enough to trigger a recompile — equivalent to:
#   touch scripts/rebuild.sh
touch "$0"
