# Fault Injection Framework

Runtime fault injection for resilience and chaos testing of the private relayer. The
framework lets E2E tests arm rules at named points in the relayer code via an HTTP API,
then watch how the relayer behaves when those points trigger a crash, panic, sleep, or
error.

> **Production safety.** This package is gated behind a build tag (`-tags faultinjection`)
> *and* a runtime config flag (`FAULT_INJECTION_ENABLED`). The default `make build-relayer`
> produces a binary in which every `Check()` call is a no-op (see `noop.go`) — the rule
> map, HTTP server, persistence, and `os.Exit()` paths are not even compiled in. Never
> ship a binary built with `-tags faultinjection` to production.

---

## Why this exists

The relayer connects external systems with very different failure modes — Postgres, NATS
JetStream, the Private Hub (CC), the Privacy Node (PL), the KOS, the Proofs API. Each
system can fail independently and at any point. The relayer must survive:

- crashes between a write to one system and a write to another (e.g., DB written, NATS
  not yet acked → message redelivers on restart),
- a stuck or slow downstream (do retries cause double effects?),
- partial-success batches (which messages survived the half-completed run?),
- replay of the same input through different code paths (idempotency).

These are notoriously hard to reproduce with synthetic test inputs alone. Fault injection
gives tests a controlled way to **stop the relayer at an exact line of code** and observe
the system's recovery behaviour. The result: a deterministic, scripted reproduction of
exactly the kind of failure that's otherwise impossible to write a unit test for.

---

## Sessions — isolation between parallel tests

Every test owns a **session**. Sessions are created server-side (one HTTP call), live
for the duration of the test, and are dropped in the test's `after()` hook. Rules,
trigger logs, and counters live inside their session. Two parallel tests on the same
relayer can each arm rules on the same point without overwriting each other's state.

```ts
import { FaultInjector, FaultSession, FAULT_POINTS } from '.../fault-injector';

let fi: FaultInjector;
let session: FaultSession;

before(async () => {
  fi = FaultInjector.forRelayer('B');
  session = await fi.newSession();
});

after(async () => {
  await session.clear();
});

it('observes a transient KOS error', async () => {
  await session.arm({ point: FAULT_POINTS.SOME_POINT, action: 'error', one_shot: true });

  // ... trigger the relayer flow that should hit the armed point ...

  // MANDATORY: confirm the relayer is alive — a neighbour session may have armed
  // a `crash` rule on the same point that fired alongside ours.
  await fi.waitUntilAlive(60_000);
  expect(await session.wasTriggered(FAULT_POINTS.SOME_POINT)).to.equal(true);
});
```

Or, equivalently, using the bundled convenience that combines the trigger-poll and the
liveness gate into one call:

```ts
await session.arm({ point: FAULT_POINTS.SOME_POINT, action: 'error', one_shot: true });
// trigger the relayer flow
await session.assertLiveAfter(FAULT_POINTS.SOME_POINT, 60_000);
```

### Authoring contract — mandatory for every FI-using test

After **any** operation that arms a rule, observes a trigger, or otherwise expects the
relayer to make progress, the test **must** confirm service health before proceeding.
This is non-negotiable because a neighbour session on the same relayer may have armed
`crash` or `panic`, and that neighbour's terminal could fire on the same `Check()` your
sleep/error rule fired on. Tests that arm only `sleep` are *not* exempt — they share the
relayer with everything else. Tests run in parallel by default; neighbour-tolerance is a
test-authoring requirement, not opt-in/opt-out.

`fi.waitUntilAlive(timeoutMs)` polls the FI HTTP endpoint. Combine with business-level
readiness (e.g., NATS consumer reattached, expected log lines) if your test depends on it.

---

## Multi-arm semantics — what happens when several sessions arm the same point

Two arms at the same point are considered *equivalent* (fire together as one event) if
they share both action type and identifying payload:

