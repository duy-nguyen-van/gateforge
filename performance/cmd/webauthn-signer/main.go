package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/descope/virtualwebauthn"
)

type fixtureFile struct {
	Users []fixtureUser `json:"users"`
}

type fixtureUser struct {
	Email         string          `json:"email"`
	RPID          string          `json:"rp_id"`
	Origin        string          `json:"origin"`
	Credential    json.RawMessage `json:"credential"`
	Authenticator json.RawMessage `json:"authenticator"`
}

type assertRequest struct {
	Email   string          `json:"email"`
	Options json.RawMessage `json:"options"`
}

type assertResponse struct {
	Credential json.RawMessage `json:"credential"`
}

func main() {
	addr := flag.String("addr", "127.0.0.1:9091", "listen address")
	fixturesPath := flag.String("fixtures", "", "passkeys fixtures JSON")
	flag.Parse()
	if *fixturesPath == "" {
		*fixturesPath = ".data/passkeys.json"
	}

	raw, err := os.ReadFile(*fixturesPath)
	if err != nil {
		log.Fatalf("read fixtures: %v", err)
	}
	var ff fixtureFile
	if err := json.Unmarshal(raw, &ff); err != nil {
		log.Fatalf("parse fixtures: %v", err)
	}
	byEmail := make(map[string]fixtureUser, len(ff.Users))
	for _, u := range ff.Users {
		byEmail[u.Email] = u
	}
	if len(byEmail) == 0 {
		log.Fatal("no fixtures loaded")
	}

	var mu sync.Mutex
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("POST /assert", func(w http.ResponseWriter, r *http.Request) {
		var req assertRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		mu.Lock()
		fx, ok := byEmail[req.Email]
		mu.Unlock()
		if !ok {
			http.Error(w, "unknown email", http.StatusNotFound)
			return
		}

		var credential virtualwebauthn.Credential
		if err := json.Unmarshal(fx.Credential, &credential); err != nil {
			http.Error(w, "bad credential fixture: "+err.Error(), http.StatusInternalServerError)
			return
		}
		var authenticator virtualwebauthn.Authenticator
		if err := json.Unmarshal(fx.Authenticator, &authenticator); err != nil {
			http.Error(w, "bad authenticator fixture: "+err.Error(), http.StatusInternalServerError)
			return
		}
		authenticator.AddCredential(credential)

		parsed, err := virtualwebauthn.ParseAssertionOptions(string(req.Options))
		if err != nil {
			http.Error(w, "bad options: "+err.Error(), http.StatusBadRequest)
			return
		}
		rp := virtualwebauthn.RelyingParty{
			Name:   "GateForge Perf",
			ID:     fx.RPID,
			Origin: fx.Origin,
		}
		assertion := virtualwebauthn.CreateAssertionResponse(rp, authenticator, credential, *parsed)

		// Persist sign counter bump back into fixture map for subsequent assertions
		mu.Lock()
		fx.Authenticator, _ = json.Marshal(authenticator)
		byEmail[req.Email] = fx
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(assertResponse{Credential: json.RawMessage(assertion)})
	})

	fmt.Printf("webauthn-signer listening on http://%s (fixtures=%d)\n", *addr, len(byEmail))
	log.Fatal(http.ListenAndServe(*addr, mux))
}
