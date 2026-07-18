package main

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type codeEntry struct {
	Code         string `json:"code"`
	CodeVerifier string `json:"code_verifier"`
	ClientID     string `json:"client_id"`
	RedirectURI  string `json:"redirect_uri"`
}

func main() {
	count := flag.Int("n", 100000, "number of authorization codes to seed")
	outPath := flag.String("out", "", "output JSON path (default: performance/.data/codes.json)")
	email := flag.String("email", "perf-token@example.com", "bench user email (must exist)")
	clientID := flag.String("client-id", "oidc-dev", "OAuth client_id")
	redirectURI := flag.String("redirect-uri", "http://localhost:3000/callback", "redirect_uri stored on codes")
	dsn := flag.String("dsn", "", "Postgres DSN (default from POSTGRES_* env)")
	flag.Parse()

	if *outPath == "" {
		*outPath = defaultDataPath("codes.json")
	}
	databaseURL := *dsn
	if databaseURL == "" {
		databaseURL = postgresDSNFromEnv()
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()
	db.SetMaxOpenConns(8)
	ctx := context.Background()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("ping db: %v", err)
	}

	var userID, tenantID string
	err = db.QueryRowContext(ctx, `SELECT id FROM users WHERE email = $1 AND deleted_at IS NULL LIMIT 1`, strings.ToLower(*email)).Scan(&userID)
	if err != nil {
		log.Fatalf("bench user %q not found — register it first (make -C performance seed-user): %v", *email, err)
	}
	err = db.QueryRowContext(ctx, `
		SELECT tenant_id FROM tenant_memberships
		WHERE user_id = $1 AND deleted_at IS NULL
		ORDER BY created_at ASC LIMIT 1`, userID).Scan(&tenantID)
	if err != nil {
		log.Fatalf("no tenant membership for %s: %v", *email, err)
	}

	var clientRecordID sql.NullString
	err = db.QueryRowContext(ctx, `
		SELECT id FROM clients
		WHERE client_id = $1 AND tenant_id = $2 AND deleted_at IS NULL
		LIMIT 1`, *clientID, tenantID).Scan(&clientRecordID)
	if err != nil {
		log.Fatalf("client %q not found for tenant %s: %v", *clientID, tenantID, err)
	}

	expiresAt := time.Now().UTC().Add(24 * time.Hour)
	entries := make([]codeEntry, 0, *count)
	const batchSize = 500

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatalf("begin: %v", err)
	}
	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO authorization_codes (
			id, created_at, updated_at, code, tenant_id, oauth_client_id, user_id,
			redirect_uri, code_challenge, code_challenge_method, scope, expires_at, client_record_id
		) VALUES (
			$1, now(), now(), $2, $3, $4, $5,
			$6, $7, 'S256', 'openid email profile', $8, $9
		)`)
	if err != nil {
		_ = tx.Rollback()
		log.Fatalf("prepare: %v", err)
	}
	defer stmt.Close()

	for i := 0; i < *count; i++ {
		verifier := randomURLToken(48)
		challenge := pkceChallenge(verifier)
		code := randomURLToken(32)
		id := uuid.Must(uuid.NewV7()).String()

		var clientRec any
		if clientRecordID.Valid {
			clientRec = clientRecordID.String
		}
		if _, err := stmt.ExecContext(ctx, id, code, tenantID, *clientID, userID, *redirectURI, challenge, expiresAt, clientRec); err != nil {
			_ = tx.Rollback()
			log.Fatalf("insert %d: %v", i, err)
		}
		entries = append(entries, codeEntry{
			Code:         code,
			CodeVerifier: verifier,
			ClientID:     *clientID,
			RedirectURI:  *redirectURI,
		})
		if (i+1)%batchSize == 0 {
			if err := tx.Commit(); err != nil {
				log.Fatalf("commit batch: %v", err)
			}
			tx, err = db.BeginTx(ctx, nil)
			if err != nil {
				log.Fatalf("begin: %v", err)
			}
			stmt, err = tx.PrepareContext(ctx, `
				INSERT INTO authorization_codes (
					id, created_at, updated_at, code, tenant_id, oauth_client_id, user_id,
					redirect_uri, code_challenge, code_challenge_method, scope, expires_at, client_record_id
				) VALUES (
					$1, now(), now(), $2, $3, $4, $5,
					$6, $7, 'S256', 'openid email profile', $8, $9
				)`)
			if err != nil {
				_ = tx.Rollback()
				log.Fatalf("prepare: %v", err)
			}
			fmt.Fprintf(os.Stderr, "seeded %d / %d codes\n", i+1, *count)
		}
	}
	if err := tx.Commit(); err != nil {
		log.Fatalf("final commit: %v", err)
	}

	if err := os.MkdirAll(filepath.Dir(*outPath), 0o755); err != nil {
		log.Fatalf("mkdir: %v", err)
	}
	f, err := os.Create(*outPath)
	if err != nil {
		log.Fatalf("create out: %v", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(entries); err != nil {
		log.Fatalf("encode: %v", err)
	}
	fmt.Printf("wrote %d codes to %s\n", len(entries), *outPath)
}

func pkceChallenge(verifier string) string {
	sum := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}

func randomURLToken(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		log.Fatalf("rand: %v", err)
	}
	return base64.RawURLEncoding.EncodeToString(b)
}

func postgresDSNFromEnv() string {
	if v := os.Getenv("DATABASE_URL"); v != "" {
		return v
	}
	host := envOr("POSTGRES_HOST", "localhost")
	port := envOr("POSTGRES_PORT", "5432")
	user := envOr("POSTGRES_USER", "postgres")
	pass := envOr("POSTGRES_PASSWORD", "postgres")
	db := envOr("POSTGRES_DB", "iam")
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, db)
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
