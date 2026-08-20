#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd -- "$(dirname -- "$0")" && pwd -P)"

resolve_path() {
    local target="$1"
    local parent_dir

    parent_dir="$(cd -- "$(dirname -- "$target")" && pwd -P)"
    printf '%s/%s\n' "$parent_dir" "$(basename -- "$target")"
}

# Set KEY=VALUE in an env file: replace the existing line or append it.
# Uses '|' as the sed delimiter so URL values (with '/') don't need escaping.
set_env_var() {
    local file="$1" key="$2" value="$3"
    if grep -qE "^${key}=" "$file"; then
        if [[ "$OSTYPE" == "darwin"* ]]; then
            sed -i "" "s|^${key}=.*|${key}=${value}|" "$file"
        else
            sed -i "s|^${key}=.*|${key}=${value}|" "$file"
        fi
    else
        printf '%s=%s\n' "$key" "$value" >> "$file"
    fi
}

on_error() {
    local exit_code=$?
    echo "Error: start_dev.sh failed at line $1 while running: $2" >&2
    exit "$exit_code"
}

trap 'on_error ${LINENO} "$BASH_COMMAND"' ERR

# The dev stack runs entirely locally: every service (privacy nodes, PNH, public
# chain, relayer, CTS, governance, ...) is spun up in Docker on this machine.

# Default values
DEV_MODE=local
NUM_PARTICIPANTS=2
GOVERNANCE_ENABLED=true
GOVERNANCE_API_ENABLED=true
EXPLORERS_ENABLED=true
OTEL_SDK_DISABLED=false
PUBLIC_CHAIN_ENABLED=true
HUB_ENABLED=true
BLOCKSCOUT_ENABLED=true
OPS_API_ENABLED=true
PLAYGROUND_ENABLED=true
# CLEAN_MODE: "none" = resume, "soft" = wipe relayer/CTS/governance DBs and rebuild dev services, "clean" = nuke everything
CLEAN_MODE=none
REBUILD_AXYL=0

usage() {
    echo "Usage: $0 [options] [num_participants]"
    echo ""
    echo "Options:"
    echo "  -sc, --soft-clean      Wipe relayer + pubrelayer databases. Keeps CTS, governance, contracts, and chain nodes (PNs/PNH/Public Chain)"
    echo "  -c, --clean            Tears down everything (volumes, contracts, nodes) and redeploys from scratch"
    echo "  --rebuild-axyl         Force rebuild of the local-axyl Docker image from ../axyl"
    echo "  --epoch-duration N     Override the chain's epoch_duration (in seconds). Default 86400 (24h)."
    echo "                         Lower this (e.g. 60) to trigger frequent epoch transitions."
    echo "                         Value is baked into chain genesis; requires --clean to change."
    echo "  --no-governance        Don't start Governance. No API, no Flagger, no Listener and no Postgres"
    echo "  --no-governance-api    Don't start Governance API or Auditor Explorer, but start Flagger, Listener and Postgres"
    echo "  --no-explorers         Don't start explorer services (e.g. Auditor Explorer). Useful for e2e tests"
    echo "  --no-otel              Don't start OTEL (OpenTelemetry)"
    echo "  --no-public-chain      Don't start public chain and public relayers"
    echo "  --no-hub               Don't start the Private Hub (PNH/Besu) or private relayers."
    echo "                         Runs CTS + public relayers against Privacy Node + Public Chain only."
    echo "  --no-blockscout        Don't start Blockscout explorers per participant"
    echo "  --no-ops-api           Don't start Ops API (per-participant API + worker) and Custody (HSM)"
    echo "  --no-playground        Don't start the Playground dApp (Next.js, port 3700)"
    echo "  -h, --help             Display this help message"
    echo ""
    echo "Arguments:"
    echo "  num_participants       Number of participants (2-6), default: 2"
    echo ""
    echo "Examples:"
    echo "  ./start_dev.sh              # Resume: restart services with existing state"
    echo "  ./start_dev.sh 4"
    echo "  ./start_dev.sh -sc 4        # Soft clean: wipe relayer + pubrelayer DBs, keep CTS/governance/contracts"
    echo "  ./start_dev.sh -c 4         # Clean: nuke everything, redeploy contracts"
    echo "  ./start_dev.sh -c --no-governance-api --no-otel 6"
}

# Parse options, including both short and long forms
while [ $# -gt 0 ]; do
    case $1 in
        -h|--help)
            usage
            exit 0
            ;;
        --no-governance)
            GOVERNANCE_ENABLED=false
            GOVERNANCE_API_ENABLED=false
            shift
            ;;
        --no-governance-api)
            GOVERNANCE_API_ENABLED=false
            shift
            ;;
        --no-explorers)
            EXPLORERS_ENABLED=false
            shift
            ;;
        --no-otel)
            OTEL_SDK_DISABLED=true
            shift
            ;;
        --no-public-chain)
            PUBLIC_CHAIN_ENABLED=false
            shift
            ;;
        --no-hub)
            HUB_ENABLED=false
            shift
            ;;
        --no-blockscout)
            BLOCKSCOUT_ENABLED=false
            shift
            ;;
        --no-ops-api)
            OPS_API_ENABLED=false
            shift
            ;;
        --no-playground)
            PLAYGROUND_ENABLED=false
            shift
            ;;
        -c|--clean)
            CLEAN_MODE=clean
            shift
            ;;
        -sc|--soft-clean)
            # Only upgrade to soft if not already clean (full nuke takes precedence)
            if [ "$CLEAN_MODE" != "clean" ]; then
                CLEAN_MODE=soft
            fi
            shift
            ;;
        --rebuild-axyl)
            REBUILD_AXYL=1
            shift
            ;;
        --epoch-duration)
            if [ -z "${2-}" ] || [[ "$2" == -* ]]; then
                echo "Error: --epoch-duration requires a numeric argument (seconds)" >&2
                exit 1
            fi
            # Exported so docker compose picks it up via ${RAYLS_AXYL_EPOCH_DURATION_SECS:-86400}
            # in the pn-*-genesis / pc-genesis service env blocks. Default in compose stays 86400.
            export RAYLS_AXYL_EPOCH_DURATION_SECS="$2"
            shift 2
            ;;
        *)
            # Check if the argument starts with '-'
            if [[ $1 == "-"* ]]; then
                echo "Error: Unknown option '$1'" >&2
                usage
                exit 1
            fi
            break
            ;;
    esac
done

# --no-hub and --no-public-chain are mutually exclusive: each selects a distinct
# compose override, and disabling both leaves no meaningful topology to start.
if [ "$HUB_ENABLED" = false ] && [ "$PUBLIC_CHAIN_ENABLED" = false ]; then
    echo "Error: --no-hub and --no-public-chain cannot be used together" >&2
    usage
    exit 1
fi

# Load OAuth credentials for ops-api dev (gitignored). When .env.oauth is
# absent, the GOOGLE_*/MICROSOFT_* vars stay unset and Compose resolves them
# to empty strings via ${VAR:-}, which keeps ops-api booting fine (the OAuth
# endpoints just respond 404 until creds are filled).
if [ -f ".env.oauth" ]; then
    set -a
    # shellcheck disable=SC1091
    source ".env.oauth"
    set +a
fi

# E2E seed identities (cmd/seed in ops-api). The seed creates a privacy-node operator
# and a normal user, mints their JWTs, and start_dev grants their on-chain roles and
# writes the tokens into the tests-automation .env. Defaults let seeding run out of the
# box; override in .env.oauth if you want stable, named identities.
OPS_OPERATOR_EMAIL="${OPS_OPERATOR_EMAIL:-e2e-operator@rayls.local}"
OPS_USER_EMAIL="${OPS_USER_EMAIL:-e2e-user@rayls.local}"

