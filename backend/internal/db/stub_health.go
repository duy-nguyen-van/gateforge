package db

import "time"

// NewHealthyPostgresDBStub returns a PostgresDB whose FastHealthCheck reports healthy.
// It is intended for handler unit tests that exercise the database health success path.
func NewHealthyPostgresDBStub() *PostgresDB {
	return &PostgresDB{
		manager: &DatabaseManager{
			healthStatus: HealthStatus{
				IsHealthy: true,
				LastCheck: time.Now(),
			},
			metrics: &ConnectionMetrics{},
		},
	}
}
