package db

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewHealthyPostgresDBStub_HealthPaths(t *testing.T) {
	pg := NewHealthyPostgresDBStub()
	require.NotNil(t, pg.GetManager())
	require.True(t, pg.FastHealthCheck().IsHealthy)
	require.Equal(t, pg.FastHealthCheck(), pg.GetManager().FastHealthCheck())
	require.NoError(t, pg.Close())
}
