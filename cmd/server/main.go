package main

import (
	"log"
	"net/http"

	"go-cicd-engine/internal/job"
	"go-cicd-engine/internal/webhook"
)

func main() {
	job.StartWorker()
	http.HandleFunc("/webhook", webhook.Handler)

	log.Println("🚀 Listening on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

