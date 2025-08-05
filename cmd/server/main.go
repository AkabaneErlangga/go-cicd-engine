package main

import (
	"log"
	"net/http"
	
	"go-cicd-engine/internal/api"
	"go-cicd-engine/internal/job"
	"go-cicd-engine/internal/webhook"
	"go-cicd-engine/internal/store"
)

func main() {
	if err := store.Init(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	log.Println("✅ Database initialized successfully")

	mux := http.NewServeMux()

	job.StartWorker()
	mux.HandleFunc("/webhook", webhook.Handler)
	api.RegisterRoutes(mux)

	log.Println("🚀 Listening on :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

