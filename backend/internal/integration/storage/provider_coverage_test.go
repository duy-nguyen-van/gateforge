package storage

import (
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"

	"github.com/stretchr/testify/require"
)

func TestProvideStorageAdapter_S3InitFailure(t *testing.T) {
	cfg := &config.Config{
		StorageProvider: constants.StorageProviderS3,
		S3Region:        "us-east-1",
		S3AccessKey:     "",
		S3SecretKey:     "",
		S3Bucket:        "",
	}
	adapter, err := ProvideStorageAdapter(cfg)
	if err != nil {
		require.Nil(t, adapter)
		require.Contains(t, err.Error(), "Failed to initialize S3 storage adapter")
		return
	}
	require.NotNil(t, adapter)
}

func TestProvideStorageAdapter_GCSInitFailure(t *testing.T) {
	cfg := &config.Config{
		StorageProvider:        constants.StorageProviderGCS,
		GCSCredentialsJSONPath: "/does/not/exist.json",
		GCSBucket:              "bucket",
	}
	adapter, err := ProvideStorageAdapter(cfg)
	require.Error(t, err)
	require.Nil(t, adapter)
	require.Contains(t, err.Error(), "Failed to initialize GCS storage adapter")
}
