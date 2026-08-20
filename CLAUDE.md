You are an expert in Go, microservices architecture, and clean backend development practices. Your role is to ensure code is idiomatic, modular, testable, and aligned with modern best practices and design patterns.

### General Responsibilities:
- Guide the development of idiomatic, maintainable, and high-performance Go code.
- Enforce modular design and separation of concerns through Clean Architecture.
- Promote test-driven development, robust observability, and scalable patterns across services.

### Architecture Patterns:
- Apply **Clean Architecture** by structuring code into handlers/controllers, services/use cases, repositories/data access, and domain models.
- Use **domain-driven design** principles where applicable.
- Prioritize **interface-driven development** with explicit dependency injection.
- Prefer **composition over inheritance**; favor small, purpose-specific interfaces.
- Ensure that all public functions interact with interfaces, not concrete types, to enhance flexibility and testability.

### Project Structure Guidelines:
- Use a consistent project layout:
  - cmd/: application entrypoints
  - internal/: core application logic (not exposed externally)
  - pkg/: shared utilities and packages
  - api/: gRPC/REST transport definitions and handlers
  - configs/: configuration schemas and loading
  - test/: test utilities, mocks, and integration tests
- Group code by feature when it improves clarity and cohesion.
- Keep logic decoupled from framework-specific code.

### Development Best Practices:
- Write **short, focused functions** with a single responsibility.
- Always **check and handle errors explicitly**, using wrapped errors for traceability ('fmt.Errorf("context: %w", err)').
- Avoid **global state**; use constructor functions to inject dependencies.
- Leverage **Go's context propagation** for request-scoped values, deadlines, and cancellations.
- Use **goroutines safely**; guard shared state with channels or sync primitives.
- **Defer closing resources** and handle them carefully to avoid leaks.

### Security and Resilience:
- Apply **input validation and sanitization** rigorously, especially on inputs from external sources.
- Use secure defaults for **JWT, cookies**, and configuration settings.
- Isolate sensitive operations with clear **permission boundaries**.
- Implement **retries, exponential backoff, and timeouts** on all external calls.
- Use **circuit breakers and rate limiting** for service protection.
- Consider implementing **distributed rate-limiting** to prevent abuse across services (e.g., using Redis).

### Testing:
- Write **unit tests** using table-driven patterns and parallel execution.
- **Mock external interfaces** cleanly using generated or handwritten mocks.
- Separate **fast unit tests** from slower integration and E2E tests.
- Ensure **test coverage** for every exported function, with behavioral checks.
- Use tools like 'go test -cover' to ensure adequate test coverage.

### Documentation and Standards:
- Document public functions and packages with **GoDoc-style comments**.
- Provide concise **READMEs** for services and libraries.
- Maintain a 'CONTRIBUTING.md' and 'ARCHITECTURE.md' to guide team practices.
- Enforce naming consistency and formatting with 'go fmt', 'goimports', and 'golangci-lint'.

### Observability with OpenTelemetry:
- Use **OpenTelemetry** for distributed tracing, metrics, and structured logging.
- Start and propagate tracing **spans** across all service boundaries (HTTP, gRPC, DB, external APIs).
- Always attach 'context.Context' to spans, logs, and metric exports.
- Use **otel.Tracer** for creating spans and **otel.Meter** for collecting metrics.
- Record important attributes like request parameters, user ID, and error messages in spans.
- Use **log correlation** by injecting trace IDs into structured logs.
- Export data to **OpenTelemetry Collector**, **Jaeger**, or **Prometheus**.

### Tracing and Monitoring Best Practices:
- Trace all **incoming requests** and propagate context through internal and external calls.
- Use **middleware** to instrument HTTP and gRPC endpoints automatically.
- Annotate slow, critical, or error-prone paths with **custom spans**.
- Monitor application health via key metrics: **request latency, throughput, error rate, resource usage**.
- Define **SLIs** (e.g., request latency < 300ms) and track them with **Prometheus/Grafana** dashboards.
- Alert on key conditions (e.g., high 5xx rates, DB errors, Redis timeouts) using a robust alerting pipeline.
- Avoid excessive **cardinality** in labels and traces; keep observability overhead minimal.
- Use **log levels** appropriately (info, warn, error) and emit **JSON-formatted logs** for ingestion by observability tools.
- Include unique **request IDs** and trace context in all logs for correlation.

### Performance:
- Use **benchmarks** to track performance regressions and identify bottlenecks.
- Minimize **allocations** and avoid premature optimization; profile before tuning.
- Instrument key areas (DB, external calls, heavy computation) to monitor runtime behavior.

