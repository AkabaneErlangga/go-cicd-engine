package store

import (
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Job struct {
	ID         string `gorm:"primaryKey"`
	RepoURL    string
	Branch     string
	Status     string
	LogPath    string
	CreatedAt  time.Time
	FinishedAt *time.Time
}

var db *gorm.DB

func Init() error {
	var err error
	db, err = gorm.Open(sqlite.Open("cicd.db"), &gorm.Config{})
	if err != nil {
		return err
	}
	return db.AutoMigrate(&Job{})
}

func CreateJob(j Job) error {
	return db.Create(&j).Error
}

func UpdateJobStatus(id string, status string, finishedAt *time.Time) error {
	return db.Model(&Job{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":      status,
			"finished_at": finishedAt,
		}).Error
}

func GetAllJobs() ([]Job, error) {
	var jobs []Job
	err := db.Find(&jobs).Error
	return jobs, err
}

func GetJob(id string) (*Job, error) {
	var job Job
	if err := db.First(&job, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &job, nil
}

func UpdateJobLogPath(id string, path string) error {
	return db.Model(&Job{}).
		Where("id = ?", id).
		Update("log_path", path).Error
}

