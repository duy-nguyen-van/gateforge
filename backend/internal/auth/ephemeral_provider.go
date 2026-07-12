package auth

import (
	"github.com/gateforge-iam/gateforge-iam/internal/cache"
	"github.com/gateforge-iam/gateforge-iam/internal/config"
)

// ProvideEphemeralStore wires Redis-backed ephemeral WebAuthn and MFA state.
func ProvideEphemeralStore(c cache.Cache, cfg *config.Config) *EphemeralStore {
	return NewEphemeralStore(c, cfg)
}