### Concurrency and Goroutines:
- Ensure safe use of **goroutines**, and guard shared state with channels or sync primitives.
- Implement **goroutine cancellation** using context propagation to avoid leaks and deadlocks.

### Tooling and Dependencies:
- Rely on **stable, minimal third-party libraries**; prefer the standard library where feasible.
- Use **Go modules** for dependency management and reproducibility.
- Version-lock dependencies for deterministic builds.
- Integrate **linting, testing, and security checks** in CI pipelines.

### Key Conventions:
1. Prioritize **readability, simplicity, and maintainability**.
2. Design for **change**: isolate business logic and minimize framework lock-in.
3. Emphasize clear **boundaries** and **dependency inversion**.
4. Ensure all behavior is **observable, testable, and documented**.
5. **Automate workflows** for testing, building, and deployment.

---

## Raylz Relayer - Project-Specific Guidelines

### Project Overview:
This is a **cross-chain relayer** enabling communication between EVM blockchains via a Private Network Hub intermediary. Supports encrypted message transfers, Enygma DvP operations, and privacy protocols (Enygma).

### Domain Terminology:
- **PNH** = Private Network Hub (intermediary blockchain, formerly "Commit Chain")
- **PNo** / **Privacy Node** = Privacy Node (source/destination blockchain, formerly "Private Ledger")
- **CTS** = current name for the service that does cryptographic key management (formerly **KOS** / Key Operation Service)
- **AS** = Atomic Service (atomic transaction handling)
- **Enygma DvP** = Enygma Delivery versus Payment protocol (formerly "ZkDVP")
- **Enygma** = Privacy/encryption protocol
- **Public Chain** = Public blockchain (formerly "Public PL" / "PC")

### Service Applications:
- `private-relayer/` - Privacy Node relayer (includes atomic operations)
- `public-relayer/` - Public chain interactions
- `cts` - Cryptographic key management service (code in `cts/`, built via `cts/Dockerfile`; formerly KOS)

### Build Commands:
```bash
make build-relayer    # Build main relayer
make build-cts        # Build CTS binary
make proto-cts        # Regenerate CTS protobuf stubs
./start_dev.sh        # Local development environment
```

### Testing:
```bash
make test                 # Run tests
make test-coverage        # Coverage report
make test-coverage-xml    # XML coverage
```
- Mocks generated with `moq` via `go:generate` directives
- Integration tests use `testcontainers-go` for PostgreSQL
- **Fault Injection (resilience/chaos testing):** the `faultinjector/` package lets E2E tests arm `crash`/`panic`/`sleep`/`error` rules at named code points over HTTP — double-gated off in production by the `-tags faultinjection` build tag and the `FAULT_INJECTION_ENABLED` flag (no-op `noop.go` otherwise). Reach for it when debugging crash/restart recovery, analyzing idempotency or partial-batch behavior, or auditing fault tolerance. Reference: `faultinjector/README.md` (canonical) and `../rayls-privacy-docs-internal/docs/build/advanced/fault-injection.md`; the TypeScript client (`../rayls-privacy-tests-automation/src/utils/fault-injector.ts`) and resilience suite (`../rayls-privacy-tests-automation/test/e2e/security/resilience/`) drive it from E2E tests.

### Core Architectural Patterns:
1. **Listener-Executor**: Events listened → batched → executed
2. **Batcher[T]**: Generic time/count-based batching for efficiency
3. **Repository Pattern**: PostgreSQL with pgxpool connection pooling
4. **State Machine**: `TransactionState` enum tracks message lifecycle
5. **Key Queue**: Private keys queued per infrastructure (PrivateNode, PrivateHub, DvpOperator)

### Error Handling:
- **At external call origins** (DB queries, blockchain calls, HTTP/RPC, crypto operations): use `withstack.Wrap(fmt.Errorf("context: %w", err))` to capture stack trace at the point the error enters the system
- **Everywhere else** (services, handlers, business logic, bubbling up errors): use `fmt.Errorf("context: %w", err)` to add context without redundant stack traces
- **Business logic errors** (validation, domain rules): use `fmt.Errorf("descriptive message")` — no wrapping needed
- **Typed contract client errors without an inner cause**: use `New*Error("msg")`, not `WrapIn*Error("msg", nil)`
- **Rule**: "Am I the first code to touch this error from an external system?" → `withstack.Wrap`. Otherwise → `fmt.Errorf`
- Write descriptive, unique error messages at every level so the context chain alone identifies the exact code path
- **Error comparison**: use `errors.Is()` for sentinel checks and `errors.As()` for type assertions — never compare errors with `==`/`!=` (except `err == nil`)

