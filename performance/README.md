# Performance benches

k6 load scripts and Go seed/signer tools for honest, reproducible single-node numbers (landing-page claims). Latest measured figures: [`RESULTS.md`](RESULTS.md).

## What this measures

| Landing claim | Publishable command | Metric |
| ------------- | ------------------- | ------ |
| Auth token issuance RPS | `make -C performance token-hold` | Sustained `http_reqs` ≈ `PERF_HOLD_RPS`, errors &lt; 1%, **p95 &lt; 200 ms** and **p99 &lt; 500 ms** |
| Idle RAM | `make -C performance rss` | Median **app** RSS only (not Postgres/Redis) |
| Passkey latency | `make -C performance signer` (other terminal) + `make -C performance passkey` | `passkey_e2e` **p95 &lt; 80 ms**, **p99 &lt; 150 ms** (multi-VU) |
| Complete OIDC journey | `make -C performance oidc-e2e` | `oidc_complete_logins` + `oidc_e2e` p95/p99 (secondary) |

Use `make -C performance token` only as a **capacity probe** (ramping arrival). Do **not** publish ramp peak RPS when p95/p99 are multi-second.

**p95 vs p99:** p95 is the marketing / landing figure; p99 is the capacity guardrail (rare stalls). Both must pass for a hold to be publishable.

## Honesty floor

- Do not count Postgres/Redis toward idle RAM.
- Do not label full OIDC logins/s as “token issuance RPS”.
- Do not publish 1-VU latency as marketing p95/p99.
- Document hardware; one app process = “single-node”.
- If measured figures beat / miss the landing claims, update landing copy to match [`RESULTS.md`](RESULTS.md).

## Prerequisites

- `k6`
- Go 1.24+
- Docker (Postgres/Redis via the backend compose stack)

Root shortcuts: `make performance-help`, `make performance-smoke`, `make performance-token-hold`, `make performance-passkey`, `make performance-oidc-e2e`, `make performance-rss`.

---

## How to run

All commands from the **repo root** unless noted.

### 1. Dependencies + env

```bash
make -C backend container-up
make -C backend migrate-up          # if DB is empty

cp performance/.env.bench.example performance/.env.bench
# Edit performance/.env.bench if needed (limits, PERF_HOLD_RPS, passwords)
```

`performance/Makefile` auto-loads `performance/.env.bench` when present.

### 2. Build and start the server under test

Prefer a lean API binary (embed SPA optional). **Merge** `backend/cmd/server/.env` with bench overrides so rate limits and WebAuthn origins apply:

```bash
cd backend && CGO_ENABLED=0 go build -o ../bin/gateforge-iam-server ./cmd/server

cd backend/cmd/server
set -a
source .env
source ../../../performance/.env.bench
set +a
# Critical: high limits (middleware reads these), localhost WebAuthn, no APM
export DEFAULT_RATE_LIMIT=1000000 AUTH_RATE_LIMIT=1000000
export WEBAUTHN_RP_ID=localhost WEBAUTHN_RP_ORIGINS=http://localhost:3000
export SENTRY_DSN= NEWRELIC_LICENSE= LOG_LEVEL=warn SERVE_EMBEDDED_FRONTEND=false

../../../bin/gateforge-iam-server
```

Leave that terminal running. Health check: `curl -fsS http://127.0.0.1:3000/api/v1/`.

### 3. Smoke (optional)

```bash
make -C performance smoke
```

Registers the bench user (with `PERF_TENANT_ID`), seeds 2k codes, runs a short `token-smoke.js` hold.

### 4. Token issuance (primary RPS)

```bash
make -C performance seed-user
PERF_SEED_CODES=500000 make -C performance seed-codes   # one-shot; takes a few minutes

# Publishable sustained RPS (constant arrival)
PERF_HOLD_RPS=350 PERF_HOLD_DURATION=45s PERF_CODE_OFFSET=0 \
  make -C performance token-hold

# Optional: ramp / cliff probe (not for marketing RPS)
PERF_HOLD_RPS=350 PERF_CODE_OFFSET=0 make -C performance token
```

Notes:

- Codes are **single-use**. After a run, bump `PERF_CODE_OFFSET` past codes already consumed, or re-seed.
- Raise `PERF_HOLD_RPS` until **p95** or **p99** approaches the thresholds (or errors appear); last good hold = published RPS.
- Thresholds in `performance/k6/token-hold.js`: `p(95)<200`, `p(99)<500`, errors &lt; 1%.