| Action  | What makes two arms equivalent                                                                                          |
|---------|-------------------------------------------------------------------------------------------------------------------------|
| `sleep` | Action alone. A 10s sleep encompasses a 5s sleep; both sessions get latency.                                            |
| `error` | Action **and** identical key — `error_code` when set on both, else `message`. The key is namespaced by source, so arms keyed on `error_code="x"` and `message="x"` form **two** classes. |
| `panic` | Same keying rule as `error`.                                                                                            |
| `crash` | Action alone. There is no `message` or `error_code` to disambiguate crashes.                                            |

For each `Check(point)` invocation, the framework:

1. **Sleep**: if any sleep arms exist, sleeps for `max(durations)`. Every sleep arm
   decrements / logs / consumes-if-one-shot — they're all satisfied by the larger delay.
2. **Pick a terminal** by priority (most-destructive-wins): `crash` > `panic` > `error`.
   - `crash`: every crash arm across all sessions decrements (one equivalence class).
   - `panic`: among panic arms, FIFO-oldest *message class* fires. Every arm in that
     class decrements. Arms with other messages stay armed for the next `Check`.
   - `error`: among error arms, FIFO-oldest *message class* fires; **only that one arm**
     decrements (single error return slot). Other error arms stay armed.

This is "fair to everyone": an arm is marked fired *only when an event matching its
identity occurred*. A session that armed `panic "rpc timeout"` and saw the framework
panic with a different message has not observed its intended fault yet — its arm stays
for the next time the point is hit.

### Worked examples — copy-paste reference

Notation: `count=N` means counter starts at N; absence of `count` and `one_shot` means
default `count=0` = unlimited (never used up). "Stays armed" lists every arm whose
counter remains positive (or is unlimited) and whose `one_shot` wasn't consumed.

| # | Armed at point.X (across sessions) | Effective sleep | Terminal that fires | Counters after Check | Stays armed |
|---|---|---|---|---|---|
| 1 | sess1: `sleep 5s` count=3 | 5s | — | sess1.sleep: 3→2 | sess1.sleep (count=2) |
| 2 | sess1: `sleep 5s` (unlimited); sess2: `sleep 10s` one_shot | 10s | — | sess1.sleep unchanged; sess2.sleep consumed | sess1.sleep |
| 3 | sess1: `error "kos"` count=2; sess2: `error "rpc"` count=2 | 0 | `"kos"` (FIFO-oldest distinct class) | sess1.error: 2→1; sess2.error: unchanged | sess1.error (count=1), sess2.error (count=2) |
| 4 | sess1: `error "kos"` count=2; sess2: `error "kos"` count=2 *(same message)* | 0 | `"kos"` (single class — both arms equivalent) | sess1.error: 2→1; sess2.error: 2→1 | sess1.error (count=1), sess2.error (count=1) |
| 5 | sess1: `sleep 5s`; sess2: `error "kos"` count=2 | 5s | `"kos"` | sess1.sleep unchanged; sess2.error: 2→1 | sess1.sleep, sess2.error (count=1) |
| 6 | sess1: `sleep 10s`; sess2: `panic "X"`; sess3: `panic "Y"` *(different messages)* | 10s | panic `"X"` (FIFO-oldest distinct class) | sess1.sleep unchanged; sess2.panic unchanged; sess3.panic unchanged | sess1.sleep, sess2.panic, sess3.panic |
| 7 | sess1: `sleep 10s`; sess2: `panic "X"`; sess3: `panic "X"` *(same message — single class)* | 10s | panic `"X"` (single class) | sess1.sleep unchanged; sess2.panic unchanged; sess3.panic unchanged | all three |
| 8 | sess1: `error "a"`; sess2: `panic "b"`; sess3: `crash` count=2 | 0 | crash | sess1.error unchanged; sess2.panic unchanged; sess3.crash: 2→1 | sess1.error, sess2.panic, sess3.crash (count=1) |
| 9 | sess1: `sleep 5s`; sess2: `error "e"`; sess3: `panic "p"`; sess4: `crash` count=2 | 5s | crash | sess1.sleep unchanged; sess4.crash: 2→1; sess2.error & sess3.panic untouched | sess1.sleep, sess2.error, sess3.panic, sess4.crash (count=1) |
| 10 | sess1: `crash` (unlimited); sess2: `crash` (unlimited); sess3: `crash` count=3 *(all crash = one equiv class)* | 0 | crash | sess1 & sess2 unchanged; sess3.crash: 3→2 | all three |
| 11 | sess1: `error "first"`; sess2: `error "second"`; sess3: `error "third"` *(all one_shot, distinct messages → 3 classes)* | 0 | `"first"` (FIFO-oldest) | sess1.error consumed; others unchanged | sess2.error, sess3.error |

