package job

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"go-cicd-engine/internal/model"
	"go-cicd-engine/internal/notifier"
	"go-cicd-engine/internal/runner"
	"go-cicd-engine/internal/store"
)

var jobQueue = make(chan model.Job, 10) // Queue ukuran 10 job

func StartWorker() {
	go func() {
		for job := range jobQueue {
			SetStatus(job.ID, StatusRunning)

			logPath := filepath.Join("logs", fmt.Sprintf("%s.log", job.ID));
			f, err := os.Create(logPath)

			if err != nil {
				log.Printf("❗ Error creating log file for job %s: %v", job.ID, err)
				continue
			}
			defer f.Close()

			store.UpdateJobStatus(job.ID, string(StatusRunning), nil)

			log.Printf("⚙️  Running job: %s\n", job.ID)
			notifyJob(job, "running")
			err = runner.Execute(job, f)
			now := time.Now()

			switch {
			case errors.Is(err, runner.ErrTimeout):
				log.Printf("⏰ Job %s timed out", job.ID)
				SetStatus(job.ID, StatusTimedOut)
				store.UpdateJobStatus(job.ID, string(StatusTimedOut), &now)
				notifyJob(job, "timed out")
			case err != nil:
				log.Printf("❌ Job %s failed: %v", job.ID, err)
				SetStatus(job.ID, StatusFailed)
				store.UpdateJobStatus(job.ID, string(StatusFailed), &now)
				notifyJob(job, "failed")

			default:
				log.Printf("✅ Job %s completed", job.ID)
				SetStatus(job.ID, StatusDone)
				store.UpdateJobStatus(job.ID, string(StatusDone), &now)
				notifyJob(job, "completed")
			}
			store.UpdateJobLogPath(job.ID, logPath)
		}
	}()
}

func Enqueue(j model.Job) {
	SetStatus(j.ID, StatusQueued)

	err := store.CreateJob(store.Job{
		ID:        j.ID,
		RepoURL:   j.RepoURL,
		Branch:    j.Branch,
		Status:    string(StatusQueued),
		CreatedAt: time.Now(),
	})
	if err != nil {
		log.Printf("❌ Failed to store job: %v", err)
	}

	select {
	case jobQueue <- j:
		log.Printf("📦 Job queued: %s", j.ID)
	default:
		log.Println("❗ Queue full. Dropping job.")
	}
}

func notifyJob(job model.Job, status string) {
	notifier.NotifyDiscord(notifier.Event{
		Status: status,
		JobID:  job.ID,
		Repo:   job.RepoURL,
		Branch: job.Branch,
		Author: job.Author,
		CommitMsg: job.CommitMsg,
		CommitURL: job.CommitURL,
		CommitTime: job.CommitTime,
		Timestamp: time.Now(),
	})
}

