package model

import "time"

type Job struct {
	ID 				string
	RepoURL		string
	Branch		string
	Author		string	
	CommitMsg	string
	CommitURL	string
	CommitTime	time.Time
	Status		string // "pending", "success", "failure"
	CreatedAt	string
}