Read the last row carefully — it shows the "errors over consecutive Check calls"
pattern. Three Checks in a row return `"first"`, `"second"`, `"third"` in that order,
each consuming its own one_shot arm exactly once.

---

## Configuration

```env
FAULT_INJECTION_ENABLED=true
FAULT_INJECTION_PORT=6660
FAULT_INJECTION_PERSIST_PATH=/tmp/faultinjector_state.json
FAULT_INJECTION_SESSION_TTL_MINUTES=60     # optional, default 60
```

`FAULT_INJECTION_PERSIST_PATH` is required for crash-recovery tests — without it, armed
rules and sessions don't survive the crash.

`FAULT_INJECTION_SESSION_TTL_MINUTES` controls the inactivity TTL. A background sweeper
runs once per minute and drops any session whose `last_activity` timestamp is older than
the TTL. Default is 60 minutes.

### Dev-local port mapping

Each fault-injectable service (private relayer, public relayer, CTS) has a unique
fault-injection port in `docker-compose.dev-local.yml`. The defaults match the
per-instance `.env` files in `docker/development/local/`:

| Service       | FI Port | Host URL                |
|---------------|---------|-------------------------|
| relayer-a     | 6660    | `http://127.0.0.1:6660` |
| relayer-b     | 6661    | `http://127.0.0.1:6661` |
| relayer-c     | 6662    | `http://127.0.0.1:6662` |
| relayer-d     | 6663    | `http://127.0.0.1:6663` |
| relayer-e     | 6664    | `http://127.0.0.1:6664` |
| relayer-f     | 6665    | `http://127.0.0.1:6665` |
| pubrelayer-a  | 6700    | `http://127.0.0.1:6700` |
| pubrelayer-b  | 6701    | `http://127.0.0.1:6701` |
| pubrelayer-c  | 6702    | `http://127.0.0.1:6702` |
| pubrelayer-d  | 6703    | `http://127.0.0.1:6703` |
| pubrelayer-e  | 6704    | `http://127.0.0.1:6704` |
| pubrelayer-f  | 6705    | `http://127.0.0.1:6705` |
| cts-a         | 6800    | `http://127.0.0.1:6800` |
| cts-b         | 6801    | `http://127.0.0.1:6801` |
| cts-c         | 6802    | `http://127.0.0.1:6802` |
| cts-d         | 6803    | `http://127.0.0.1:6803` |
| cts-e         | 6804    | `http://127.0.0.1:6804` |
| cts-f         | 6805    | `http://127.0.0.1:6805` |

### Build tag

The `start_dev.sh` flow builds with `-tags faultinjection` (see `Dockerfile.relayer-dev`,
`Dockerfile.kos-dev`, `Dockerfile.public-relayer-dev`, and the matching `air.*.toml`
files). Production builds use `Dockerfile.relayer` / `Dockerfile.kos` /
`Dockerfile.public-relayer` without the tag — the no-op shim in `noop.go` is compiled
instead.

#### `BUILD_TAGS` mechanism

