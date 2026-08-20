#!/usr/bin/env bash
#
# Generates a self-signed CA + leaf certificates for local mTLS development
# between the relayers, CTS, and NATS.
#
# Layout:
#   cts/certs/                       ca.crt, server.crt, server.key (gRPC server)
#                                    cts.crt,    cts.key            (NATS client)
#   private-relayer/certs/           ca.crt, private-relayer.crt,
#                                    private-relayer.key            (gRPC + NATS client)
#   public-relayer/certs/            ca.crt, public-relayer.crt,
#                                    public-relayer.key             (gRPC + NATS client)
#   docker/development/nats-certs/   ca.crt, server.crt, server.key (NATS server)
#   ../rayls-privacy-pnh-governance-api/certs/
#                                    ca.crt, governance.crt,
#                                    governance.key                 (NATS client; copied
#                                                                    only if the dir exists)
#
# All leaves are signed by the same CA. The CA's public half (ca.crt) is
# distributed everywhere so each peer can verify the others.
#
# Idempotent: skips if docker/development/nats-certs/server.crt already
# exists. Pass --force to regenerate.

set -euo pipefail

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd -P)"
REPO_ROOT="$(cd -- "$SCRIPT_DIR/.." && pwd -P)"

CTS_DIR="$REPO_ROOT/cts/certs"
PRIV_DIR="$REPO_ROOT/private-relayer/certs"
PUB_DIR="$REPO_ROOT/public-relayer/certs"
NATS_DIR="$REPO_ROOT/docker/development/nats-certs"
GOV_REPO="$REPO_ROOT/../rayls-privacy-pnh-governance-api"
GOV_DIR="$GOV_REPO/certs"

FORCE=false
if [[ "${1:-}" == "--force" ]]; then
    FORCE=true
fi

if [[ "$FORCE" != true && -f "$NATS_DIR/server.crt" ]]; then
    echo "Dev certs already exist at $NATS_DIR/server.crt; skipping (use --force to regenerate)."
    exit 0
fi

if ! command -v openssl >/dev/null 2>&1; then
    echo "Error: openssl is required but not installed." >&2
    exit 1
fi

mkdir -p "$CTS_DIR" "$PRIV_DIR" "$PUB_DIR" "$NATS_DIR"

WORK="$(mktemp -d)"
trap 'rm -rf "$WORK"' EXIT

DAYS=825

echo "==> Generating Rayls Dev CA..."
openssl genrsa -out "$WORK/ca.key" 4096 >/dev/null 2>&1
openssl req -x509 -new -nodes -key "$WORK/ca.key" -sha256 -days "$DAYS" \
    -subj "/CN=Rayls Dev CA/O=Rayls/OU=Local Development" \
    -out "$WORK/ca.crt" >/dev/null 2>&1

gen_leaf() {
    local name="$1"
    local cn="$2"
    local ext_section="$3"

    openssl genrsa -out "$WORK/$name.key" 4096 >/dev/null 2>&1
    openssl req -new -key "$WORK/$name.key" \
        -subj "/CN=$cn/O=Rayls/OU=Local Development" \
        -out "$WORK/$name.csr" >/dev/null 2>&1
    openssl x509 -req -in "$WORK/$name.csr" \
        -CA "$WORK/ca.crt" -CAkey "$WORK/ca.key" -CAcreateserial \
        -days "$DAYS" -sha256 \
        -extfile "$WORK/$name.ext" -extensions "$ext_section" \
        -out "$WORK/$name.crt" >/dev/null 2>&1
}

write_client_ext() {
    cat > "$WORK/$1.ext" <<EOF
[v3_client]
basicConstraints     = CA:FALSE
keyUsage             = digitalSignature, keyEncipherment
extendedKeyUsage     = clientAuth
EOF
}

# --- CTS gRPC server cert ---
echo "==> Generating CTS gRPC server cert (SANs: localhost, cts, cts-a..cts-f, 127.0.0.1)..."
cat > "$WORK/cts-server.ext" <<'EOF'
[v3_server]
basicConstraints     = CA:FALSE
keyUsage             = digitalSignature, keyEncipherment
extendedKeyUsage     = serverAuth
subjectAltName       = @alt_names

