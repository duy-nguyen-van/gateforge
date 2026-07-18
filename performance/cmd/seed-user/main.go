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
	"time"
)

func main() {
	baseURL := flag.String("base", envOr("PERF_BASE_URL", "http://localhost:3000"), "API base URL")
	email := flag.String("email", "perf-token@example.com", "bench user email")
	password := flag.String("password", "perf-token-passphrase-long!!", "bench user password")
	flag.Parse()

	client := &http.Client{Timeout: 30 * time.Second}
	body := map[string]any{
		"email":      *email,
		"password":   *password,
		"first_name": "Perf",
		"last_name":  "Token",
		"tenant_id":  envOr("PERF_TENANT_ID", "00000000-0000-0000-0000-000000000001"),
	}
	b, _ := json.Marshal(body)
	req, err := http.NewRequest(http.MethodPost, *baseURL+"/api/v1/register", bytes.NewReader(b))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		fmt.Printf("registered %s\n", *email)
		return
	}
	// Already exists is fine
	fmt.Printf("register status %d (ok if user exists): %s\n", resp.StatusCode, truncate(string(raw), 200))
}

func envOr(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
