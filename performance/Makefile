.PHONY: deps seed-user seed-codes seed-passkeys signer token token-hold passkey oidc-e2e password-login rss smoke help

PERF_DIR := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
ROOT := $(abspath $(PERF_DIR)/..)
DATA := $(PERF_DIR)/.data
K6 := k6
GO := go
export PERF_DATA_DIR := $(DATA)
export PERF_CODES_FILE := $(DATA)/codes.json
export PERF_PASSKEYS_FILE := $(DATA)/passkeys.json

-include $(PERF_DIR)/.env.bench
export

help:
	@echo "GateForge IAM performance benches"
	@echo "  make deps              # go mod download for performance tools"
	@echo "  make seed-user         # register bench user via API"
	@echo "  make seed-codes        # insert auth codes + write codes.json"
	@echo "  make seed-passkeys     # register passkey fixtures"
	@echo "  make signer            # run webauthn soft-authenticator :9091"
	@echo "  make token             # k6 token-issuance ramp (capacity probe)"
	@echo "  make token-hold        # k6 constant hold (publishable RPS)"
	@echo "  make passkey           # k6 passkey p95/p99 (signer must be up)"
	@echo "  make oidc-e2e          # k6 full OIDC journey"
	@echo "  make password-login    # k6 password login/refresh"
	@echo "  make rss               # idle RSS of running server (PERF_PID=...)"
	@echo "  make smoke             # short seed + tiny k6 sanity (needs server)"

deps:
	cd $(PERF_DIR) && $(GO) mod tidy

seed-user:
	cd $(PERF_DIR) && $(GO) run ./cmd/seed-user -base "$${PERF_BASE_URL:-http://127.0.0.1:3000}"

seed-codes: deps
	mkdir -p $(DATA)
	cd $(PERF_DIR) && $(GO) run ./cmd/seed-auth-codes -n $${PERF_SEED_CODES:-100000} -out $(DATA)/codes.json

seed-passkeys: deps
	mkdir -p $(DATA)
	cd $(PERF_DIR) && $(GO) run ./cmd/seed-passkeys -n $${PERF_PASSKEY_USERS:-20} -out $(DATA)/passkeys.json \
		-base "$${PERF_BASE_URL:-http://127.0.0.1:3000}" \
		-origin "$${WEBAUTHN_RP_ORIGIN:-http://localhost:3000}" \
		-rp-id "$${WEBAUTHN_RP_ID:-localhost}"

signer: deps
	cd $(PERF_DIR) && $(GO) run ./cmd/webauthn-signer -fixtures $(DATA)/passkeys.json -addr 127.0.0.1:9091

token:
	cd $(PERF_DIR)/k6 && $(K6) run -e PERF_CODES_FILE=$(DATA)/codes.json token-issuance.js

token-hold:
	cd $(PERF_DIR)/k6 && $(K6) run -e PERF_CODES_FILE=$(DATA)/codes.json token-hold.js

passkey:
	cd $(PERF_DIR)/k6 && $(K6) run -e PERF_PASSKEYS_FILE=$(DATA)/passkeys.json passkey-login.js

oidc-e2e:
	cd $(PERF_DIR)/k6 && $(K6) run oidc-e2e.js

password-login:
	cd $(PERF_DIR)/k6 && $(K6) run password-login.js

rss:
	chmod +x $(PERF_DIR)/scripts/measure-idle-rss.sh
	$(PERF_DIR)/scripts/measure-idle-rss.sh

smoke: seed-user
	mkdir -p $(DATA)
	cd $(PERF_DIR) && $(GO) run ./cmd/seed-auth-codes -n 2000 -out $(DATA)/codes.json
	cd $(PERF_DIR)/k6 && $(K6) run -e PERF_CODES_FILE=$(DATA)/codes.json token-smoke.js
	@echo "smoke OK"