The dev Dockerfiles read a `BUILD_TAGS` build-arg / env-var and pass it through to both
the initial container build and `air`'s hot-reload rebuild. This lets you flip
fault-injection support on or off without editing the Dockerfile.

- **Build-arg form** (one-off build):
  ```bash
  docker build --build-arg BUILD_TAGS=faultinjection -f Dockerfile.relayer-dev .
  ```
- **Compose env form** (the path `start_dev.sh` uses):
  ```yaml
  # docker-compose.dev-local.yml — already wired for the relayer / kos /
  # public-relayer services that need FI.
  relayer-a:
    build:
      context: .
      dockerfile: Dockerfile.relayer-dev
      args: { BUILD_TAGS: faultinjection }
    environment:
      BUILD_TAGS: faultinjection   # air rebuilds inside the container pick this up
  ```
- **Local `go build`**:
  ```bash
  go build -tags faultinjection ./...     # full FI tree
  go build ./...                          # no tag → noop.go path
  ```

When `BUILD_TAGS=""` (the default) neither the binary nor the air rebuild loop pulls in
the fault-injection machinery, and `noop.go` is compiled — Check() is a no-op and the
HTTP control server is never started.

---

## HTTP API

| Method   | Path                                          | Action                                                                  |
|----------|-----------------------------------------------|-------------------------------------------------------------------------|
| `POST`   | `/sessions`                                   | Create a session. Returns `{ "id": "<uuid>", "created_at": "..." }`.    |
| `GET`    | `/sessions`                                   | List session summaries (debug / observability).                         |
| `GET`    | `/sessions/{sid}`                             | Full snapshot of one session: rules + log + metadata.                   |
| `DELETE` | `/sessions/{sid}`                             | Drop the session entirely.                                              |
| `POST`   | `/sessions/{sid}/faults`                      | Arm a rule in this session.                                             |
| `DELETE` | `/sessions/{sid}/faults`                      | Clear every rule in this session (keeps the session and its log).       |
| `DELETE` | `/sessions/{sid}/faults/{point}`              | Clear one rule by point (URL-encode if it contains slashes).            |
| `DELETE` | `/sessions/{sid}/log`                         | Empty the trigger log without dropping rules.                           |

The flat `/faults` endpoints from earlier revisions are no longer served — calls return
`404`.

### Arm a fault

```bash
SID=$(curl -s -X POST http://127.0.0.1:6661/sessions | jq -r .id)

curl -X POST http://127.0.0.1:6661/sessions/$SID/faults \
  -H 'Content-Type: application/json' \
  -d '{
    "point": "enygma.handler.Receiver.HandleEnygmaCrossTransfer.after_insert_history",
    "action": "crash",
    "one_shot": true
  }'
```

### FaultRule schema

| Field         | Type   | Notes                                                                                                                                                |
|---------------|--------|------------------------------------------------------------------------------------------------------------------------------------------------------|
| `point`       | string | **Required.** Fault point name. Must match a registered point in the running binary.                                                                 |
| `action`      | string | **Required.** One of `crash`, `sleep`, `panic`, `error`.                                                                                             |
| `duration_ms` | int    | For `sleep`: how long to block in milliseconds.                                                                                                      |
| `message`     | string | For `error`/`panic`: human-readable payload. Acts as the equivalence-class key when `error_code` is not set.                                         |
| `error_code`  | string | For `error`/`panic`: machine-readable discriminator the production cutpoint can switch on via `CodeOf(err)`. When set, also overrides `message` as the equivalence-class key. |
| `one_shot`    | bool   | If `true`, the rule is consumed after the first trigger. Atomic-with-trigger persistence applies.                                                    |
| `max_count`   | int    | If `>0`, the rule fires at most this many times. `0` = unlimited (default).                                                                          |

### Actions

