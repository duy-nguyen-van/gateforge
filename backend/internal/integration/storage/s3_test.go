package storage

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/testutil"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/require"
)

func newMockS3Adapter(t *testing.T, handler http.HandlerFunc) *S3Adapter {
	t.Helper()
	testutil.InitLogger()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	cfg, err := awsconfig.LoadDefaultConfig(
		context.Background(),
		awsconfig.WithRegion("us-east-1"),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("AKIATESTKEY", "secret", "")),
	)
	require.NoError(t, err)

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(server.URL)
		o.UsePathStyle = true
	})

	return &S3Adapter{
		config: &config.Config{
			S3Region:               "us-east-1",
			S3Bucket:               "test-bucket",
			S3PresignedURLDuration: time.Hour,
		},
		client: client,
		bucket: "test-bucket",
	}
}

func newMultipartFileHeader(t *testing.T, filename, contentType string, data []byte) *multipart.FileHeader {
	t.Helper()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile("file", filename)
	require.NoError(t, err)
	_, err = part.Write(data)
	require.NoError(t, err)
	require.NoError(t, writer.Close())

	reader := multipart.NewReader(&buf, writer.Boundary())
	form, err := reader.ReadForm(10 << 20)
	require.NoError(t, err)
	t.Cleanup(func() { _ = form.RemoveAll() })

	return form.File["file"][0]
}

func TestNewS3Adapter_Success(t *testing.T) {
	t.Parallel()
	testutil.InitLogger()
	cfg := &config.Config{
		S3Region:    "us-east-1",
		S3AccessKey: "AKIATESTKEY",
		S3SecretKey: "secret",
		S3Bucket:    "test-bucket",
	}
	adapter, err := NewS3Adapter(cfg)
	require.NoError(t, err)
	require.NotNil(t, adapter)
}

func TestS3Adapter_GetObjectURL(t *testing.T) {
	t.Parallel()
	adapter := &S3Adapter{
		config: &config.Config{S3Region: "us-east-1"},
		bucket: "my-bucket",
	}
	require.Equal(t, "https://my-bucket.s3.us-east-1.amazonaws.com/object.txt", adapter.GetObjectURL("object.txt"))
}

func TestS3Adapter_GetPresignedURL(t *testing.T) {
	t.Parallel()
	testutil.InitLogger()
	adapter := newMockS3Adapter(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	url, err := adapter.GetPresignedURL(context.Background(), "file.txt", 30*time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, url)
	require.Contains(t, url, "file.txt")
}

func TestS3Adapter_UploadFile_Success(t *testing.T) {
	t.Parallel()
	var uploaded []byte
	adapter := newMockS3Adapter(t, func(w http.ResponseWriter, r *http.Request) {
		uploaded, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
	})

	file := newMultipartFileHeader(t, "upload.txt", "text/plain", []byte("file-data"))
	result, err := adapter.UploadFile(context.Background(), file, "uploads/upload.txt")
	require.NoError(t, err)
	require.Equal(t, "file-data", string(uploaded))
	require.Equal(t, "uploads/upload.txt", result.Key)
	require.Equal(t, "test-bucket", result.Bucket)
	require.Contains(t, result.URL, "uploads/upload.txt")
}

func TestS3Adapter_UploadFile_OpenError(t *testing.T) {
	t.Parallel()
	testutil.InitLogger()
	adapter := &S3Adapter{
		config: &config.Config{S3Region: "us-east-1"},
		client: s3.NewFromConfig(aws.Config{Region: "us-east-1"}),
		bucket: "test-bucket",
	}

	_, err := adapter.UploadFile(context.Background(), &multipart.FileHeader{Filename: ""}, "key")
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to open file")
}

func TestS3Adapter_UploadFile_UploadError(t *testing.T) {
	t.Parallel()
	adapter := newMockS3Adapter(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	file := newMultipartFileHeader(t, "upload.txt", "text/plain", []byte("file-data"))
	_, err := adapter.UploadFile(context.Background(), file, "uploads/upload.txt")
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to upload to S3")
}

func TestS3Adapter_UploadFiles_Success(t *testing.T) {
	t.Parallel()
	adapter := newMockS3Adapter(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	files := []*multipart.FileHeader{
		newMultipartFileHeader(t, "a.txt", "text/plain", []byte("a")),
		newMultipartFileHeader(t, "b.txt", "text/plain", []byte("b")),
	}

	result, err := adapter.UploadFiles(context.Background(), files)
	require.NoError(t, err)
	require.Len(t, result.Files, 2)
}

func TestS3Adapter_UploadFiles_Error(t *testing.T) {
	t.Parallel()
	adapter := newMockS3Adapter(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	files := []*multipart.FileHeader{
		newMultipartFileHeader(t, "a.txt", "text/plain", []byte("a")),
	}

	_, err := adapter.UploadFiles(context.Background(), files)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to upload file")
}

func TestS3Adapter_GetPresignedURL_Error(t *testing.T) {
	t.Parallel()
	testutil.InitLogger()

	cfg, err := awsconfig.LoadDefaultConfig(
		context.Background(),
		awsconfig.WithRegion("us-east-1"),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("", "", "")),
	)
	require.NoError(t, err)

	adapter := &S3Adapter{
		config: &config.Config{S3Region: "us-east-1", S3PresignedURLDuration: time.Hour},
		client: s3.NewFromConfig(cfg),
		bucket: "test-bucket",
	}

	_, err = adapter.GetPresignedURL(context.Background(), "missing.txt")
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to generate presigned URL")
}
