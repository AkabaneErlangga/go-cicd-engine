package job

import (
	"log"

	"go-cicd-engine/internal/model"
	"go-cicd-engine/internal/runner"
)

var jobQueue = make(chan model.Job, 10) // Queue ukuran 10 job

func StartWorker() {
	go func() {
		for job := range jobQueue {
			log.Printf("⚙️  Running job: %s\n", job.ID)
			err := runner.Execute(job)
			if err != nil {
				log.Printf("❌ Job failed: %v", err)
			} else {
				log.Println("✅ Job completed")
			}
		}
	}()
}

func Enqueue(job model.Job) {
	select {
	case jobQueue <- job:
		log.Printf("📦 Job queued: %s", job.ID)
	default:
		log.Println("❗ Queue full. Dropping job.")
	}
}

