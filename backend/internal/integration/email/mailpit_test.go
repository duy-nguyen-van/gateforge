package email

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"

	"github.com/stretchr/testify/require"
)

func TestNewMailpitSender(t *testing.T) {
	t.Run("rejects invalid url", func(t *testing.T) {
		sender, err := NewMailpitSender(config.Config{MailpitBaseURL: "not-a-url", EmailFrom: "ops@gateforge.test"})
		require.Error(t, err)
		require.Nil(t, sender)
	})

	t.Run("requires from address", func(t *testing.T) {
		sender, err := NewMailpitSender(config.Config{MailpitBaseURL: "http://localhost:8025"})
		require.Error(t, err)
		require.Nil(t, sender)
	})

	t.Run("parses from mailbox", func(t *testing.T) {
		sender, err := NewMailpitSender(config.Config{
			MailpitBaseURL: "http://localhost:8025",
			EmailFrom:      "GateForge <support@gateforge.test>",
		})
		require.NoError(t, err)
		require.Equal(t, "support@gateforge.test", sender.from.Email)
		require.Equal(t, "GateForge", sender.from.Name)
	})
}

func TestProvideEmailSender_Mailpit(t *testing.T) {
	sender, err := ProvideEmailSender(&config.Config{
		EmailProvider:  constants.EmailProviderMailpit,
		MailpitBaseURL: "http://localhost:8025",
		EmailFrom:      "ops@gateforge.test",
	})
	require.NoError(t, err)
	_, ok := sender.(*MailpitSender)
	require.True(t, ok)
}

func TestMailpitSender_SendEmail(t *testing.T) {
	t.Run("posts json send payload", func(t *testing.T) {
		var got mailpitSendRequest
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "/api/v1/send", r.URL.Path)
			body, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			require.NoError(t, json.Unmarshal(body, &got))
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"ID":"mailpit-msg-1"}`))
		}))
		t.Cleanup(server.Close)

		sender, err := NewMailpitSender(config.Config{
			MailpitBaseURL: server.URL,
			EmailFrom:      "GateForge <ops@gateforge.test>",
		})
		require.NoError(t, err)

		resp, err := sender.SendEmail(context.Background(), EmailRequest{
			To:       []string{"member@example.com"},
			Cc:       []string{"cc@example.com"},
			Subject:  "You were added to Acme",
			TextBody: "plain",
			HTMLBody: "<p>html</p>",
		})
		require.NoError(t, err)
		require.Equal(t, "mailpit-msg-1", resp.MessageID)
		require.Equal(t, constants.EmailProviderMailpit, resp.Provider)
		require.Equal(t, "ops@gateforge.test", got.From.Email)
		require.Equal(t, "member@example.com", got.To[0].Email)
		require.Equal(t, "cc@example.com", got.Cc[0].Email)
		require.Equal(t, "You were added to Acme", got.Subject)
		require.Equal(t, "plain", got.Text)
		require.Equal(t, "<p>html</p>", got.HTML)
	})

	t.Run("maps http errors", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error":"invalid from"}`))
		}))
		t.Cleanup(server.Close)

		sender, err := NewMailpitSender(config.Config{MailpitBaseURL: server.URL, EmailFrom: "ops@gateforge.test"})
		require.NoError(t, err)

		_, err = sender.SendEmail(context.Background(), EmailRequest{
			To:       []string{"member@example.com"},
			Subject:  "You were added to Acme",
			TextBody: "plain",
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "Failed to send email via Mailpit")
	})
}