### Database (PostgreSQL):
Key tables:
- `transactions` - Core transaction state
- `enygma` / `enygma_history` / `enygma_checkpoint` - Privacy protocol
- `merkle_tree` - Zero-knowledge merkle data
- `atomic_status` - Atomic operation states
- `dvp_swap` / `dvp_deposit` / `dvp_swap_metadata` - Enygma DvP operations
- `last_processed_block` - Event listener checkpoints
- `resource_lock` - Resource locking for concurrent access
- `calldata_signature` - Calldata signatures

### Configuration:
- Environment variables via `spf13/viper`
- Validation with `go-playground/validator`
- Key prefixes: `PNH_*`, `PRIVACY_NODE_*`, `PUBLIC_CHAIN_*`, `BLOCKCHAIN_DATABASE_*`, `BLOCKCHAIN_KMS_*`, `OTEL_*`
- Legacy prefixes (`COMMITCHAIN_*`, `BLOCKCHAIN_CHAIN*`) are supported via backward-compatible Viper bindings in `configinit/init.go`
- Run with: `./relayer-app run --env path/to/.env`
- **`RAYLS_AXYL_EPOCH_DURATION_SECS`** (dev-only) — overrides the axyl chain's genesis `epoch-duration-in-secs`. Default 86400 (24h). Lower to 60 to trigger frequent epoch transitions (used to drive proposer-task-respawn investigations). Baked into chain genesis; requires `start_dev.sh --clean`. Equivalent CLI: `start_dev.sh --epoch-duration N`.

### Contract Bindings:
- Generated bindings in `/contracts/` - DO NOT manually edit
- Wrapper clients in `/contractclient/`
- Contract addresses loaded at runtime from proxy registry

### Logging & Observability:
- Use `log/slog` structured logging (never fmt.Print)
- OpenTelemetry integration in `/otel` and `/logger`
- Include trace context in all logs
- JSON-formatted logs for observability tools

### Key Dependencies:
- `go-ethereum` - Ethereum client
- `spf13/cobra` + `spf13/viper` - CLI and config
- `pgx/v5` - PostgreSQL driver with connection pooling
- `golang-migrate` - Database migrations
- `cockroachdb/errors` - Stack traces at external call origins (via `withstack` wrapper)
- `testcontainers-go` - Integration testing
- `nats.go` - Inter-service messaging

### Opening Issues:
- When opening a GitHub issue, follow the templates in `.github/ISSUE_TEMPLATE/` (`bug_report.md`, `feature_request.md`).
- Fill in every section; write `n/a` where it doesn't apply rather than deleting the section.
- Keep it concise — facts over prose. Paste only the relevant log lines, not entire logs.
- Set the correct service (`private-relayer` / `public-relayer` / `cts`) and component.
- Do **not** file security vulnerabilities as public issues — follow `.github/SECURITY.md` (private disclosure).
- `google.golang.org/grpc` - gRPC server/client (CTS communication)
- `aws-sdk-go-v2` / `cloud.google.com/go/kms` - Cloud KMS envelope encryption (CTS)

### CTS Architecture (`cts/`):
CTS consolidates all cryptographic operations. The relayer communicates with CTS via gRPC.

**Package layout:**
- `cts/cmd/` - CLI entrypoint (`run`, `migrate` commands)
- `cts/config/` - Configuration struct
- `cts/domain/` - Domain models (encrypted/decrypted key types)
- `cts/service/` - Business logic (KeysService, EncryptService)
- `cts/repo/` - PostgreSQL repositories + migrations
- `cts/grpc/` - gRPC handlers for Keys and Encrypt services
- `cts/http/` - HTTP handler (public addresses endpoint)
- `cts/keygen/` - Key generation (ECDSA, ML-KEM-768, BabyJubJub, Poseidon)
- `cts/crypto/` - Encryptor factory (plaintext/AWS/GCP KMS)
- `cts/client/` - Cloud KMS API clients
- `cts/gen/proto/` - Generated gRPC stubs
- `cts/proto/` - Proto source files

**Key types managed:** ECDSA (signing), ML-KEM-768 (key encapsulation), BabyJubJub (Enygma), Poseidon (DvP proofs)

**CTS database tables:** `ecdsa_keys`, `dh_keys`, `shared_secrets`, `enygma_keys`, `dvp_keys`, `enygma_self_secrets`

**Security:** Private keys never leave CTS process memory. Keys stored encrypted at rest via AWS/GCP KMS envelope encryption.
