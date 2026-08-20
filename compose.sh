#!/usr/bin/env bash

# Thin wrapper around `docker compose` that selects the right compose files for the
# current dev mode. Both modes share docker-compose.dev-local.yml as the base; remote
# mode layers docker-compose.remote.override.yml (which drops the local-axyl coupling
# in contracts.depends_on). The authoritative mode is the `mode=` field in the
# .deploy-ok marker written by start_dev.sh; fall back to the DEV_MODE hint in .A.env.
DEV_MODE=""
if [ -f .deploy-ok ]; then
    DEV_MODE=$(grep -o 'mode=[^ ]*' .deploy-ok | cut -d= -f2)
fi
if [ -z "$DEV_MODE" ] && [ -f .A.env ]; then
    DEV_MODE=$(grep -E '^[[:space:]]*#?[[:space:]]*DEV_MODE=' .A.env | head -n1 | sed -E 's/^[[:space:]]*#?[[:space:]]*DEV_MODE=[[:space:]]*//; s/[[:space:]]+$//')
fi

COMPOSE_ARGS=(-f docker-compose.dev-local.yml)
if [ "$DEV_MODE" = "remote" ]; then
    COMPOSE_ARGS+=(-f docker-compose.remote.override.yml)
fi

docker compose "${COMPOSE_ARGS[@]}" "$@"