| Action  | Behavior                                                                    | Process exits?           |
|---------|-----------------------------------------------------------------------------|--------------------------|
| `crash` | `os.Exit(1)` immediately after slog flushes                                 | yes (rc=1)               |
| `panic` | `panic(<message>)` — panics the goroutine; if not recovered, terminates     | yes if not recovered     |
| `sleep` | `time.Sleep(duration_ms)` — blocks the calling goroutine                    | no                       |
| `error` | Returns a Go error from `Check()` — caller decides what to do               | no                       |

**Choosing the action**
- Use `crash` to test NATS redelivery, recovery_data resume, and restart-driven paths.
- Use `error` to test error-handling without terminating the process.
- Use `sleep` to test timeout behaviour or widen race windows in concurrency tests.
- Use `panic` only when you specifically want to test panic-recovery machinery.

---

## Discriminating error types from Go (`error_code`)

Production code often handles different *kinds* of errors at the same line of code with
different behaviour — retry on a transient timeout, escalate on a permanent failure, back
off on a rate-limit. The `error_code` field lets a test arm a single cutpoint to drive
each of those branches in turn.

`action: "error"` returns a typed `*faultinjector.Error` value carrying both `Code` (the
test-supplied `error_code`, or `""` when none) and `Message`. Production callers
discriminate via `errors.As` or the `CodeOf(err)` one-liner:

```go
if err := faultinjector.Check("private_relayer.source.db.before_persist"); err != nil {
    switch faultinjector.CodeOf(err) {
    case "timeout":
        // transient: retry with backoff
    case "db_locked":
        // contention: yield this batch, take the next on the queue
    case "permanent_failure":
        // give up: surface upstream so the orchestrator reverts
    default:
        // no FI / unknown code: production-grade fallback
    }
    return fmt.Errorf("persist: %w", err)
}
```

Three arms at one cutpoint, each driving a different production branch:

```ts
// session A
await sessionA.arm({ point: P, action: 'error', error_code: 'timeout',           one_shot: true });
// session B
await sessionB.arm({ point: P, action: 'error', error_code: 'db_locked',         one_shot: true });
// session C
await sessionC.arm({ point: P, action: 'error', error_code: 'permanent_failure', one_shot: true });
```

Each call to `Check(P)` returns one error in FIFO order; each `error_code` forms its own
equivalence class. Tests can confirm *which* arm fired by reading the trigger log:

```ts
const events = await sessionA.triggerEvents(P);
expect(events[0].code).to.equal('timeout');  // the arm A observed
```

**Backward compatibility.** A rule armed without `error_code` returns an error whose
`.Error()` string is byte-identical to the legacy `"FAULT INJECTION at <point>:
<message>"` format. Existing tests that string-match on that prefix keep working.

---

## Persistence

When `FAULT_INJECTION_PERSIST_PATH` is set:

1. The entire session table is fsynced to disk on every mutation: `NewSession`,
   `SetRuleInSession`, `ClearRuleInSession`, `ClearRulesInSession`, `ClearLogInSession`,
   `ClearSession`, **and on rule consumption inside `Check()`** (one-shot or counter
   decrement).
2. On `Enable()`, the session table is restored from disk. If the file is missing,
   unreadable, or in an unexpected shape, a WARN is logged and the framework starts
   with an empty session table (no migration code; the file is transient test state).

**One-shot crash rules are safe** because consumption is persisted *before* `os.Exit`:

1. Test arms a `crash` one-shot rule for point P in session S.
2. Relayer reaches P, `Check()` decrements/consumes the rule, persists the updated
   session table to disk, releases the lock, then calls `os.Exit(1)`.
3. Relayer docker container is restarted. **Not automatic in the local dev environment**
   — the test must explicitly bring the container back up (e.g. via `compose.start`)
   and wait for the relayer to become healthy. Cloud deployments (dev/qa/prod) restart
   automatically.
4. On startup, persisted state shows the one-shot rule is gone → `Check(P)` returns nil
   → the retry path completes.

Persistent (non-one-shot) rules survive restarts and continue firing — useful for
"keep this point failing across N restarts" scenarios.

