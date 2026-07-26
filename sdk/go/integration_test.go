//go:build integration

package gateforge_test

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	gateforge "github.com/gateforge-iam/gateforge-iam/sdk/go"
)

// Integration tests against a live GateForge IAM instance.
// Example:
//
//	GATEFORGE_BASE_URL=http://localhost:3000 go test -tags=integration ./...
func TestIntegration_Health(t *testing.T) {
	base := os.Getenv("GATEFORGE_BASE_URL")
	if base == "" {
		t.Skip("GATEFORGE_BASE_URL not set")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	client := gateforge.NewClient(base, gateforge.WithHTTPClient(&http.Client{Timeout: 10 * time.Second}))
	env, resp, err := client.GetHealth(ctx)
	if err != nil {
		t.Fatalf("GetHealth: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d", resp.StatusCode)
	}
	data := env.GetData()
	if data.GetStatus() == "" {
		t.Fatalf("empty health status: %+v", env)
	}
}

func TestIntegration_OpenIDConfiguration(t *testing.T) {
	base := os.Getenv("GATEFORGE_BASE_URL")
	if base == "" {
		t.Skip("GATEFORGE_BASE_URL not set")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	client := gateforge.NewClient(base)
	cfg, resp, err := client.GetOpenIDConfiguration(ctx)
	if err != nil {
		t.Fatalf("GetOpenIDConfiguration: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d", resp.StatusCode)
	}
	if cfg.GetIssuer() == "" || cfg.GetTokenEndpoint() == "" {
		t.Fatalf("incomplete discovery: %+v", cfg)
	}
}
