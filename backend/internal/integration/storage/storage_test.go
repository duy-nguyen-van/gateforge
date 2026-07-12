package storage

import (
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"

	"github.com/stretchr/testify/require"
)

func TestProvideStorageAdapter_InvalidProvider(t *testing.T) {
	cfg := &config.Config{StorageProvider: "invalid"}
	adapter, err := ProvideStorageAdapter(cfg)
	require.Error(t, err)
	require.Nil(t, adapter)
	require.Contains(t, err.Error(), "invalid storage provider")
}

func TestProvideStorageAdapter_GCSMissingCredentials(t *testing.T) {
	cfg := &config.Config{
		StorageProvider:        constants.StorageProviderGCS,
		GCSCredentialsJSONPath: "/does/not/exist/credentials.json",
		GCSBucket:              "bucket",
	}
	adapter, err := ProvideStorageAdapter(cfg)
	require.Error(t, err)
	require.Nil(t, adapter)
}

func TestProvideStorageAdapter_S3MinimalConfig(t *testing.T) {
	cfg := &config.Config{
		StorageProvider: constants.StorageProviderS3,
		S3Region:        "us-east-1",
		S3AccessKey:     "AKIATESTKEY",
		S3SecretKey:     "secret",
		S3Bucket:        "test-bucket",
	}
	adapter, err := ProvideStorageAdapter(cfg)
	require.NoError(t, err)
	require.NotNil(t, adapter)
}
