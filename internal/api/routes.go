package api

import (
	"go-cicd-engine/internal/store"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/jobs", handleListJobs)
	mux.HandleFunc("/api/jobs/", handleJobDetail) // also handles /log
}

func handleListJobs(w http.ResponseWriter, r *http.Request) {
	jobs, err := store.GetAllJobs()
	if err != nil {
		log.Printf("❌ Failed to get jobs: %v", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	respondJSON(w, jobs)
}

func handleJobDetail(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/jobs/")

	if strings.HasSuffix(id, "/log") {
		jobID := strings.TrimSuffix(id, "/log")
		sendLog(w, jobID)
		return
	}

	job, err := store.GetJob(id)
	if err != nil {
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	respondJSON(w, job)
}


func sendLog(w http.ResponseWriter, jobID string) {
	path := filepath.Join("logs", jobID+".log")
	f, err := os.Open(path)
	if err != nil {
		http.Error(w, "Log not found", http.StatusNotFound)
		return
	}
	defer f.Close()

	w.Header().Set("Content-Type", "text/plain")
	io.Copy(w, f)
}

func respondJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("❌ JSON encode error: %v", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
	}
}

