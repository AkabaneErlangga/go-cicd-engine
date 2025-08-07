package model

type Job struct {
	ID 				string
	RepoURL		string
	Branch		string
	Author		string	
	commitMsg	string
	status		string // "pending", "success", "failure"
	CreatedAt	string
}