# Handle remaining arguments (num_participants)
if [ $# -gt 0 ]; then
    # Verify that the argument is a valid integer
    if ! [[ $1 =~ ^[0-9]+$ ]]; then
        echo "Error: '$1' is not a valid number" >&2
        echo "Number of participants must be an integer between 2 and 6" >&2
        usage
        exit 1
    fi
    NUM_PARTICIPANTS=$1
else
    NUM_PARTICIPANTS=2
fi

# Validate NUM_PARTICIPANTS is between 2 and 6
if [ $NUM_PARTICIPANTS -lt 2 ] || [ $NUM_PARTICIPANTS -gt 6 ]; then
    echo "Error: Number of participants must be between 2 and 6" >&2
    usage
    exit 1
fi

# Check if ../rayls-sovereign-contracts exists and error if not
if [ ! -d "../rayls-sovereign-contracts" ]; then
    RAYLS_CONTRACTS_DIR=$(resolve_path "$SCRIPT_DIR/../rayls-sovereign-contracts")
    echo "Error: rayls-sovereign-contracts directory not found at $RAYLS_CONTRACTS_DIR" >&2
    echo "       Try to git pull it over there." >&2
    exit 1
fi

# Check if ../governance-api exists and error if not
if [ "$GOVERNANCE_ENABLED" = "true" ] && [ ! -d "../rayls-sovereign-pnh-governance" ]; then
    GOVERNANCE_DIR=$(resolve_path "$SCRIPT_DIR/../rayls-sovereign-pnh-governance")
    echo "Error: rayls-sovereign-pnh-governance directory not found at $GOVERNANCE_DIR" >&2
    echo "       Try to git pull it over there." >&2
    exit 1
fi

# Check if ../rayls-sovereign-pnh-auditor-ui exists and error if not
if [ "$GOVERNANCE_API_ENABLED" = "true" ] && [ "$EXPLORERS_ENABLED" = "true" ] && [ ! -d "../rayls-sovereign-pnh-auditor-ui" ]; then
    AUDITOR_UI_DIR=$(resolve_path "$SCRIPT_DIR/../rayls-sovereign-pnh-auditor-ui")
    echo "Error: rayls-sovereign-pnh-auditor-ui directory not found at $AUDITOR_UI_DIR" >&2
    echo "       Try to git pull it over there." >&2
    exit 1
fi

# Check if ../rayls-privacy-blockscout exists when Blockscout is enabled
if [ "$BLOCKSCOUT_ENABLED" = "true" ] && [ ! -d "../rayls-privacy-blockscout" ]; then
    BLOCKSCOUT_DIR=$(resolve_path "$SCRIPT_DIR/../rayls-privacy-blockscout")
    echo "Error: rayls-privacy-blockscout directory not found at $BLOCKSCOUT_DIR" >&2
    echo "       git clone it there or re-run with --no-blockscout." >&2
    exit 1
fi

# Check if ../rayls-sovereign-ops-api and ../rayls-sovereign-custody-light exist when Ops API is enabled.
# Custody-light is a transitive dep of ops-api (HSM URL); both are bind-mounted at build time.
if [ "$OPS_API_ENABLED" = "true" ]; then
    if [ ! -d "../rayls-sovereign-ops-api" ]; then
        OPS_API_DIR=$(resolve_path "$SCRIPT_DIR/../rayls-sovereign-ops-api")
        echo "Error: rayls-sovereign-ops-api directory not found at $OPS_API_DIR" >&2
        echo "       git clone it there or re-run with --no-ops-api." >&2
        exit 1
    fi
    if [ ! -d "../rayls-sovereign-custody-light" ]; then
        CUSTODY_DIR=$(resolve_path "$SCRIPT_DIR/../rayls-sovereign-custody-light")
        echo "Error: rayls-sovereign-custody-light directory not found at $CUSTODY_DIR" >&2
        echo "       It is required by the ops-api (HSM dependency)." >&2
        echo "       git clone it there or re-run with --no-ops-api." >&2
        exit 1
    fi
fi

# Check if ../rayls-privacy-playground exists when Playground is enabled.
# The playground is a Next.js dApp bind-mounted into a node:20-alpine container
# (hot-reload via `npm run dev`) — see docker-compose.dev-local.yml `playground:` service.
if [ "$PLAYGROUND_ENABLED" = "true" ] && [ ! -d "../rayls-privacy-playground" ]; then
    PLAYGROUND_DIR=$(resolve_path "$SCRIPT_DIR/../rayls-privacy-playground")
    echo "Error: rayls-privacy-playground directory not found at $PLAYGROUND_DIR" >&2
    echo "       git clone it there or re-run with --no-playground." >&2
    exit 1
fi

# Axyl is required to run privacy nodes, the private hub, and the public chain locally.
if [ ! -d "../axyl" ]; then
    AXYL_DIR=$(resolve_path "$SCRIPT_DIR/../axyl")
    echo "Error: axyl directory not found at $AXYL_DIR . It is required to run privacy nodes, the private hub, and the public chain locally." >&2
    exit 1
fi

# When using DEV_MODE=local, ensure the local-axyl image is up-to-date.
# Rebuild only if: image is missing or ../axyl has newer commits than the image.
# --clean wipes environment state (volumes, contracts, nodes) but does NOT force an axyl rebuild,
# since building Rust from scratch is expensive. Use --rebuild-axyl to force it explicitly.
if [ "$DEV_MODE" = "local" ]; then
    AXYL_BUILD_NEEDED=0

    if ! docker image inspect local-axyl:latest &> /dev/null; then
        echo "local-axyl:latest image not found."
        AXYL_BUILD_NEEDED=1
    elif [ "$REBUILD_AXYL" = 1 ]; then
        echo "Axyl rebuild explicitly requested."
        AXYL_BUILD_NEEDED=1
    else
        # Compare the latest commit timestamp in ../axyl with the image creation time
        AXYL_COMMIT_TS=$(git -C "../axyl" log -1 --format=%ct 2>/dev/null || echo "0")
        IMAGE_CREATED=$(docker image inspect local-axyl:latest --format '{{.Created}}' 2>/dev/null || echo "")
        if [ -n "$IMAGE_CREATED" ] && [ "$AXYL_COMMIT_TS" != "0" ]; then
            IMAGE_TS=$(date -jf "%Y-%m-%dT%H:%M:%S" "$(echo "$IMAGE_CREATED" | cut -d. -f1)" "+%s" 2>/dev/null \
                    || date -d "$(echo "$IMAGE_CREATED" | cut -d. -f1)" "+%s" 2>/dev/null \
                    || echo "0")
            if [ "$AXYL_COMMIT_TS" -gt "$IMAGE_TS" ] 2>/dev/null; then
                echo "Axyl repo has commits newer than local-axyl:latest image — rebuilding."
                AXYL_BUILD_NEEDED=1
            fi
        fi
    fi

    if [ "$AXYL_BUILD_NEEDED" = 1 ]; then
        echo "Building local-axyl:latest from ../axyl..."
        # reth's build.rs reads VERGEN_GIT_SHA from the env and slices it (&sha[0..7]),
        # panicking on an empty value. CI passes this build-arg; mirror that locally by
        # stamping the axyl HEAD commit. Falls back to all-zeros (>= 8 chars) if unknown.
        export AXYL_GIT_SHA=$(git -C "../axyl" rev-parse HEAD 2>/dev/null || echo "0000000000000000000000000000000000000000")
        docker compose -f docker-compose.dev-local.yml build axyl-base
    fi
fi

# Generate local mTLS certs for the relayer<->CTS channel if missing.
# The script is idempotent (skips when cts/certs/ca.crt already exists).
"$SCRIPT_DIR/scripts/gen-dev-certs.sh"

# Check if docker and docker compose are installed on the computer
set +e
if ! command -v docker &> /dev/null; then
    echo "Error: Docker is not installed or not in PATH." >&2
    exit 1
fi

# Check if Docker Compose (as a plugin) is available
if ! docker compose version &> /dev/null; then
    echo "Error: Docker Compose is not installed or not available." >&2
    exit 1
fi
set -e

# Create an array with letters A-F up to NUM_PARTICIPANTS
PARTICIPANTS=()
for ((i=0; i<NUM_PARTICIPANTS; i++)); do
    PARTICIPANTS+=($(printf "\\x$(printf %02x $((65 + i)))"))
done

# --- Auto-escalation to clean ---
# Use the deploy marker as the single source of truth for what was previously deployed.
# If the participant count changed since last deployment, a full clean is required.
DEPLOY_MARKER=".deploy-ok"

if [ "$CLEAN_MODE" != "clean" ]; then
    if [ ! -f "$DEPLOY_MARKER" ]; then
        echo "No previous deployment found (missing $DEPLOY_MARKER). Escalating to clean..."
        CLEAN_MODE=clean
    else
        # Read previous deployment parameters from marker
        PREV_DEPLOY_PARTICIPANTS=$(grep -o 'participants=[^ ]*' "$DEPLOY_MARKER" | cut -d= -f2)

        PARTICIPANT_LIST_CHECK=$(IFS=,; echo "${PARTICIPANTS[*]}")

        if [ "$PREV_DEPLOY_PARTICIPANTS" != "$PARTICIPANT_LIST_CHECK" ]; then
            echo "Participants changed ($PREV_DEPLOY_PARTICIPANTS -> $PARTICIPANT_LIST_CHECK). Escalating to clean..."
            CLEAN_MODE=clean
        fi
    fi
fi

# --- Prepare env files ---
# On clean: reset all env files from templates.
# On soft-clean/resume: reuse existing env files as-is.
if [ "$CLEAN_MODE" = "clean" ]; then
    cp "../rayls-sovereign-contracts/.env.example-local" "../rayls-sovereign-contracts/.env"

    for name in "${PARTICIPANTS[@]}"; do
        PARTICIPANT_ENV_FILE=".${name}.env"
        cp "docker/development/local/$PARTICIPANT_ENV_FILE" "$PARTICIPANT_ENV_FILE"

        if [[ "$OSTYPE" == "darwin"* ]]; then
            sed -i "" "s/OTEL_SDK_DISABLED=.*/OTEL_SDK_DISABLED=$OTEL_SDK_DISABLED/" "$PARTICIPANT_ENV_FILE"
        else
            sed -i "s/OTEL_SDK_DISABLED=.*/OTEL_SDK_DISABLED=$OTEL_SDK_DISABLED/" "$PARTICIPANT_ENV_FILE"
        fi
    done
fi

# Build DB_NAMES from participant .env files
DB_NAMES=""
for participant in "${PARTICIPANTS[@]}"; do
    PARTICIPANT_ENV_FILE=".${participant}.env"

    # Extract database names from PostgreSQL connection strings
    # Format: postgres://user:pass@host:port/DBNAME?params
    PRIVATE_DB=$(grep -E '^PRIVATE_RELAYER_DATABASE_CONNECTIONSTRING=' "$PARTICIPANT_ENV_FILE" | sed -n 's|.*://[^/]*/\([^?]*\).*|\1|p')
    PUBLIC_DB=$(grep -E '^RAYLS_NODE_DATABASE_CONNECTIONSTRING=' "$PARTICIPANT_ENV_FILE" | sed -n 's|.*://[^/]*/\([^?]*\).*|\1|p')
    CTS_DB=$(grep -E '^CTS_DATABASE_CONNECTIONSTRING=' "$PARTICIPANT_ENV_FILE" | sed -n 's|.*://[^/]*/\([^?]*\).*|\1|p')

    # Blockscout shares this postgres instead of running its own (shared-db /
    # shared-stats-db). Pre-create the per-participant DBs so the lifecycle
    # matches the relayer/CTS DBs; the backend's create_and_migrate() and
    # stats' STATS__CREATE_DATABASE=true still take care of empty volumes.
    BLOCKSCOUT_DBS=""
    if [ "$BLOCKSCOUT_ENABLED" = "true" ]; then
        LOWERCASE_PARTICIPANT=$(echo "$participant" | tr '[:upper:]' '[:lower:]')
        BLOCKSCOUT_DBS="blockscout_pn-${LOWERCASE_PARTICIPANT},stats_pn-${LOWERCASE_PARTICIPANT},"
    fi

    # Ops API also shares this postgres (per-participant DB ops_api_<lc>).
    # golang-migrate runs automatically at boot, so the DB just needs to exist.
    OPS_API_DBS=""
    if [ "$OPS_API_ENABLED" = "true" ]; then
        LOWERCASE_PARTICIPANT=$(echo "$participant" | tr '[:upper:]' '[:lower:]')
        # Per-participant: ops-api DB + custody DB (each custody-<i> has its own).
        OPS_API_DBS="ops_api_${LOWERCASE_PARTICIPANT},raylzdb_${LOWERCASE_PARTICIPANT},"
    fi

    DB_NAMES+="${PRIVATE_DB},${PUBLIC_DB},${CTS_DB},${BLOCKSCOUT_DBS}${OPS_API_DBS}"
done

# Remove trailing comma
DB_NAMES=${DB_NAMES%,}

# Create config for governance
if [ "$GOVERNANCE_ENABLED" = "true" ] && ( [ ! -f "../rayls-sovereign-pnh-governance/.env" ] || [ "$CLEAN_MODE" = "clean" ] ); then
    cp ../rayls-sovereign-pnh-governance/config/.env.example ../rayls-sovereign-pnh-governance/.env
fi

# Add governance database to DB_NAMES if governance is enabled
if [ "$GOVERNANCE_ENABLED" = "true" ]; then
    GOVERNANCE_DB=$(grep -E '^DATABASE_CONNECTIONSTRING=' "../rayls-sovereign-pnh-governance/.env" | sed -n 's|.*://[^/]*/\([^?]*\).*|\1|p')
    if [ -n "$GOVERNANCE_DB" ]; then
        DB_NAMES="${DB_NAMES},${GOVERNANCE_DB}"
    fi
fi

# Generate BASE_SERVICES and DEV_SERVICES dynamically
PARTICIPANT_LIST=""
BASE_SERVICES="postgres mongodb contracts "
HUB_SERVICE=""
if [ "$HUB_ENABLED" = "true" ]; then
    HUB_SERVICE="private-hub proofs-api "
fi
BASE_SERVICES+="${HUB_SERVICE}"
if [ "$PUBLIC_CHAIN_ENABLED" = "true" ]; then
    BASE_SERVICES+="public-chain "
fi

# Attach otel
if [ "$OTEL_SDK_DISABLED" != "true" ]; then
    BASE_SERVICES+="otel "
fi

DEV_SERVICES=""
for participant in "${PARTICIPANTS[@]}"; do
    PARTICIPANT_LIST+="$participant,"
    LOWERCASE_PARTICIPANT=$(echo "$participant" | tr '[:upper:]' '[:lower:]')
    DEV_SERVICES+="cts-$LOWERCASE_PARTICIPANT "
    # The private relayer requires the hub (PNH_* config is required) — omit it
    # when --no-hub is set, keeping only CTS + the public relayer.
    if [ "$HUB_ENABLED" = "true" ]; then
        DEV_SERVICES+="relayer-$LOWERCASE_PARTICIPANT "
    fi
    if [ "$PUBLIC_CHAIN_ENABLED" = "true" ]; then
        DEV_SERVICES+="pubrelayer-$LOWERCASE_PARTICIPANT "
    fi

    BASE_SERVICES+="pn-$LOWERCASE_PARTICIPANT "
done

# Clean up any trailing whitespace
BASE_SERVICES="${BASE_SERVICES}"
DEV_SERVICES="${DEV_SERVICES}"

# Remove last comma from string. example: PARTICIPANT_LIST="A,B," becomes PARTICIPANT_LIST="A,B"
PARTICIPANT_LIST=${PARTICIPANT_LIST%,}

# echo "PARTICIPANTS=${PARTICIPANTS[@]}"
# echo "PARTICIPANT_LIST=$PARTICIPANT_LIST"
# echo "BASE_SERVICES=$BASE_SERVICES"
# echo "DEV_SERVICES=$DEV_SERVICES"

# --no-hub adds an override that drops private-hub from contracts.depends_on so the
# deploy no longer waits on the (absent) hub — see docker-compose.no-hub.override.yml.
COMPOSE_FILE="docker-compose.dev-local.yml"
if [ "$HUB_ENABLED" = "true" ]; then
    COMPOSE_OVERRIDE=""
else
    COMPOSE_OVERRIDE="docker-compose.no-hub.override.yml"
fi

# --- Podman (rootless) compatibility. Opt-in; Docker is the default. ---
# The team runs Docker, so that's the default behavior. Rootless-Podman users opt in
# with USE_PODMAN=true (e.g. `export USE_PODMAN=true` in your shell). We deliberately do
# NOT auto-switch from a `docker --version` string match: a Docker member must never get
# Podman tweaks (wrong UIDs) by accident. We only print a hint if Podman looks present
# but the flag wasn't set. The tweaks below are no-ops unless USE_PODMAN=true:
#
# 1) Podman's healthcheck exec can't resolve an image's NAMED user in /etc/passwd, so
#    those containers are wrongly marked unhealthy; forcing the numeric uid skips the
#    lookup. The values below MUST match the users baked into the source images:
#      private-hub: docker/development/local/Dockerfile.besu (FROM hyperledger/besu, USER besu = 1000:1000)
#      proofs-api:  ../rayls-sovereign-gnark-api/Dockerfile     (USER appuser = 999:999)
#    If an image changes its UID, the container goes unhealthy again — update it here.
#    Empty -> compose uses the image default via the `:-besu`/`:-appuser` sentinels.
# 2) Contracts deploy: run the container as root (CUSTOM_UID/GID=0) under Podman so it
#    can write the image-baked root-owned dirs (cache_forge) AND so its writes/chowns
#    land on the host user — rootless Podman maps container root <-> host uid by default,
#    so files in the bind-mounted repos stay owned by us instead of a subuid (100999).
#    On Docker the daemon is root, so the real uid/gid already keep host ownership.
USE_PODMAN="${USE_PODMAN:-false}"
if [ "$USE_PODMAN" != "true" ] && docker --version 2>/dev/null | grep -qi podman; then
    echo "Note: Podman detected but USE_PODMAN is not set — running in Docker-compatible mode." >&2
    echo "      On rootless Podman, re-run with: USE_PODMAN=true $0 $*" >&2
fi

# PRIVATE_HUB_USER/PROOFS_API_USER are constant across calls, so they're exported once
# for every `docker compose` invocation to pick up. CUSTOM_UID/GID differ between build
# (host uid) and runtime (0 on Podman), so they're passed inline per-call via DEPLOY_UID.
PRIVATE_HUB_USER=""
PROOFS_API_USER=""
DEPLOY_UID=$(id -u); DEPLOY_GID=$(id -g)
if [ "$USE_PODMAN" = "true" ]; then
    PRIVATE_HUB_USER="1000:1000"   # hyperledger/besu: USER besu
    PROOFS_API_USER="999:999"      # gnark-api: USER appuser
    DEPLOY_UID=0; DEPLOY_GID=0
fi
export PRIVATE_HUB_USER PROOFS_API_USER

# Helper: compose command with all required env vars
compose_env() {
    DB_NAMES=$DB_NAMES PUBLIC_CHAIN_ENABLED=$PUBLIC_CHAIN_ENABLED HUB_ENABLED=$HUB_ENABLED OTEL_SDK_DISABLED=$OTEL_SDK_DISABLED GOVERNANCE_ENABLED=$GOVERNANCE_ENABLED OPS_API_ENABLED=$OPS_API_ENABLED PLAYGROUND_ENABLED=$PLAYGROUND_ENABLED GOOGLE_CLIENT_ID=${GOOGLE_CLIENT_ID:-} GOOGLE_CLIENT_SECRET=${GOOGLE_CLIENT_SECRET:-} MICROSOFT_CLIENT_ID=${MICROSOFT_CLIENT_ID:-} MICROSOFT_CLIENT_SECRET=${MICROSOFT_CLIENT_SECRET:-} PNH_DEPLOYMENT_PROXY_REGISTRY=${PNH_DEPLOYMENT_PROXY_REGISTRY:-} PH_AUDITOR_ADDRESS=${PH_AUDITOR_ADDRESS:-} PN_A_AUDITOR_ADDRESS=${PN_A_AUDITOR_ADDRESS:-} PN_B_AUDITOR_ADDRESS=${PN_B_AUDITOR_ADDRESS:-} PN_C_AUDITOR_ADDRESS=${PN_C_AUDITOR_ADDRESS:-} PN_D_AUDITOR_ADDRESS=${PN_D_AUDITOR_ADDRESS:-} PN_E_AUDITOR_ADDRESS=${PN_E_AUDITOR_ADDRESS:-} PN_F_AUDITOR_ADDRESS=${PN_F_AUDITOR_ADDRESS:-} WALLETCONNECT_PROJECT_ID=${WALLETCONNECT_PROJECT_ID:-} PARTICIPANT_LIST=$PARTICIPANT_LIST DEV_MODE=$DEV_MODE CUSTOM_UID=$DEPLOY_UID CUSTOM_GID=$DEPLOY_GID docker compose -f $COMPOSE_FILE ${COMPOSE_OVERRIDE:+-f $COMPOSE_OVERRIDE} "$@"
}

# Build unique images via dedicated builder services (profiles: ["build"]).
# Maps participant services (e.g., cts-a, relayer-b) to their shared builder
# service (cts, relayer), deduplicates, and builds only the unique images.
build_images() {
    local builders=""
    for svc in $1; do
        case "$svc" in
            pn-*)        builders+="axyl-base " ;;
            cts-*)       builders+="cts " ;;
            relayer-*)   builders+="relayer " ;;
            pubrelayer-*) builders+="pubrelayer " ;;
            *)           builders+="$svc " ;;
        esac
    done
    # Deduplicate (e.g., "cts cts relayer relayer" -> "cts relayer")
    builders=$(echo "$builders" | tr ' ' '\n' | sort -u | tr '\n' ' ')
    # Build uses the host uid for the CUSTOM_UID arg (it only sets an ENV the runtime
    # overrides) to keep the contracts image cache stable; the root override (DEPLOY_UID)
    # is applied at runtime in compose_env, not here.
    DB_NAMES=$DB_NAMES PUBLIC_CHAIN_ENABLED=$PUBLIC_CHAIN_ENABLED HUB_ENABLED=$HUB_ENABLED OTEL_SDK_DISABLED=$OTEL_SDK_DISABLED GOVERNANCE_ENABLED=$GOVERNANCE_ENABLED OPS_API_ENABLED=$OPS_API_ENABLED PLAYGROUND_ENABLED=$PLAYGROUND_ENABLED GOOGLE_CLIENT_ID=${GOOGLE_CLIENT_ID:-} GOOGLE_CLIENT_SECRET=${GOOGLE_CLIENT_SECRET:-} MICROSOFT_CLIENT_ID=${MICROSOFT_CLIENT_ID:-} MICROSOFT_CLIENT_SECRET=${MICROSOFT_CLIENT_SECRET:-} PNH_DEPLOYMENT_PROXY_REGISTRY=${PNH_DEPLOYMENT_PROXY_REGISTRY:-} PH_AUDITOR_ADDRESS=${PH_AUDITOR_ADDRESS:-} PN_A_AUDITOR_ADDRESS=${PN_A_AUDITOR_ADDRESS:-} PN_B_AUDITOR_ADDRESS=${PN_B_AUDITOR_ADDRESS:-} PN_C_AUDITOR_ADDRESS=${PN_C_AUDITOR_ADDRESS:-} PN_D_AUDITOR_ADDRESS=${PN_D_AUDITOR_ADDRESS:-} PN_E_AUDITOR_ADDRESS=${PN_E_AUDITOR_ADDRESS:-} PN_F_AUDITOR_ADDRESS=${PN_F_AUDITOR_ADDRESS:-} WALLETCONNECT_PROJECT_ID=${WALLETCONNECT_PROJECT_ID:-} PARTICIPANT_LIST=$PARTICIPANT_LIST DEV_MODE=$DEV_MODE CUSTOM_UID=$(id -u) CUSTOM_GID=$(id -g) docker compose -f $COMPOSE_FILE ${COMPOSE_OVERRIDE:+-f $COMPOSE_OVERRIDE} --profile build build ${BUILD_TAGS:+--build-arg BUILD_TAGS=$BUILD_TAGS} $builders
}