### 5. Passkey latency (p95 + p99)

```bash
make -C performance seed-passkeys   # needs WEBAUTHN_RP_ORIGINS matching -origin (default http://localhost:3000)

# Terminal A — soft authenticator
make -C performance signer          # listens on 127.0.0.1:9091

# Terminal B
PERF_PASSKEY_VUS=10 PERF_PASSKEY_DURATION=60s make -C performance passkey
```

Publish `passkey_e2e` **p95** (landing) and confirm **p99** from the same multi-VU run (`p(95)<80`, `p(99)<150`).

### 6. OIDC e2e + password login (secondary)

```bash
make -C performance oidc-e2e        # login → authorize → token → userinfo
make -C performance password-login  # /api/v1/login (+ refresh)
```

Uses `PERF_USER_EMAIL` / `PERF_USER_PASSWORD` (same user as token seed).

### 7. Idle RSS

Prefer a **fresh** process (restart the binary after heavy load so heap isn’t inflated):

```bash
# In another shell, after the server has been warm ~30s with no traffic:
PERF_PID=$(ps aux | awk '/bin\/gateforge-iam-server/ && !/awk|zsh/ {print $2; exit}')
PERF_PID=$PERF_PID make -C performance rss
```

Reads median RSS of that PID only → `.data/idle-rss.txt`.

---

## Make targets

| Target | What it does |
| ------ | ------------ |
| `make -C performance help` | List targets |
| `make -C performance deps` | `go mod tidy` for seed/signer tools |
| `make -C performance seed-user` | `POST /api/v1/register` (needs `PERF_TENANT_ID`) |
| `make -C performance seed-codes` | Insert auth codes + write `.data/codes.json` |
| `make -C performance seed-passkeys` | Register users + passkeys → `.data/passkeys.json` |
| `make -C performance signer` | WebAuthn assertion helper (`:9091`) |
| `make -C performance token-hold` | Constant-arrival token RPS (publish) |
| `make -C performance token` | Ramping token probe |
| `make -C performance smoke` | Tiny end-to-end sanity |
| `make -C performance passkey` | Passkey start→signer→finish |
| `make -C performance oidc-e2e` | Full OIDC journey |
| `make -C performance password-login` | App login/refresh |
| `make -C performance rss` | Idle RSS sample |

Useful env vars (see [`.env.bench.example`](.env.bench.example)):

| Variable | Meaning |
| -------- | ------- |
| `PERF_BASE_URL` | API base (default `http://127.0.0.1:3000`) |
| `PERF_HOLD_RPS` | Target arrival rate for `token-hold` |
| `PERF_HOLD_DURATION` | Hold length (default in script `45s`) |
| `PERF_CODE_OFFSET` | Skip already-consumed codes in `codes.json` |
| `PERF_SEED_CODES` | How many codes to mint |
| `PERF_PASSKEY_VUS` / `PERF_PASSKEY_DURATION` | Passkey load shape |
| `PERF_TENANT_ID` | Required for register (default seed tenant UUID) |
| `DEFAULT_RATE_LIMIT` / `AUTH_RATE_LIMIT` | Must be high on the **server** process |

## Layout

```
performance/
├── README.md           # this file — how it works and how to run
├── RESULTS.md          # latest publishable numbers
├── Makefile            # seed, k6, and RSS targets
├── .env.bench.example  # bench env template
├── cmd/                # Go helpers (seed user/codes/passkeys, WebAuthn signer)
├── k6/                 # load scripts
└── scripts/            # idle RSS sampler
```

## Tuning

1. Raise `PERF_HOLD_RPS` until latency (p95 **or** p99) or errors cliff; re-seed / bump offset between runs.
2. Raise `DATABASE_MAX_OPEN_CONNS` (and Redis) if the app waits on pools.
3. Keep k6 and the server on localhost; set `GOMAXPROCS` to physical cores for RPS showpieces.
4. Idle RSS: warm connect only, `LOG_LEVEL=warn`, empty APM DSNs, measure before heavy load (or restart).

## Recording results

Copy machine specs (CPU, RAM, OS), commit SHA, Postgres/Redis versions, and k6 summaries (include **p95 and p99**) into [`RESULTS.md`](RESULTS.md). Prefer `token-hold` + passkey p95/p99 + idle RSS as the three landing numbers.
