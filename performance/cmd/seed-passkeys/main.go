package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/descope/virtualwebauthn"
)

type envelope struct {
	Data json.RawMessage `json:"data"`
}

type loginData struct {
	AccessToken string `json:"access_token"`
}

type waStartData struct {
	Options      json.RawMessage `json:"options"`
	SessionToken string         `json:"session_token"`
}

type fixtureCredential struct {
	Email        string          `json:"email"`
	Password     string          `json:"password"`
	RPID         string          `json:"rp_id"`
	Origin       string          `json:"origin"`
	UserHandle   string          `json:"user_handle,omitempty"`
	Credential   json.RawMessage `json:"credential"`
	Authenticator json.RawMessage `json:"authenticator"`
}

type fixtureFile struct {
	Users []fixtureCredential `json:"users"`
}

func main() {
	baseURL := flag.String("base", envOr("PERF_BASE_URL", "http://localhost:3000"), "API base URL")
	count := flag.Int("n", 20, "number of passkey users to seed")
	outPath := flag.String("out", "", "output fixtures JSON")
	password := flag.String("password", "perf-passkey-passphrase-long!!", "password for seeded users")
	rpID := flag.String("rp-id", envOr("WEBAUTHN_RP_ID", "localhost"), "WebAuthn RP ID")
	origin := flag.String("origin", envOr("WEBAUTHN_RP_ORIGIN", "http://localhost:3000"), "WebAuthn origin")
	flag.Parse()
	if *outPath == "" {
		*outPath = defaultDataPath("passkeys.json")
	}

	client := &http.Client{Timeout: 30 * time.Second}
	rp := virtualwebauthn.RelyingParty{
		Name:   "GateForge Perf",
		ID:     *rpID,
		Origin: *origin,
	}

	out := fixtureFile{Users: make([]fixtureCredential, 0, *count)}
	for i := 0; i < *count; i++ {
		email := fmt.Sprintf("perf-passkey-%03d@example.com", i)
		if err := ensureRegistered(client, *baseURL, email, *password); err != nil {
			log.Fatalf("register %s: %v", email, err)
		}
		access, err := login(client, *baseURL, email, *password)
		if err != nil {
			log.Fatalf("login %s: %v", email, err)
		}

		authenticator := virtualwebauthn.NewAuthenticator()
		credential := virtualwebauthn.NewCredential(virtualwebauthn.KeyTypeEC2)

		startBody, err := postJSON(client, *baseURL+"/api/v1/webauthn/register/start", map[string]any{
			"device_name": "perf-soft-auth",
		}, access)
		if err != nil {
			log.Fatalf("register/start %s: %v", email, err)
		}
		var start waStartData
		if err := json.Unmarshal(startBody, &start); err != nil {
			log.Fatalf("parse register/start: %v", err)
		}
		parsed, err := virtualwebauthn.ParseAttestationOptions(string(start.Options))
		if err != nil {
			log.Fatalf("parse attestation options: %v", err)
		}
		attestation := virtualwebauthn.CreateAttestationResponse(rp, authenticator, credential, *parsed)
		finishBody := struct {
			SessionToken string          `json:"session_token"`
			Credential   json.RawMessage `json:"credential"`
		}{
			SessionToken: start.SessionToken,
			Credential:   json.RawMessage(attestation),
		}
		_, err = postJSON(client, *baseURL+"/api/v1/webauthn/register/finish", finishBody, access)
		if err != nil {
			log.Fatalf("register/finish %s: %v", email, err)
		}

		credJSON, _ := json.Marshal(credential)
		authJSON, _ := json.Marshal(authenticator)
		out.Users = append(out.Users, fixtureCredential{
			Email:         email,
			Password:      *password,
			RPID:          *rpID,
			Origin:        *origin,
			Credential:    credJSON,
			Authenticator: authJSON,
		})
		fmt.Fprintf(os.Stderr, "seeded passkey user %s\n", email)
	}

	if err := os.MkdirAll(filepath.Dir(*outPath), 0o755); err != nil {
		log.Fatalf("mkdir: %v", err)
	}
	f, err := os.Create(*outPath)
	if err != nil {
		log.Fatalf("create: %v", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(out); err != nil {
		log.Fatalf("encode: %v", err)
	}
	fmt.Printf("wrote %d fixtures to %s\n", len(out.Users), *outPath)
}

func ensureRegistered(client *http.Client, base, email, password string) error {
	body := map[string]any{
		"email":      email,
		"password":   password,
		"first_name": "Perf",
		"last_name":  "Passkey",
		"tenant_id":  envOr("PERF_TENANT_ID", "00000000-0000-0000-0000-000000000001"),
	}
	b, _ := json.Marshal(body)
	req, err := http.NewRequest(http.MethodPost, base+"/api/v1/register", bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		return nil
	}
	if resp.StatusCode == http.StatusConflict || strings.Contains(string(raw), "already") {
		return nil
	}
	// Some installs return 400 for duplicate email
	if resp.StatusCode == http.StatusBadRequest && strings.Contains(strings.ToLower(string(raw)), "exist") {
		return nil
	}
	return fmt.Errorf("status %d: %s", resp.StatusCode, truncate(string(raw), 300))
}

func login(client *http.Client, base, email, password string) (string, error) {
	body := map[string]any{"email": email, "password": password}
	data, err := postJSON(client, base+"/api/v1/login", body, "")
	if err != nil {
		return "", err
	}
	var ld loginData
	if err := json.Unmarshal(data, &ld); err != nil {
		return "", err
	}
	if ld.AccessToken == "" {
		return "", fmt.Errorf("no access_token in login response: %s", truncate(string(data), 300))
	}
	return ld.AccessToken, nil
}

func postJSON(client *http.Client, url string, body any, bearer string) (json.RawMessage, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if bearer != "" {
		req.Header.Set("Authorization", "Bearer "+bearer)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("status %d: %s", resp.StatusCode, truncate(string(raw), 400))
	}
	var env envelope
	if err := json.Unmarshal(raw, &env); err != nil {
		return nil, fmt.Errorf("envelope: %w body=%s", err, truncate(string(raw), 200))
	}
	return env.Data, nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

func envOr(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func defaultDataPath(name string) string {
	if root := os.Getenv("PERF_DATA_DIR"); root != "" {
		return filepath.Join(root, name)
	}
	return filepath.Join(".", ".data", name)
}
