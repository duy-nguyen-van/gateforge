package services

import (
	"strings"

	"github.com/gateforge-iam/gateforge-iam/internal/config"
)

// FederationCallbackURL returns the OAuth redirect URI registered with upstream IdPs.
func FederationCallbackURL(cfg *config.Config, providerID string) string {
	base := strings.TrimSuffix(strings.TrimSpace(cfg.AppBaseURL), "/")
	if base == "" {
		base = "http://localhost:3000"
	}
	return base + "/oidc/federation/" + strings.ToLower(strings.TrimSpace(providerID)) + "/callback"
}

func federationEncryptionKey(cfg *config.Config) string {
	if cfg.MFAEncryptionKey != "" {
		return cfg.MFAEncryptionKey
	}
	return cfg.JWTSecret
}
