<div align="center">

# Rayls Relayer

**Cross-chain messaging with full end-to-end encryption, routed through the Private Network Hub.**

[![License: Apache 2.0][license-badge]][license-url]
[![Go][go-badge]][go-url]

[![Discord][discord-badge]][discord-url]
[![X][x-badge]][x-url]
[![LinkedIn][linkedin-badge]][linkedin-url]
[![YouTube][youtube-badge]][youtube-url]

[Getting started](#getting-started) | [Configuration](#configuration) | [Commands](#commands) | [Developing locally](./docs/dev/README.md)

</div>

## What is the Rayls Relayer?

The Rayls relayer moves messages and assets between Privacy Nodes without exposing their contents. It routes traffic through an intermediary chain — the Private Network Hub (PNH) — and keeps the payloads encrypted end to end, so the hub relays messages it cannot read. Message handling is generic and follows [EIP-5164](https://eips.ethereum.org/EIPS/eip-5164).

It runs as three Go binaries:

- **`relayer`** (private relayer) — listens to a Privacy Node and executes messages on it.
- **`pubrelayer`** (public relayer) — bridges to the public chain.
- **`cts`** (Cryptography Trust Suite) — key management and signing.

## Features

- Cross-chain transfers of assets and messages.
- Privacy-preserving cross-chain communication — payloads stay encrypted through the hub.
- Generic message handling (EIP-5164).

## Getting started

### Prerequisites

- Go 1.26.4+
- Two running EVM-based Privacy Nodes
- An EVM-based Private Network Hub

### Build

```sh
git clone https://github.com/raylsnetwork/rayls-sovereign-relayer.git
cd rayls-sovereign-relayer
make
```

`make` builds all three binaries into `build/`.

### Configure

Set up your env file. During development the three services can share a single env file. Copy [`.env.example`](./.env.example), fill in your own endpoints and credentials, and pass it with `--env`. See [Configuration](#configuration) for what the variables mean.

### Run

```sh
./build/relayer run --env path/to/env
```

See [Commands](#commands) for the other two binaries.

## Configuration

Configuration is via environment variables (loaded with [`spf13/viper`](https://github.com/spf13/viper)). The authoritative, up-to-date list of every variable — with example values — is in [`.env.example`](./.env.example).

The variables are grouped by component:

| Prefix | Component |
| --- | --- |
| `PRIVACY_NODE_*` | The Privacy Node the private relayer listens to and executes on |
| `PNH_*` | The Private Network Hub (the intermediary chain) |
| `PUBLIC_CHAIN_*` | The public chain (used by the public relayer) |
| `PRIVATE_RELAYER_*` / `RAYLS_NODE_*` | Relayer database and node settings |
| `CTS_*` | Cryptography Trust Suite (key management / KMS) |
| `OTEL_*` | OpenTelemetry export |

## Commands

Run each binary with the `run` subcommand and an env file:

```sh
./build/relayer run --env path/to/env      # private relayer
./build/pubrelayer run --env path/to/env   # public relayer
./build/cts run --env path/to/env          # Cryptography Trust Suite
```

## Developing

To run the Cryptography Trust Suite, private relayer, and public relayer locally, see the [Developing locally with Docker](./docs/dev/README.md) guide.

## Contributing

We are not accepting external contributions at this time — see [CONTRIBUTING.md](./CONTRIBUTING.md). Please also read our [Code of Conduct](./CODE_OF_CONDUCT.md).

## Security

To report a security vulnerability, see [SECURITY.md](./SECURITY.md) — please do not open a public issue.

## License

Licensed under the Apache License, Version 2.0 — see [LICENSE](./LICENSE).

This project links third-party libraries that remain under their own licenses; the complete inventory is in [THIRD_PARTY_LICENSES.md](./THIRD_PARTY_LICENSES.md). Notably, it uses [go-ethereum](https://github.com/ethereum/go-ethereum) under the LGPL-3.0 (library packages only) and incorporates a Merkle-Patricia-Trie implementation derived from [zhangchiqing/merkle-patricia-trie](https://github.com/zhangchiqing/merkle-patricia-trie) (MIT). See [NOTICE](./NOTICE).

Copyright 2026 Rayls Core Ltd.

[license-badge]: https://img.shields.io/badge/License-Apache_2.0-blue.svg
[license-url]: ./LICENSE
[go-badge]: https://img.shields.io/badge/Go-1.26.4-00ADD8?logo=go&logoColor=white
[go-url]: ./go.mod
[discord-badge]: https://img.shields.io/badge/Discord-join%20chat-5865F2?logo=discord&logoColor=white
[discord-url]: https://discord.gg/6THZ96357r
[x-badge]: https://img.shields.io/badge/X-%40RaylsLabs-000000?logo=x&logoColor=white
[x-url]: https://x.com/RaylsLabs
[linkedin-badge]: https://img.shields.io/badge/LinkedIn-Rayls-0A66C2?logo=linkedin&logoColor=white
[linkedin-url]: https://www.linkedin.com/company/rayls/
[youtube-badge]: https://img.shields.io/badge/YouTube-Rayls-FF0000?logo=youtube&logoColor=white
[youtube-url]: https://www.youtube.com/@Rayls_blockchain
