package storage

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"mime/multipart"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/testutil"

	"github.com/stretchr/testify/require"
)

func writeTestServiceAccountJSON(t *testing.T) string {
	t.Helper()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	keyBytes := x509.MarshalPKCS1PrivateKey(key)
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: keyBytes})

	payload := map[string]string{
		"type":                        "service_account",
		"project_id":                  "test-project",
		"private_key_id":              "key-id",
		"private_key":                 string(keyPEM),
		"client_email":                "test@test-project.iam.gserviceaccount.com",
		"client_id":                   "123",
		"auth_uri":                    "https://accounts.google.com/o/oauth2/auth",
		"token_uri":                   "https://oauth2.googleapis.com/token",
		"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
	}
	data, err := json.Marshal(payload)
	require.NoError(t, err)

	path := filepath.Join(t.TempDir(), "credentials.json")
	require.NoError(t, os.WriteFile(path, data, 0o600))
	return path
}

func TestNewGCSAdapter_MissingCredentialsFile(t *testing.T) {
	t.Parallel()
	testutil.InitLogger()
	cfg := &config.Config{
		GCSCredentialsJSONPath: "/does/not/exist/credentials.json",
		GCSBucket:              "bucket",
	}
	adapter, err := NewGCSAdapter(cfg)
	require.Error(t, err)
	require.Nil(t, adapter)
	require.Contains(t, err.Error(), "failed to load GCS credentials")
}

func TestNewGCSAdapter_InvalidCredentialsJSON(t *testing.T) {
	t.Parallel()
	testutil.InitLogger()

	dir := t.TempDir()
	path := filepath.Join(dir, "credentials.json")
	require.NoError(t, os.WriteFile(path, []byte("{invalid-json"), 0o600))

	cfg := &config.Config{
		GCSCredentialsJSONPath: path,
		GCSBucket:              "bucket",
	}
	adapter, err := NewGCSAdapter(cfg)
	require.Error(t, err)
	require.Nil(t, adapter)
}

func TestNewGCSAdapter_InvalidServiceAccountJSON(t *testing.T) {
	t.Parallel()
	testutil.InitLogger()

	dir := t.TempDir()
	path := filepath.Join(dir, "credentials.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"client_email":"test@example.com","private_key":"invalid","project_id":"proj"}`), 0o600))

	cfg := &config.Config{
		GCSCredentialsJSONPath: path,
		GCSBucket:              "bucket",
	}
	adapter, err := NewGCSAdapter(cfg)
	require.Error(t, err)
	require.Nil(t, adapter)
	require.Contains(t, err.Error(), "failed to create GCS client")
}

func TestGCSAdapter_GetObjectURL(t *testing.T) {
	t.Parallel()
	adapter := &GCSAdapter{
		config: &config.Config{GCSBucket: "my-bucket"},
	}
	require.Equal(t, "https://storage.googleapis.com/my-bucket/object.txt", adapter.GetObjectURL("object.txt"))
}

func TestGCSAdapter_GetPresignedURL_MissingCredentialsFile(t *testing.T) {
	t.Parallel()
	testutil.InitLogger()
	adapter := &GCSAdapter{
		config: &config.Config{
			GCSBucket:               "my-bucket",
			GCSCredentialsJSONPath:  "/does/not/exist/credentials.json",
			GCSPresignedURLDuration: time.Hour,
		},
	}

	_, err := adapter.GetPresignedURL(context.Background(), "object.txt")
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to read service account JSON")
}

func TestGCSAdapter_GetPresignedURL_InvalidJSON(t *testing.T) {
	t.Parallel()
	testutil.InitLogger()

	dir := t.TempDir()
	path := filepath.Join(dir, "credentials.json")
	require.NoError(t, os.WriteFile(path, []byte("{invalid-json"), 0o600))

	adapter := &GCSAdapter{
		config: &config.Config{
			GCSBucket:               "my-bucket",
			GCSCredentialsJSONPath:  path,
			GCSPresignedURLDuration: time.Hour,
		},
	}

	_, err := adapter.GetPresignedURL(context.Background(), "object.txt")
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to parse service account JSON")
}

func TestGCSAdapter_GetPresignedURL_Success(t *testing.T) {
	t.Parallel()
	testutil.InitLogger()

	path := writeTestServiceAccountJSON(t)
	adapter := &GCSAdapter{
		config: &config.Config{
			GCSBucket:               "test-bucket",
			GCSCredentialsJSONPath:  path,
			GCSPresignedURLDuration: 30 * time.Minute,
		},
	}

	url, err := adapter.GetPresignedURL(context.Background(), "object.txt", 15*time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, url)
	require.Contains(t, url, "test-bucket")
	require.Contains(t, url, "object.txt")
}

func TestGCSAdapter_UploadFile_OpenError(t *testing.T) {
	t.Parallel()
	testutil.InitLogger()
	adapter := &GCSAdapter{
		config: &config.Config{GCSBucket: "my-bucket"},
	}

	_, err := adapter.UploadFile(context.Background(), &multipart.FileHeader{Filename: ""}, "key")
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to open file")
}

func TestGCSAdapter_UploadFiles_Error(t *testing.T) {
	t.Parallel()
	testutil.InitLogger()
	adapter := &GCSAdapter{
		config: &config.Config{GCSBucket: "my-bucket"},
	}

	_, err := adapter.UploadFiles(context.Background(), []*multipart.FileHeader{{Filename: ""}})
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to upload file")
}
