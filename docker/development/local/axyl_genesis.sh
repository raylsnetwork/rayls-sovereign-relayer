#!/bin/bash
CHAIN_ID=${CHAIN_ID:-0x3039}
DEV_FUNDED_ACCOUNT=${DEV_FUNDED_ACCOUNT:-0xE4F2eB9B68cd50c604B5C16258167F77B6b23F31}
GAS_LIMIT=${GAS_LIMIT:-0xffffffffffff}

# Run-once guard, mirroring axyl's canonical local-testnet.sh (`if [ -d ROOTDIR ]`)
# and the relayer's own setup_validator.sh (`if [ ! -d node-keys ]`). The genesis
# ceremony must run exactly once per chain lifetime: `rayls genesis` stamps `now()` into
# the genesis block, so re-running it produces a different genesis hash and breaks
# validator restarts ("genesis hash in the storage does not match the specified
# chainspec", reth init.rs:238). Compose re-runs this one-shot when re-resolving the
# `contracts -> pn-a/pn-b -> pn-*-genesis` dependency graph during `up --force-recreate`;
# this guard makes that re-run a no-op. A clean deploy (`start_dev.sh -c`) runs
# `down --volumes`, wiping genesis-data, so a fresh genesis is generated as intended.
if [ -f /home/nonroot/data/genesis/genesis.yaml ]; then
    echo "Genesis already generated at /home/nonroot/data/genesis/genesis.yaml -- skipping (run-once)."
    exit 0
fi

GASLESS_FLAGS=""
if [ "$GASLESS" = "true" ]; then
    GASLESS_FLAGS="--base-fee 0 --min-base-fee 0"
fi

ACCOUNTS_FLAG=""
if [ -f /accounts.yaml ]; then
    ACCOUNTS_FLAG="--accounts /accounts.yaml"
fi

mkdir -p /home/nonroot/data/genesis/validators
for i in 1; do
  cp /home/nonroot/data/validator-$i/node-info.yaml \
     /home/nonroot/data/genesis/validators/validator-$i.yaml
done

# RAYLS_AXYL_EPOCH_DURATION_SECS overrides the chain's genesis epoch duration.
# Default 86400s (24h). Lowering this triggers axyl epoch transitions more
# frequently, which respawns the proposer task. Persisted into the chain's
# genesis config, so it's baked in for the lifetime of the dev env.
EPOCH_DURATION_SECS=${RAYLS_AXYL_EPOCH_DURATION_SECS:-86400}

/usr/local/bin/rayls genesis \
    --datadir /home/nonroot/data/ \
    --chain-id $CHAIN_ID \
    --epoch-duration-in-secs $EPOCH_DURATION_SECS \
    --dev-funded-account $DEV_FUNDED_ACCOUNT \
    --max-header-delay-ms 1000 \
    --min-header-delay-ms 1000 \
    --consensus-registry-owner $DEV_FUNDED_ACCOUNT \
    --gas-limit $GAS_LIMIT \
    $GASLESS_FLAGS \
    $ACCOUNTS_FLAG

for i in 1; do
    mkdir -p /home/nonroot/data/validator-$i/genesis/
    cp /home/nonroot/data/genesis/genesis.yaml \
       /home/nonroot/data/genesis/committee.yaml \
       /home/nonroot/data/validator-$i/genesis/
    cp /home/nonroot/data/parameters.yaml /home/nonroot/data/validator-$i/
    # `node --dev` auto-bootstrap treats a datadir without this sentinel as empty and
    # re-runs its own dev ceremony — regenerating the validator key and overwriting
    # THIS genesis with the built-in dev one (chain-id 2017, gasless). The datadir is
    # ceremony-initialized here, so write the sentinel to make the bootstrap a no-op.
    # Sentinel name: BOOTSTRAP_SENTINEL in axyl crates/infrastructure/network-cli/src/dev/mod.rs.
    touch /home/nonroot/data/validator-$i/.dev-bootstrap-complete
done
chown -R 1101:1101 /home/nonroot/data

echo "Genesis created: chain-id=$CHAIN_ID, gasless=$GASLESS, gas-limit=$GAS_LIMIT, funded=$DEV_FUNDED_ACCOUNT"