# Bootstrap each ops-api admin (POST /admin/bootstrap) and grant the on-chain
# business roles on the matching Privacy Node. Runs in the background so the
# foreground `docker compose up --watch` can start; once the target ops-api
# becomes healthy the function makes its HTTP call.
#
# Skipped when OPS_API_ENABLED is not "true" or OPS_ADMIN_EMAIL is empty.
# Port mapping mirrors docker-compose.dev-local.yml: 9780 + index*100.
#
# Roles granted (in order):
#   - PRIVACY_NODE_OPERATOR:     required for /api/me to return 200 (login gate); signs createUser,
#                                approvals, and the PN registry-admin writes ops-api makes —
#                                updatePrivacyNodeStatus (approve), submitToHub, submitToPublicChain,
#                                and freezeOn{PrivacyNode,PublicChain}/unfreezeOn{…} — which
#                                PNTokenRegistryV1 now authorizes for PRIVACY_NODE_OPERATOR (see
#                                activate-business-roles-pn in rayls-sovereign-contracts). PN_TOKEN_REGISTRY_ADMIN
#                                is no longer needed here — the deploy task keeps it on initialOwner as break-glass.
#   - FACTORY_DEPLOYER:          required for POST /api/tokens (deploy/register).

# All are idempotent on-chain (grant-business-role hasRole-checks before grant).
bootstrap_ops_api() {
    set +e
    if [ "$OPS_API_ENABLED" != "true" ]; then return 0; fi
    if [ -z "${OPS_ADMIN_EMAIL:-}" ]; then
        echo "[bootstrap] OPS_ADMIN_EMAIL not set in .env.oauth — skipping ops-api admin bootstrap"
        return 0
    fi
    if ! command -v jq >/dev/null 2>&1; then
        echo "[bootstrap] jq not installed on host — install with 'apt install jq' or 'brew install jq' to enable auto-bootstrap"
        return 0
    fi

    local ADMIN_ROLES="PRIVACY_NODE_OPERATOR FACTORY_DEPLOYER"
    local index=0
    for participant in "${PARTICIPANTS[@]}"; do
        local lc upper port url tmp code addr elapsed
        lc=$(echo "$participant" | tr '[:upper:]' '[:lower:]')
        upper=$(echo "$participant" | tr '[:lower:]' '[:upper:]')
        port=$((9780 + index * 100))
        url="http://localhost:${port}/admin/bootstrap"

        echo "[bootstrap] Waiting for ops-api-${lc} on :${port} to become healthy…"
        elapsed=0
        while [ $elapsed -lt 1200 ]; do
            if curl -sf -o /dev/null --max-time 2 "http://localhost:${port}/health"; then
                break
            fi
            sleep 5
            elapsed=$((elapsed + 5))
        done

        if [ $elapsed -ge 1200 ]; then
            echo "[bootstrap] ! ops-api-${lc} not healthy after 20 min — skipping"
            index=$((index + 1))
            continue
        fi

        echo "[bootstrap] POST ${url} for ${OPS_ADMIN_EMAIL}"
        tmp=$(mktemp)
        local body
        body=$(jq -cn --arg email "$OPS_ADMIN_EMAIL" '{email: $email}')
        code=$(curl -sS -o "$tmp" -w "%{http_code}" \
            -X POST "$url" \
            -H 'Content-Type: application/json' \
            -d "$body" 2>/dev/null)
        [ -z "$code" ] && code="000"

        case "$code" in
            201)
                addr=$(jq -r '.address' "$tmp" 2>/dev/null)
                if [ -z "$addr" ] || [ "$addr" = "null" ]; then
                    echo "[bootstrap] ! ops-api-${lc} returned 201 with no address; skipping grant"
                else
                    echo "[bootstrap] → ops-api-${lc} admin wallet: ${addr}"
                    local role
                    for role in $ADMIN_ROLES; do
                        echo "[bootstrap] Granting ${role} on PN ${upper}…"
                        if ( cd ../rayls-sovereign-contracts &&
                             npx hardhat grant-business-role \
                                 --pn "$upper" \
                                 --role "$role" \
                                 --account "$addr" \
                                 --delay 0 ); then
                            echo "[bootstrap] ✓ ${role} granted on ops-api-${lc}"
                        else
                            echo "[bootstrap] ! grant-business-role ${role} failed for ${lc} — retry manually:"
                            echo "    (cd ../rayls-sovereign-contracts && npx hardhat grant-business-role --pn ${upper} --role ${role} --account ${addr} --delay 0)"
                        fi
                    done
                fi
                ;;
            409)
                echo "[bootstrap] → ops-api-${lc} already bootstrapped — skipping"
                ;;
            *)
                echo "[bootstrap] ! ops-api-${lc} returned HTTP ${code}:"
                cat "$tmp"
                echo
                ;;
        esac

        rm -f "$tmp"
        index=$((index + 1))
    done
}

