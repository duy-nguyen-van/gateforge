package email

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/mail"
	"net/url"
	"strings"
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"
)

type MailpitSender struct {
	baseURL  string
	from     mailpitAddress
	client   *http.Client
	sendPath string
}

type mailpitAddress struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

type mailpitSendRequest struct {
	From    mailpitAddress   `json:"from"`
	To      []mailpitAddress `json:"to"`
	Cc      []mailpitAddress `json:"cc,omitempty"`
	Bcc     []mailpitAddress `json:"bcc,omitempty"`
	Subject string           `json:"subject"`
	Text    string           `json:"text,omitempty"`
	HTML    string           `json:"html,omitempty"`
}

type mailpitSendResponse struct {
	ID string `json:"ID"`
	Id string `json:"id"`
}

func NewMailpitSender(cfg config.Config) (*MailpitSender, error) {
	base := strings.TrimRight(strings.TrimSpace(cfg.MailpitBaseURL), "/")
	if base == "" {
		base = "http://localhost:8025"
	}
	parsed, err := url.Parse(base)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		cause := err
		if cause == nil {
			cause = fmt.Errorf("missing scheme or host")
		}
		return nil, errors.InternalError("Mailpit base URL is invalid", cause).
			WithOperation("initialize_mailpit").
			WithResource("mailpit")
	}
	fromEmail, fromName := mailboxFromConfig(cfg.EmailFrom, cfg.EmailFromName)
	if fromEmail == "" {
		return nil, errors.InternalError("Email from address is required", fmt.Errorf("EMAIL_FROM is empty")).
			WithOperation("initialize_mailpit").
			WithResource("mailpit")
	}
	timeout := cfg.HTTPClientTimeout
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	return &MailpitSender{
		baseURL:  strings.TrimRight(parsed.String(), "/"),
		from:     mailpitAddress{Email: fromEmail, Name: fromName},
		client:   &http.Client{Timeout: timeout},
		sendPath: "/api/v1/send",
	}, nil
}

func (s *MailpitSender) SendEmail(ctx context.Context, request EmailRequest) (*EmailResponse, error) {
	if len(request.To) == 0 {
		return nil, errors.ExternalServiceError("at least one recipient is required", nil).
			WithOperation("send_email").
			WithResource("mailpit")
	}
	if request.HTMLBody == "" && request.TextBody == "" {
		return nil, errors.ExternalServiceError("either HTML or text content must be provided", nil).
			WithOperation("send_email").
			WithResource("mailpit")
	}
	payload, err := json.Marshal(mailpitSendRequest{
		From:    s.from,
		To:      mailpitRecipients(request.To),
		Cc:      mailpitRecipients(request.Cc),
		Bcc:     mailpitRecipients(request.Bcc),
		Subject: request.Subject,
		Text:    request.TextBody,
		HTML:    request.HTMLBody,
	})
	if err != nil {
		return nil, errors.InternalError("Failed to encode Mailpit payload", err).
			WithOperation("send_email").
			WithResource("mailpit")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.baseURL+s.sendPath, bytes.NewReader(payload))
	if err != nil {
		return nil, errors.InternalError("Failed to build Mailpit request", err).
			WithOperation("send_email").
			WithResource("mailpit")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, errors.ExternalServiceError("Failed to send email via Mailpit", err).
			WithOperation("send_email").
			WithResource("mailpit")
	}
	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, errors.ExternalServiceError("Failed to read Mailpit response", err).
			WithOperation("send_email").
			WithResource("mailpit")
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errors.ExternalServiceError("Failed to send email via Mailpit", fmt.Errorf("status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))).
			WithOperation("send_email").
			WithResource("mailpit")
	}
	var parsed mailpitSendResponse
	if len(body) > 0 {
		if err := json.Unmarshal(body, &parsed); err != nil {
			return nil, errors.ExternalServiceError("Failed to decode Mailpit response", err).
				WithOperation("send_email").
				WithResource("mailpit")
		}
	}
	id := parsed.ID
	if id == "" {
		id = parsed.Id
	}
	return &EmailResponse{MessageID: id, Provider: constants.EmailProviderMailpit, Status: "sent"}, nil
}

func (s *MailpitSender) SendRawEmail(ctx context.Context, rawData []byte) (*EmailResponse, error) {
	return nil, errors.InternalError("Mailpit HTTP send does not support raw MIME", nil).
		WithOperation("send_raw_email").
		WithResource("mailpit")
}

func mailpitRecipients(addrs []string) []mailpitAddress {
	if len(addrs) == 0 {
		return nil
	}
	out := make([]mailpitAddress, 0, len(addrs))
	for _, addr := range addrs {
		emailAddr, name := mailboxFromConfig(addr, "")
		if emailAddr == "" {
			continue
		}
		out = append(out, mailpitAddress{Email: emailAddr, Name: name})
	}
	return out
}

func mailboxFromConfig(raw, fallbackName string) (emailAddr, name string) {
	raw = strings.TrimSpace(raw)
	fallbackName = strings.TrimSpace(fallbackName)
	if raw == "" {
		return "", fallbackName
	}
	if addr, err := mail.ParseAddress(raw); err == nil {
		emailAddr = addr.Address
		name = addr.Name
	} else {
		emailAddr = raw
	}
	if name == "" {
		name = fallbackName
	}
	return emailAddr, name
}
