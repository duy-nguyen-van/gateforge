# GateForge Go SDK

Handwritten helpers for GateForge IAM (`package gateforge`). Generated OpenAPI models and API clients live in [`openapi/`](openapi/) (regenerate with `make sdk-generate`).

## Install

```bash
go get github.com/gateforge-iam/gateforge-iam/sdk/go@sdk/go/v0.1.0
```

## Usage

```go
package main

import (
	"context"
	"fmt"
	"log"

	gateforge "github.com/gateforge-iam/gateforge-iam/sdk/go"
	"github.com/gateforge-iam/gateforge-iam/sdk/go/openapi"
)

func main() {
	ctx := context.Background()
	client := gateforge.NewClient("https://iam.example.com",
		gateforge.WithTokenProvider(func(context.Context) (string, error) {
			return "ACCESS_TOKEN", nil
		}),
	)

	health, _, err := client.GetHealth(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(health.GetData().GetStatus())

	me, _, err := client.GetMe(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(me.GetData().GetEmail())

	// Or call generated APIs directly:
	_ = client.API().AdminAPI
	_ = openapi.LoginRequest{}
}
```

## OIDC (server-side)

```go
pkce, err := gateforge.PKCEGenerate()
if err != nil {
	log.Fatal(err)
}
authorizeURL, err := gateforge.BuildAuthorizeURL("https://iam.example.com", gateforge.AuthorizeParams{
	ClientID:      "my-app",
	RedirectURI:   "https://app.example.com/callback",
	Scope:         "openid profile",
	CodeChallenge: pkce.Challenge,
})
// After the browser returns ?code=...
tok, _, err := client.ExchangeAuthorizationCode(ctx, "my-app", "https://app.example.com/callback", code, pkce.Verifier, "")

// Resource server: introspect access or refresh tokens (confidential client required).
active, _, err := client.IntrospectToken(ctx, "rs-client", "rs-secret", tok.GetAccessToken(), "access_token")
if err != nil {
	log.Fatal(err)
}
fmt.Println(active.GetActive(), active.GetSub())
```

WebAuthn ceremonies require a browser authenticator — the Go SDK exposes generated start/finish request types only.

## Errors

App API failures typically use the Meta envelope (`error_code`, `message`, `code`). Generated transport errors are `*openapi.GenericOpenAPIError` (see `AsOpenAPIError`).

## Regenerate

```bash
make sdk-generate   # from monorepo root (Docker required)
make sdk-go-test
```