[alt_names]
DNS.1 = localhost
DNS.2 = cts
DNS.3 = cts-a
DNS.4 = cts-b
DNS.5 = cts-c
DNS.6 = cts-d
DNS.7 = cts-e
DNS.8 = cts-f
IP.1  = 127.0.0.1
EOF
gen_leaf "cts-server" "cts" "v3_server"

# --- NATS server cert ---
echo "==> Generating NATS server cert (SANs: localhost, nats, 127.0.0.1)..."
cat > "$WORK/nats-server.ext" <<'EOF'
[v3_server]
basicConstraints     = CA:FALSE
keyUsage             = digitalSignature, keyEncipherment
extendedKeyUsage     = serverAuth
subjectAltName       = @alt_names

[alt_names]
DNS.1 = localhost
DNS.2 = nats
IP.1  = 127.0.0.1
EOF
gen_leaf "nats-server" "nats" "v3_server"

# --- client certs ---
echo "==> Generating private-relayer client cert..."
write_client_ext "private-relayer"
gen_leaf "private-relayer" "private-relayer" "v3_client"

echo "==> Generating public-relayer client cert..."
write_client_ext "public-relayer"
gen_leaf "public-relayer" "public-relayer" "v3_client"

echo "==> Generating cts client cert (NATS)..."
write_client_ext "cts-client"
gen_leaf "cts-client" "cts" "v3_client"

echo "==> Generating governance client cert (NATS)..."
write_client_ext "governance"
gen_leaf "governance" "governance" "v3_client"

install_pub() { install -m 0644 "$1" "$2"; }
install_key() { install -m 0600 "$1" "$2"; }

echo "==> Distributing certs..."

# CTS: gRPC server cert + NATS client cert
install_pub "$WORK/ca.crt"               "$CTS_DIR/ca.crt"
install_pub "$WORK/cts-server.crt"       "$CTS_DIR/server.crt"
install_key "$WORK/cts-server.key"       "$CTS_DIR/server.key"
install_pub "$WORK/cts-client.crt"       "$CTS_DIR/cts.crt"
install_key "$WORK/cts-client.key"       "$CTS_DIR/cts.key"

# Private relayer: gRPC + NATS client (same cert)
install_pub "$WORK/ca.crt"               "$PRIV_DIR/ca.crt"
install_pub "$WORK/private-relayer.crt"  "$PRIV_DIR/private-relayer.crt"
install_key "$WORK/private-relayer.key"  "$PRIV_DIR/private-relayer.key"

# Public relayer: gRPC + NATS client (same cert)
install_pub "$WORK/ca.crt"               "$PUB_DIR/ca.crt"
install_pub "$WORK/public-relayer.crt"   "$PUB_DIR/public-relayer.crt"
install_key "$WORK/public-relayer.key"   "$PUB_DIR/public-relayer.key"

# NATS server
install_pub "$WORK/ca.crt"               "$NATS_DIR/ca.crt"
install_pub "$WORK/nats-server.crt"      "$NATS_DIR/server.crt"
install_key "$WORK/nats-server.key"      "$NATS_DIR/server.key"

# Governance (best-effort — only if the sibling repo exists)
if [[ -d "$GOV_REPO" ]]; then
    mkdir -p "$GOV_DIR"
    install_pub "$WORK/ca.crt"           "$GOV_DIR/ca.crt"
    install_pub "$WORK/governance.crt"   "$GOV_DIR/governance.crt"
    install_key "$WORK/governance.key"   "$GOV_DIR/governance.key"
    GOV_MSG="  $GOV_DIR (governance repo)"
else
    GOV_MSG="  governance repo not found at $GOV_REPO — governance.crt not distributed"
fi

echo
echo "Done. Certs distributed to:"
echo "  $CTS_DIR"
echo "  $PRIV_DIR"
echo "  $PUB_DIR"
echo "  $NATS_DIR"
echo "$GOV_MSG"