# Seed the e2e operator + normal user (ops-api cmd/seed) on EVERY participant's
# ops-api, grant their on-chain business roles on the matching Privacy Node, and
# inject the per-node JWTs into the tests-automation .env.
#
# Per-participant (not once): ops-api resolves the authenticated user by the JWT's
# UUID and requires that user/wallet row to exist in that node's own DB, and cmd/seed
# mints a distinct random identity per database — so each node needs its own seed run
# and its own token. Mirrors bootstrap_ops_api's per-participant loop and port mapping
# (9780 + index*100).
#
# Roles: the operator gets the same roles as the admin (PRIVACY_NODE_OPERATOR +
# FACTORY_DEPLOYER + BANK_EMPLOYEE — the operator wallet signs every governance write, and
# addAddressPair / addToken are gated to BANK_EMPLOYEE on-chain). The PN registry-admin writes
# ops-api makes — updatePrivacyNodeStatus (approve), submitToHub, submitToPublicChain, and
# freezeOn{PrivacyNode,PublicChain}/unfreezeOn{…} — are now authorized for PRIVACY_NODE_OPERATOR by
# PNTokenRegistryV1 (see activate-business-roles-pn in rayls-sovereign-contracts), so PN_TOKEN_REGISTRY_ADMIN
# is no longer granted to the operator (it stays on initialOwner as deployer break-glass). The normal
# user gets BANK_EMPLOYEE (ops-api's default BLOCKCHAIN_OSA_USER_ROLE_NAME). Skipped when OPS_API_ENABLED != true.
#
# Resume-safe (runs on every start_dev.sh, not just --clean): cmd/seed's ensureUser is
# find-or-create keyed on email, so a re-run returns the existing user with the SAME
# wallet address (no duplicate, no mutation) and just mints a fresh JWT. Old tokens are
# stateless and stay valid until their TTL lapses; the .env is simply overwritten with
# the new one. grant-business-role is likewise idempotent (hasRole-checked).
seed_e2e_users() {
    set +e
    if [ "$OPS_API_ENABLED" != "true" ]; then return 0; fi

    local ADMIN_ROLES="PRIVACY_NODE_OPERATOR FACTORY_DEPLOYER BANK_EMPLOYEE"
    local USER_ROLES="BANK_EMPLOYEE"
    local automation_env="$AUTOMATION_PATH/.env"
    local index=0
    for participant in "${PARTICIPANTS[@]}"; do
        local lc upper port elapsed out seed_err operator_addr user_addr operator_key user_key role
        lc=$(echo "$participant" | tr '[:upper:]' '[:lower:]')
        upper=$(echo "$participant" | tr '[:lower:]' '[:upper:]')
        port=$((9780 + index * 100))

        # bootstrap_ops_api already polled this node's /health up to 20 min before
        # returning, so a short re-check suffices here — it only guards against the
        # node blipping during the grant-business-role calls in between, without
        # risking another full 20-min hang on an already-broken node.
        echo "[seed] Verifying ops-api-${lc} on :${port} is still healthy…"
        elapsed=0
        while [ $elapsed -lt 60 ]; do
            if curl -sf -o /dev/null --max-time 2 "http://localhost:${port}/health"; then
                break
            fi
            sleep 5
            elapsed=$((elapsed + 5))
        done
        if [ $elapsed -ge 60 ]; then
            echo "[seed] ! ops-api-${lc} unexpectedly unhealthy — skipping e2e seed for ${upper}"
            index=$((index + 1))
            continue
        fi

        echo "[seed] Running cmd/seed on ops-api-${lc} (operator=${OPS_OPERATOR_EMAIL}, user=${OPS_USER_EMAIL})…"
        # Capture stderr (compiler errors, panics, structured logs) to a temp file so it
        # can be surfaced on the error path below instead of being discarded.
        seed_err=$(mktemp)
        out=$(compose_env exec -T "ops-api-${lc}" go run ./cmd/seed \
            --operator-email "$OPS_OPERATOR_EMAIL" \
            --user-email "$OPS_USER_EMAIL" \
            --ttl 24h 2>"$seed_err")

        # Pull the four OPS_SERVICE_* lines out of the seed output (ignore any log noise
        # on the same stream). Values are split on the first '=' to keep JWTs intact.
        operator_addr=$(printf '%s\n' "$out" | grep -E '^OPS_SERVICE_OPERATOR_ADDRESS=' | head -1 | cut -d= -f2-)
        user_addr=$(printf '%s\n' "$out" | grep -E '^OPS_SERVICE_USER_ADDRESS=' | head -1 | cut -d= -f2-)
        operator_key=$(printf '%s\n' "$out" | grep -E '^OPS_SERVICE_OPERATOR_AUTH_KEY=' | head -1 | cut -d= -f2-)
        user_key=$(printf '%s\n' "$out" | grep -E '^OPS_SERVICE_USER_AUTH_KEY=' | head -1 | cut -d= -f2-)

        if [ -z "$operator_addr" ] || [ -z "$user_addr" ] || [ -z "$operator_key" ] || [ -z "$user_key" ]; then
            echo "[seed] ! cmd/seed on ops-api-${lc} did not return the expected OPS_SERVICE_* values — skipping ${upper}. Raw stdout:"
            printf '%s\n' "$out"
            if [ -s "$seed_err" ]; then
                echo "[seed] ! cmd/seed stderr:"
                cat "$seed_err"
            fi
            rm -f "$seed_err"
            index=$((index + 1))
            continue
        fi
        rm -f "$seed_err"
        echo "[seed] ops-api-${lc} → operator ${operator_addr}, user ${user_addr}"

        # Grant on-chain roles on this participant's PN. Idempotent (grant-business-role
        # hasRole-checks before granting). Operator gets ADMIN_ROLES, user gets USER_ROLES.
        for role in $ADMIN_ROLES; do
            echo "[seed] Granting ${role} to operator on PN ${upper}…"
            if ( cd ../rayls-sovereign-contracts &&
                 npx hardhat grant-business-role \
                     --pn "$upper" --role "$role" --account "$operator_addr" --delay 0 ); then
                echo "[seed] ✓ ${role} granted to operator on PN ${upper}"
            else
                echo "[seed] ! grant-business-role ${role} failed for operator on PN ${upper} — retry manually:"
                echo "    (cd ../rayls-sovereign-contracts && npx hardhat grant-business-role --pn ${upper} --role ${role} --account ${operator_addr} --delay 0)"
            fi
        done
        for role in $USER_ROLES; do
            echo "[seed] Granting ${role} to user on PN ${upper}…"
            if ( cd ../rayls-sovereign-contracts &&
                 npx hardhat grant-business-role \
                     --pn "$upper" --role "$role" --account "$user_addr" --delay 0 ); then
                echo "[seed] ✓ ${role} granted to user on PN ${upper}"
            else
                echo "[seed] ! grant-business-role ${role} failed for user on PN ${upper} — retry manually:"
                echo "    (cd ../rayls-sovereign-contracts && npx hardhat grant-business-role --pn ${upper} --role ${role} --account ${user_addr} --delay 0)"
            fi
        done

        # Inject this node's minted JWTs into the tests-automation .env as per-node
        # auth keys (mirrors the OPS_SERVICE_<NODE>_URL convention). The .env was
        # overwritten from the contracts .env earlier, so set_env_var updates/appends
        # on top of it.
        if [ -f "$automation_env" ]; then
            set_env_var "$automation_env" "OPS_SERVICE_${upper}_OPERATOR_AUTH_KEY" "$operator_key"
            set_env_var "$automation_env" "OPS_SERVICE_${upper}_USER_AUTH_KEY" "$user_key"
            echo "[seed] ✓ Wrote OPS_SERVICE_${upper}_{OPERATOR,USER}_AUTH_KEY into ${automation_env}"
        else
            echo "[seed] tests-automation .env not found at ${automation_env} — skipping token injection for ${upper}"
        fi

        index=$((index + 1))
    done
}

