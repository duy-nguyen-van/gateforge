# API types & SDK

Hand-maintained SPA types live in [`../types.ts`](../types.ts). Prefer migrating call sites to the monorepo TypeScript SDK when possible.

## Official SDK

The frontend depends on [`@gateforge/sdk`](../../../sdk/typescript/) via a local `file:` dependency:

```json
"@gateforge/sdk": "file:../sdk/typescript"
```

Shared client wiring: [`../sdk.ts`](../sdk.ts) (`gateforge` singleton).

OpenAPI source of truth: [`api/openapi.yaml`](../../../api/openapi.yaml) at the monorepo root.

## Regenerate the SDK client

From the **monorepo root** (requires Docker):

```bash
make sdk-generate
```

That regenerates `sdk/typescript/src/generated/` from `api/openapi.yaml` and patches the fetch client for NodeNext. Then rebuild the SDK package:

```bash
cd sdk/typescript && npm run build
```

The SPA’s `npm run generate:api` script only documents this flow; it does not generate into `frontend/src/api/generated/`.
