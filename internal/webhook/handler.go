package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"go-cicd-engine/internal/job"
	"go-cicd-engine/internal/model"

	"github.com/google/uuid"
)

// Ambil secret dari env
var secret = os.Getenv("GITHUB_WEBHOOK_SECRET")

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}

	// Signature dari header
	signature := r.Header.Get("X-Hub-Signature-256")
	if !verifySignature(body, signature) {
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		return
	}

	// Parse payload
	var payload struct {
		Repository struct {
			CloneURL string `json:"clone_url"`
		} `json:"repository"`
		Ref string `json:"ref"` // Format: refs/heads/main
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	branch := extractBranch(payload.Ref)
	if branch == "" {
		http.Error(w, "Unsupported ref", http.StatusBadRequest)
		return
	}

	jobID := uuid.NewString()
	log.Printf("📥 Webhook received for repo: %s, branch: %s", payload.Repository.CloneURL, branch)

	job.Enqueue(model.Job{
		ID:       jobID,
		RepoURL:  payload.Repository.CloneURL,
		Branch:   branch,
		Commands: nil, // Akan diisi dari .cicd.yaml saat run
	})

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Job accepted"))
}

func verifySignature(body []byte, signature string) bool {
	if secret == "" {
		log.Println("⚠️  No webhook secret set (SKIPPING VERIFY)")
		return true
	}
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(signature))
}

func extractBranch(ref string) string {
	const prefix = "refs/heads/"
	if len(ref) <= len(prefix) || ref[:len(prefix)] != prefix {
		return ""
	}
	return ref[len(prefix):]
}

