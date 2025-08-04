package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"go-cicd-engine/internal/job"
	"go-cicd-engine/internal/model"
	"io"
	"log"
	"net/http"
	"os"
	"github.com/google/uuid"	
)

// GITHUB_WEBHOOK_SECRET disimpan di env var
var secret = os.Getenv("GITHUB_WEBHOOK_SECRET")

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Read error", http.StatusBadRequest)
		return
	}

	sig := r.Header.Get("X-Hub-Signature-256")
	if !verifySignature(body, sig) {
		http.Error(w, "Signature mismatch", http.StatusUnauthorized)
		return
	}

	// TODO: Parse payload, push job ke queue
	log.Println("✅ Webhook received & verified!")

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Webhook accepted"))

	job.Enqueue(model.Job{
		ID: uuid.NewString(),
		RepoURL: "https://github.com/AkabaneErlangga/tes-cicd-engine",
		Branch: "main",
		Commands: []string{"cat README.md"},
	})
}

func verifySignature(payload []byte, signature string) bool {
	if secret == "" {
		log.Println("⚠️  Warning: no secret set")
		return true // skip verification (DEV ONLY)
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(expected), []byte(signature))
}