---

## Naming convention

Fault points use `<go_import_path>.<Type>.<Method>.<position>`. Position names are verb
phrases: `before_X`, `after_X`, `before_X_after_Y`. Be specific —
`after_insert_history` beats `after_db`.

---

## Adding a new fault point

1. **Identify the boundary.** Read the code near a write to an external system. Where
   would a crash leave you in an inconsistent state?
2. **Pick a name** following the convention above.
3. **Add the import** in the target Go file:
   ```go
   import "github.com/raylsnetwork/rayls-sovereign-relayer/faultinjector"
   ```
4. **Insert the call** at the cut point and handle the returned error appropriately:
   ```go
   // Boundary: after DB row is committed, before chain TX is broadcast.
   // A crash here exercises the recovery-resume path on restart.
   if err := faultinjector.Check(
       "enygma.service.EnygmaExecutor.persistAndBroadcast.after_insert_recovery_before_history",
   ); err != nil {
       return fmt.Errorf("fault injection at after_insert_recovery_before_history: %w", err)
   }
   ```
   Behaviour:
   - With `crash` armed → `os.Exit(1)` runs *inside* `Check()`; the line after never runs.
   - With `error` armed → `Check()` returns the configured error; the caller decides.
   - With no rule armed → `Check()` returns `nil` immediately; zero behavioural impact.
5. **Document the cut point in code.** Leave a brief comment immediately above the
   `Check()` call describing the boundary it sits on and the failure mode a crash there
   exposes. The comment is the local documentation for the point.
6. **Expose the constant to TypeScript callers** in the test-automation repo's
   `fault-injector.ts` helper, alongside the other point names.

The `Check()` call is a no-op when fault injection is disabled (single mutex
acquisition, zero behavioural impact in production with the default build tag).

---

## TypeScript Client

The client is split between a `FaultInjector` (per-relayer health + session factory) and
a `FaultSession` (per-test isolation).

```typescript
import { FaultInjector, FaultSession, FAULT_POINTS } from './fault-injector';
import { compose } from './docker-compose';

// Per-relayer client. Pass a host override (e.g. '127.0.0.1') to talk from outside
// the Docker network.
const fi = FaultInjector.forRelayer('B');

// One session per test file.
const session: FaultSession = await fi.newSession();

// Arm a one-shot crash at a fault point.
await session.arm({ point: FAULT_POINTS.SOME_POINT, action: 'crash', one_shot: true });

// ... trigger the relayer flow that should hit the armed point ...

// Start the container back up. In local dev, docker-compose does NOT auto-restart
// crashed containers (intentional, so a crash can be inspected). The test must
// explicitly bring it back. In cloud deployments (dev/qa/prod) the orchestrator
// would restart automatically. This step is a local dev environment.
compose.start('relayer-b');

// Wait for the fault-injection HTTP API to answer again. This is a proxy for
// "service is up and the binary has finished initializing". For full
// business-level readiness assert specific operational signals as well
// (NATS consumer reattached, expected log lines emitted, etc.).
await fi.waitUntilAlive(60_000);

// Per-session trigger inspection — semantics documented inline in the TS file.
await session.wasTriggered(FAULT_POINTS.SOME_POINT);    // boolean
await session.triggerCount(FAULT_POINTS.SOME_POINT);    // exact count
await session.triggerEvents();                          // full log
await session.triggerEvents(FAULT_POINTS.SOME_POINT);   // filtered log
await session.clearLog();                                // empties log, keeps rules
await session.clearAllRules();                           // empties rules, keeps log

// MANDATORY teardown: drop the session entirely. Per-test isolation depends on it.
// `clear()` swallows 404 so it's idempotent — call it from `after()` / `afterAll()`.
await session.clear();
```

Pair `FaultInjector` with the docker-compose helper from the same `src/utils/`
directory when a test needs to restart the container after a `crash` action.
