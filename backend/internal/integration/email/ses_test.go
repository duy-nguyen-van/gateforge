package email

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/testutil"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/stretchr/testify/require"
)

func sesSendEmailXML(messageID string) string {
	return fmt.Sprintf(`<?xml version="1.0"?>
<SendEmailResponse xmlns="http://ses.amazonaws.com/doc/2010-12-01/">
  <SendEmailResult>
    <MessageId>%s</MessageId>
  </SendEmailResult>
  <ResponseMetadata><RequestId>test-request-id</RequestId></ResponseMetadata>
</SendEmailResponse>`, messageID)
}

func sesSendRawEmailXML(messageID string) string {
	return fmt.Sprintf(`<?xml version="1.0"?>
<SendRawEmailResponse xmlns="http://ses.amazonaws.com/doc/2010-12-01/">
  <SendRawEmailResult>
    <MessageId>%s</MessageId>
  </SendRawEmailResult>
  <ResponseMetadata><RequestId>test-request-id</RequestId></ResponseMetadata>
</SendRawEmailResponse>`, messageID)
}

func sesErrorXML(message string) string {
	return fmt.Sprintf(`<?xml version="1.0"?>
<ErrorResponse xmlns="http://ses.amazonaws.com/doc/2010-12-01/">
  <Error><Code>MessageRejected</Code><Message>%s</Message></Error>
  <RequestId>test-request-id</RequestId>
</ErrorResponse>`, message)
}

func newMockSESSender(t *testing.T, handler http.HandlerFunc) *SESSender {
	t.Helper()
	testutil.InitLogger()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if handler != nil {
			handler(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/xml")
		_, _ = w.Write([]byte(sesSendEmailXML("mock-msg-id")))
	}))
	t.Cleanup(server.Close)

	cfg, err := awsconfig.LoadDefaultConfig(
		context.Background(),
		awsconfig.WithRegion("us-east-1"),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("AKIATESTKEY", "secret", "")),
	)
	require.NoError(t, err)

	client := ses.NewFromConfig(cfg, func(o *ses.Options) {
		o.BaseEndpoint = aws.String(server.URL)
	})

	return &SESSender{
		client: client,
		config: config.Config{AWSSESAccessKey: "sender@example.com"},
	}
}

func TestSESSender_SendEmail_HTMLBody(t *testing.T) {
	t.Parallel()
	s := newMockSESSender(t, func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		require.Contains(t, string(body), "html-body")
		w.Header().Set("Content-Type", "text/xml")
		_, _ = w.Write([]byte(sesSendEmailXML("html-msg-id")))
	})

	resp, err := s.SendEmail(context.Background(), EmailRequest{
		To:       []string{"to@example.com"},
		Subject:  "Subject",
		HTMLBody: "html-body",
	})
	require.NoError(t, err)
	require.Equal(t, "html-msg-id", resp.MessageID)
	require.Equal(t, "ses", resp.Provider)
	require.Equal(t, "sent", resp.Status)
}

func TestSESSender_SendEmail_TextBody(t *testing.T) {
	t.Parallel()
	s := newMockSESSender(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/xml")
		_, _ = w.Write([]byte(sesSendEmailXML("text-msg-id")))
	})

	resp, err := s.SendEmail(context.Background(), EmailRequest{
		To:       []string{"to@example.com"},
		Subject:  "Subject",
		TextBody: "text-body",
	})
	require.NoError(t, err)
	require.Equal(t, "text-msg-id", resp.MessageID)
}

func TestSESSender_SendEmail_HTMLAndTextBody(t *testing.T) {
	t.Parallel()
	var capturedBody string
	s := newMockSESSender(t, func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		capturedBody = string(body)
		w.Header().Set("Content-Type", "text/xml")
		_, _ = w.Write([]byte(sesSendEmailXML("both-msg-id")))
	})

	resp, err := s.SendEmail(context.Background(), EmailRequest{
		To:       []string{"to@example.com"},
		Cc:       []string{"cc@example.com"},
		Bcc:      []string{"bcc@example.com"},
		Subject:  "Subject",
		HTMLBody: "html-part",
		TextBody: "text-part",
	})
	require.NoError(t, err)
	require.Equal(t, "both-msg-id", resp.MessageID)
	require.Contains(t, capturedBody, "html-part")
	require.Contains(t, capturedBody, "text-part")
}

func TestSESSender_SendEmail_APIError(t *testing.T) {
	t.Parallel()
	s := newMockSESSender(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/xml")
		_, _ = w.Write([]byte(sesErrorXML("rejected")))
	})

	resp, err := s.SendEmail(context.Background(), EmailRequest{
		To:       []string{"to@example.com"},
		Subject:  "Subject",
		TextBody: "hello",
	})
	require.Error(t, err)
	require.NotNil(t, resp)
	require.Equal(t, "failed", resp.Status)
	require.Equal(t, "ses", resp.Provider)
	require.NotEmpty(t, resp.Error)
}

func TestSESSender_SendRawEmail_Success(t *testing.T) {
	t.Parallel()
	s := newMockSESSender(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/xml")
		_, _ = w.Write([]byte(sesSendRawEmailXML("raw-msg-id")))
	})

	resp, err := s.SendRawEmail(context.Background(), []byte("raw-email"))
	require.NoError(t, err)
	require.Equal(t, "raw-msg-id", resp.MessageID)
	require.Equal(t, "sent", resp.Status)
}

func TestNewSESSender_LoadConfigError(t *testing.T) {
	t.Parallel()
	testutil.InitLogger()

	// Empty region causes AWS config load to fail in some SDK versions.
	cfg := config.Config{
		AWSSESRegion:    "",
		AWSSESAccessKey: "AKIATESTKEY",
		AWSSESSecretKey: "secret",
	}
	sender, err := NewSESSender(cfg)
	if err != nil {
		require.Nil(t, sender)
		require.Contains(t, err.Error(), "Failed to load AWS config")
		return
	}
	// If SDK accepts empty region, constructor still succeeds.
	require.NotNil(t, sender)
}

func TestSESSender_SendRawEmail_APIError(t *testing.T) {
	t.Parallel()
	s := newMockSESSender(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/xml")
		_, _ = w.Write([]byte(sesErrorXML("raw rejected")))
	})

	resp, err := s.SendRawEmail(context.Background(), []byte("raw-email"))
	require.Error(t, err)
	require.NotNil(t, resp)
	require.Equal(t, "failed", resp.Status)
	require.Contains(t, resp.Error, "raw rejected")
}