# Tear down EVERY Blockscout container the script may have created, regardless of
# BLOCKSCOUT_ENABLED — the per-instance projects AND the shared-external-db project
# live outside the main compose project, so the main `down --remove-orphans` never
# reaches them. Runs on --clean so a differently-flagged prior run leaves nothing
# behind. No-op when the blockscout repo isn't cloned.
clean_all_blockscout() {
    local BS_DIR="../rayls-privacy-blockscout"
    [ -d "$BS_DIR/docker-compose" ] || return 0
    echo "==> CLEAN: tearing down all Blockscout containers (shared + instances)..."

    # Per-instance projects
    if [ -d "$BS_DIR/docker-compose/instances" ]; then
        for d in "$BS_DIR"/docker-compose/instances/*/; do
            [ -d "$d" ] || continue
            ( cd "$d" && docker compose down --remove-orphans --volumes 2>/dev/null ) || true
        done
        # Wipe generated instance dirs so new-instance.sh regenerates them for the
        # current participants/ports on the next up.
        find "$BS_DIR/docker-compose/instances" -mindepth 1 -maxdepth 1 -type d -exec rm -rf {} +
    fi

    # Shared services project (redis/visualizer/sig-provider + shared-redis-data)
    ( cd "$BS_DIR/docker-compose" && docker compose -f shared-external-db.yml down --remove-orphans --volumes 2>/dev/null ) || true
}

clean_base_services() {
    # Remove marker — will be recreated after successful deployment
    rm -f "$DEPLOY_MARKER"

    # Sync gnark verifier contracts to rayls-sovereign-contracts before building.
    # This ensures the on-chain verifiers match the gnark proving keys.
    GNARK_BUILD="../rayls-sovereign-gnark-api/last_build"
    CONTRACTS_PAYMENTS="../rayls-sovereign-contracts/src/rayls-protocol/Enygma/Enygma-Payments"
    CONTRACTS_DVP="../rayls-sovereign-contracts/src/rayls-protocol/Enygma/Enygma-DVP"
    if [ -d "$GNARK_BUILD" ] && [ -d "$CONTRACTS_PAYMENTS" ]; then
        echo "Syncing gnark verifier contracts to rayls-sovereign-contracts..."
        DVP_FILES="EnygmaJoinSplitVerifier Erc721OwnershipVerifier Erc1155JoinSplitVerifier"
        synced=0
        for sol_file in "$GNARK_BUILD"/*.sol; do
            bn=$(basename "$sol_file")
            [[ "$bn" == *"_raw"* ]] && continue
            dst="$CONTRACTS_PAYMENTS/$bn"
            for dvp in $DVP_FILES; do
                [ "${bn%.sol}" = "$dvp" ] && dst="$CONTRACTS_DVP/$bn" && break
            done
            if ! diff -q "$sol_file" "$dst" > /dev/null 2>&1; then
                cp "$sol_file" "$dst"
                synced=$((synced + 1))
            fi
        done
        echo "  Synced $synced verifier contract(s)."
    else
        echo "Warning: gnark-api last_build or contracts dir not found, skipping verifier sync."
    fi

    # shut down all services. Both modes share docker-compose.dev-local.yml; --remove-orphans
    # also clears any containers from the previous mode.
    docker compose -f docker-compose.dev-local.yml down --remove-orphans --volumes
    build_images "$BASE_SERVICES"
    compose_env up -d --force-recreate --remove-orphans --wait $BASE_SERVICES

    # Mark successful deployment
    echo "deployed=$(date -u +%Y-%m-%dT%H:%M:%SZ) mode=$DEV_MODE participants=$PARTICIPANT_LIST" > "$DEPLOY_MARKER"
}

# Bring up Blockscout per participant. Uses the standalone compose project in
# ../rayls-privacy-blockscout: one set of shared services (postgres, redis,
# visualizer, sig-provider) on the `blockscout-shared` network plus one instance
# (backend + frontend + proxy + stats + stats-db) per privacy node.
start_blockscout() {
    [ "$BLOCKSCOUT_ENABLED" = "true" ] || return 0

    local BS_DIR="../rayls-privacy-blockscout"
    echo "==> BLOCKSCOUT: setting up explorer instances..."

    # On --clean the teardown + generated-dir wipe is handled unconditionally by
    # clean_all_blockscout (called from the clean branch), so nothing to do here.

    # Shared services: redis, visualizer, sig-provider. The Postgres for both
    # blockscout_pn-* and stats_pn-* lives in the relayer's `postgres` service
    # (see DB_NAMES above), so we use shared-external-db.yml instead of the
    # standalone shared.yml. Idempotent.
    ( cd "$BS_DIR/docker-compose" && docker compose -f shared-external-db.yml up -d )

    if [ "$CLEAN_MODE" = "clean" ]; then
        # Privacy node chains were just wiped; drop the matching blockscout
        # databases so the backend re-indexes from genesis on a clean slate.
        # clean_base_services already did `down --volumes` on postgres, so on
        # a true clean these DROPs are no-ops; they cover re-runs that don't
        # nuke the volume.
        for participant in "${PARTICIPANTS[@]}"; do
            local lc
            lc=$(echo "$participant" | tr '[:upper:]' '[:lower:]')
            compose_env exec -T postgres psql -U admin -d postgres \
                -c "DROP DATABASE IF EXISTS \"blockscout_pn-$lc\" WITH (FORCE);" >/dev/null 2>&1 || true
            compose_env exec -T postgres psql -U admin -d postgres \
                -c "DROP DATABASE IF EXISTS \"stats_pn-$lc\" WITH (FORCE);" >/dev/null 2>&1 || true
        done
    fi

    local idx=0
    for participant in "${PARTICIPANTS[@]}"; do
        local lc
        lc=$(echo "$participant" | tr '[:upper:]' '[:lower:]')
        local name="pn-$lc"
        local http_port=$((8480 + idx * 100))
        local env_file=".${participant}.env"

        local chain_id
        chain_id=$(grep -E '^PRIVACY_NODE_CHAIN_ID=' "$env_file" | cut -d= -f2)
        local rpc_url
        rpc_url=$(grep -E '^PRIVACY_NODE_RPC_URL=' "$env_file" | cut -d= -f2)

        # Use the direct compose hostname (e.g. http://pn-a:8545) — the
        # external-db override attaches Blockscout's backend to the relayer's
        # compose network, so it can resolve the chain by its service name.
        # The previous host.docker.internal rewrite is no longer needed and
        # would fail anyway, since the chain binds 127.0.0.1 on the host while
        # host.docker.internal points at the bridge gateway.

        # FIRST_BLOCK: where Blockscout's catchup indexer starts. The local chains are
        # born at genesis with a tiny history, so 0 (full index) is cheap.
        local first_block=0

        local instance_dir="$BS_DIR/docker-compose/instances/$name"
        if [ ! -d "$instance_dir" ]; then
            # Pass redis-db-index explicitly (idx+1, range 1..6) — new-instance.sh's
            # auto-assignment trips on `set -o pipefail` when the instances dir is empty.
            ( cd "$BS_DIR" && ./scripts/new-instance.sh "$name" "$http_port" "$rpc_url" "$chain_id" "$first_block" $((idx + 1)) )
        fi

        # Generate the external-db override for this instance. It re-points
        # backend's DATABASE_URL and stats' DB URLs to the relayer's `postgres`
        # (admin/admin) and attaches both services to the relayer's compose
        # network so they can resolve `postgres` by hostname. Written inline
        # (rather than sed'ing a checked-in template) so the blockscout repo
        # needs no extra hand-placed file — everything under instances/ is
        # generated by this script.
        cat > "$instance_dir/docker-compose-external-db.override.yml" <<EOF
# Generated by rayls-sovereign-relayer/start_dev.sh — do not edit by hand.
# Re-points Blockscout's backend + stats databases at the relayer stack's
# \`postgres\` service (admin/admin) and attaches them to its compose network.
networks:
  rayls-sovereign-relayer_default:
    external: true
    name: rayls-sovereign-relayer_default

services:
  backend:
    environment:
      DATABASE_URL: postgresql://admin:admin@postgres:5432/blockscout_${name}
    networks:
      - blockscout-shared
      - rayls-sovereign-relayer_default

  # chain-guard runs BEFORE backend (backend depends_on it completing) and must
  # reach the privacy node's RPC by service name. The chain resolves only on the
  # relayer's compose network, so attach the guard there too — otherwise it can't
  # resolve the RPC host, curl returns empty, and the guard fails the whole stack.
  chain-guard:
    networks:
      - blockscout-shared
      - rayls-sovereign-relayer_default

  stats:
    environment:
      STATS__DB_URL: postgresql://admin:admin@postgres:5432/stats_${name}
      STATS__BLOCKSCOUT_DB_URL: postgresql://admin:admin@postgres:5432/blockscout_${name}
      STATS__CREATE_DATABASE: "true"
      STATS__RUN_MIGRATIONS: "true"
    networks:
      - blockscout-shared
      - rayls-sovereign-relayer_default
EOF

        ( cd "$instance_dir" && docker compose -f docker-compose.yml -f docker-compose-external-db.override.yml up -d )

        echo "    Blockscout for $name: http://localhost:$http_port"
        idx=$((idx + 1))
    done
}

if [ "$CLEAN_MODE" = "clean" ]; then
    # CLEAN: nuke everything — volumes, contracts, nodes — and redeploy from scratch
    echo "==> CLEAN: tearing down everything and redeploying from scratch..."
    clean_all_blockscout
    clean_base_services

elif [ "$CLEAN_MODE" = "soft" ]; then
    # SOFT CLEAN: wipe relayer + pubrelayer databases. Keep CTS, governance, contracts, and chain nodes.
    echo "==> SOFT CLEAN: wiping relayer/pubrelayer databases (keeping CTS, governance, contracts, chain nodes)..."

    # Postgres must be up to issue DROP/CREATE
    compose_env up -d --wait postgres

    # Stop only the services that hold connections to the relayer DBs.
    # CTS, governance, and chain nodes stay running.
    RELAYER_SERVICES=""
    for participant in "${PARTICIPANTS[@]}"; do
        LOWERCASE_PARTICIPANT=$(echo "$participant" | tr '[:upper:]' '[:lower:]')
        RELAYER_SERVICES+="relayer-$LOWERCASE_PARTICIPANT pubrelayer-$LOWERCASE_PARTICIPANT "
    done
    compose_env stop $RELAYER_SERVICES 2>/dev/null || true

    # Drop and recreate relayer + pubrelayer databases. CTS DBs are preserved.
    for participant in "${PARTICIPANTS[@]}"; do
        PARTICIPANT_ENV_FILE=".${participant}.env"
        PRIVATE_DB=$(grep -E '^PRIVATE_RELAYER_DATABASE_CONNECTIONSTRING=' "$PARTICIPANT_ENV_FILE" | sed -n 's|.*://[^/]*/\([^?]*\).*|\1|p')
        PUBLIC_DB=$(grep -E '^RAYLS_NODE_DATABASE_CONNECTIONSTRING=' "$PARTICIPANT_ENV_FILE" | sed -n 's|.*://[^/]*/\([^?]*\).*|\1|p')
        for db in "$PRIVATE_DB" "$PUBLIC_DB"; do
            [ -z "$db" ] && continue
            echo "  Recreating database: $db"
            compose_env exec -T postgres psql -U admin -d postgres -c "DROP DATABASE IF EXISTS \"$db\" WITH (FORCE);" >/dev/null
            compose_env exec -T postgres psql -U admin -d postgres -c "CREATE DATABASE \"$db\";" >/dev/null
        done
    done

    # Ensure remaining base services are up
    compose_env up -d --wait $BASE_SERVICES

else
    # RESUME: just ensure base services are running, restart dev services with existing state
    echo "==> RESUME: restarting dev services with existing state..."

    # Ensure base services are running (they may have been stopped)
    compose_env up -d --wait $BASE_SERVICES
fi

# Shutdown Otel if disabled
if [ "$OTEL_SDK_DISABLED" = "true" ]; then
    docker compose -f $COMPOSE_FILE down otel
else
    # Launch otel services if they're not up yet.
    OTEL_NEEDS_LAUNCHING=true

    # Check if service otel-collector is "Up x time"
    set +e
    OTEL_COLLECTOR_STATUS=$(./compose.sh ps -a otel | awk 'NR>1 {print $9}')
    OTEL_CHECK_RC=$?
    set -e
    if [ $OTEL_CHECK_RC -ne 0 ]; then
        OTEL_NEEDS_LAUNCHING=true
    elif [ "$OTEL_COLLECTOR_STATUS" = "Up" ]; then
        OTEL_NEEDS_LAUNCHING=false
    fi

    if [ "$OTEL_NEEDS_LAUNCHING" = "true" ]; then
        compose_env up -d --remove-orphans otel
    fi
fi

# Handle Governance services

ALL_GOVERNANCE_SERVICES+="governance-listener governance-flagger governance-api auditor-explorer"
NON_WATCHABLE_GOVERNANCE_SERVICES=""
WATCHABLE_GOVERNANCE_SERVICES+="governance-listener governance-flagger "

#ENABLED_GOVERNANCE_SERVICES=$NON_WATCHABLE_GOVERNANCE_SERVICES

if [ "$GOVERNANCE_API_ENABLED" = "true" ]; then
    WATCHABLE_GOVERNANCE_SERVICES+="governance-api "
    if [ "$EXPLORERS_ENABLED" = "true" ]; then
        WATCHABLE_GOVERNANCE_SERVICES+="auditor-explorer "
    else
        docker compose -f $COMPOSE_FILE down auditor-explorer
    fi
else
    docker compose -f $COMPOSE_FILE down governance-api auditor-explorer
fi

if [ "$GOVERNANCE_ENABLED" = "true" ]; then
    DEV_SERVICES+=$WATCHABLE_GOVERNANCE_SERVICES
else
    docker compose -f $COMPOSE_FILE down $ALL_GOVERNANCE_SERVICES
fi

# Handle Ops API services (custody + api + worker, one trio per participant).
# Built per-participant like the playground loop below — the compose file defines
# custody/ops-api/ops-worker blocks for every letter a–f.
ALL_OPS_API_SERVICES=""
for participant in "${PARTICIPANTS[@]}"; do
    olc=$(echo "$participant" | tr '[:upper:]' '[:lower:]')
    ALL_OPS_API_SERVICES+=" custody-$olc ops-api-$olc ops-worker-$olc"
done
if [ "$OPS_API_ENABLED" = "true" ]; then
    DEV_SERVICES+=" $ALL_OPS_API_SERVICES "
else
    docker compose -f $COMPOSE_FILE down $ALL_OPS_API_SERVICES 2>/dev/null || true
fi

# Handle Playground (Next.js dApp). One container per chain: Hub + 1 per
# Privacy Node. Each binds to a single chain via NEXT_PUBLIC_CHAIN_* envs.
ALL_PLAYGROUND_SERVICES="playground-hub"
for participant in "${PARTICIPANTS[@]}"; do
    plc=$(echo "$participant" | tr '[:upper:]' '[:lower:]')
    ALL_PLAYGROUND_SERVICES+=" playground-$plc"
done
if [ "$PLAYGROUND_ENABLED" = "true" ]; then
    DEV_SERVICES+=" $ALL_PLAYGROUND_SERVICES "
else
    docker compose -f $COMPOSE_FILE down $ALL_PLAYGROUND_SERVICES 2>/dev/null || true
fi

# Sync contracts and .env to the tests-automation repo (if present)
AUTOMATION_PATH="../rayls-sovereign-tests-automation"
CONTRACTS_PATH="../rayls-sovereign-contracts"

if [ -d "$AUTOMATION_PATH" ]; then
    echo "📋 Syncing tests-automation repo with deployed contracts..."

    # Check if the automation's copy of contracts diverges from the contracts repo
    CONTRACTS_SRC_HASH=$(git -C "$CONTRACTS_PATH" log -1 --format=%H -- src/)
    mkdir -p "$AUTOMATION_PATH/cache"
    AUTOMATION_SYNC_MARKER="$AUTOMATION_PATH/cache/.contracts-sync-hash"

    NEEDS_SYNC=false
    if [ ! -f "$AUTOMATION_SYNC_MARKER" ]; then
        NEEDS_SYNC=true
    elif [ "$(cat "$AUTOMATION_SYNC_MARKER")" != "$CONTRACTS_SRC_HASH" ]; then
        NEEDS_SYNC=true
    fi

    if [ "$NEEDS_SYNC" = "true" ]; then
        echo "🔄 Contracts source diverged — running sync-contracts-local in automation repo..."
        (cd "$AUTOMATION_PATH" && npm run sync-contracts-local)
        echo "$CONTRACTS_SRC_HASH" > "$AUTOMATION_SYNC_MARKER"
        echo "✅ Contracts synced and compiled in tests-automation"
    else
        echo "✅ Automation contracts already in sync (hash: ${CONTRACTS_SRC_HASH:0:7})"
    fi

    # Copy .env from contracts to automation (contains all deployed addresses)
    if [ -f "$CONTRACTS_PATH/.env" ]; then
        cp "$CONTRACTS_PATH/.env" "$AUTOMATION_PATH/.env"
        echo "✅ Copied contracts .env → tests-automation/.env"
    else
        echo "⚠️  No .env found in contracts repo — automation tests may fail"
    fi
else
    echo "ℹ️  tests-automation repo not found at $AUTOMATION_PATH — skipping sync"
fi

# Start CTS and Relayer
#
# BUILD_TAGS=faultinjection is passed to the dev-image builds so the
# relayer/pubrelayer/cts binaries include the fault-injection HTTP server.
# The Dockerfiles (private-relayer/Dockerfile.dev, public-relayer/Dockerfile.dev,
# cts/Dockerfile.dev) default the ARG to empty — meaning a manual
# `docker compose build` keeps FI compiled out. start_dev.sh is the canonical
# dev entrypoint and the e2e resilience suite expects FI to be reachable on the
# FAULT_INJECTION_PORT, so we opt in here. Override with BUILD_TAGS="" if you
# want a pure dev image without chaos instrumentation (e.g. for profiling).
BUILD_TAGS=${BUILD_TAGS-faultinjection}

build_images "$DEV_SERVICES"

# Clean up stale Docker resources to prevent unbounded storage growth.
# - dangling images: old untagged layers left behind after rebuilds
# - build cache older than 24h: intermediate layers from multi-stage builds
docker image prune -f 2>/dev/null || true
docker builder prune -f -a --filter "until=24h" 2>/dev/null || true

# Start Blockscout BEFORE the dev services come up: it creates and migrates the
# blockscout_pn-* databases, which ops-api/ops-worker depend on (they install
# LISTEN/NOTIFY triggers and backfill tokens against those tables).
start_blockscout

# Extract the per-participant ops-api values that the contracts deploy wrote into
# ../rayls-sovereign-contracts/.env:
#   - DEPLOYMENT_PROXY_REGISTRY: lets the ops-api resolve its RaylsAccessManager.
#   - ACCESS_MANAGER_STARTING_BLOCK: the chain tip captured just before the AM deploy
#     (<= the AM deploy block), so the ops-worker's AM role-event indexer backfills from
#     deploy-time instead of genesis (critical on long-lived remote PNs). The DB cursor
#     overrides it once the first backfill completes, so it only matters on a fresh DB.
# Empty fallbacks keep the script working when ops-api is disabled or the values haven't
# been deployed yet (compose defaults STARTING_BLOCK to 0 -> backfill from genesis).
if [ "$OPS_API_ENABLED" = "true" ] && [ -f "../rayls-sovereign-contracts/.env" ]; then
    for participant in "${PARTICIPANTS[@]}"; do
        lc=$(echo "$participant" | tr '[:upper:]' '[:lower:]')
        upper=$(echo "$participant" | tr '[:lower:]' '[:upper:]')
        addr=$(grep -E "^PRIVACY_NODE_${upper}_DEPLOYMENT_PROXY_REGISTRY=" \
            "../rayls-sovereign-contracts/.env" | cut -d= -f2)
        export "OPS_API_${upper}_DEPLOYMENT_PROXY_REGISTRY=$addr"
        am_block=$(grep -E "^PRIVACY_NODE_${upper}_ACCESS_MANAGER_STARTING_BLOCK=" \
            "../rayls-sovereign-contracts/.env" | cut -d= -f2)
        export "OPS_API_${upper}_STARTING_BLOCK=$am_block"
    done
fi

# Hub DeploymentProxyRegistry — used by playground-hub for the Private Hub
# chain's contract resolution. Same source file as the per-PN registries above.
if [ -f "../rayls-sovereign-contracts/.env" ]; then
    pnh_addr=$(grep -E "^PNH_DEPLOYMENT_PROXY_REGISTRY=" \
        "../rayls-sovereign-contracts/.env" | cut -d= -f2)
    if [ -n "$pnh_addr" ]; then
        export PNH_DEPLOYMENT_PROXY_REGISTRY="$pnh_addr"
    fi
fi

# Bring the dev services up DETACHED first, so the ops-api admin bootstrap + e2e
# seeding can reach them over HTTP/RPC without their logs streaming yet.
compose_env up -d --force-recreate --remove-orphans $DEV_SERVICES

# Run the ops-api admin bootstrap + e2e user seeding as a single sequential background
# job, then WAIT for it to finish so its output is shown cleanly before the container
# log stream starts. bootstrap_ops_api polls each ops-api's /health, calls
# /admin/bootstrap and grants the admin's on-chain roles; seed_e2e_users then seeds the
# e2e operator/user, grants their roles, and writes the JWTs into the tests-automation
# .env. They run sequentially (not in parallel) because both sign grant-business-role
# txs with the same admin key — concurrent grants would race on its nonce. No-ops when
# OPS_API_ENABLED is not "true". `|| true` keeps a non-zero job status from tripping the
# ERR trap (the functions log their own failures non-fatally).
{ bootstrap_ops_api; seed_e2e_users; } &
wait $! || true

# Now attach: stream the running container logs and enable source file-watch for dev.
# Services are already running from the detached `up` above, so this is a no-op
# recreate-wise (same compose config) and just attaches + watches. `--progress quiet`
# mutes Compose's dependency-reconciliation reporter (the `Container ... Healthy/
# Waiting/Exited` lines for the already-up deps and the finished one-shot init
# containers), which would otherwise interleave with the attached log stream. The log
# stream, `--watch` rebuild output, and Go compile/panic output are unaffected.
compose_env --progress quiet up --watch $DEV_SERVICES
